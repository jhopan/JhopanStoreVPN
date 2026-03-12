package ping

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"strconv"
	"sync"
	"time"
)

const (
	DefaultPingURL = "https://dns.google"
	pingTimeout    = 5 * time.Second
	burstCount     = 3
	burstInterval  = 800 * time.Millisecond
	socksProxyAddr = "127.0.0.1:10808"
)

// Pinger fires a short burst of pings on connect, reports the best result, then stops.
// No background goroutine remains after the burst completes.
type Pinger struct {
	mu       sync.Mutex
	cancel   context.CancelFunc
	running  bool
	pingURL  string
	onResult func(latency time.Duration, err error)
}

// NewPinger creates a Pinger. If pingURL is empty, uses DefaultPingURL.
func NewPinger(pingURL string, onResult func(latency time.Duration, err error)) *Pinger {
	if pingURL == "" {
		pingURL = DefaultPingURL
	}
	return &Pinger{pingURL: pingURL, onResult: onResult}
}

// Start fires burstCount pings through the SOCKS5 proxy, reports the best
// latency result to onResult, then exits — no ongoing CPU or radio usage.
func (p *Pinger) Start() {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.running {
		return
	}

	ctx, cancel := context.WithCancel(context.Background())
	p.cancel = cancel
	p.running = true

	transport := &http.Transport{
		DialContext: func(ctx context.Context, _, addr string) (net.Conn, error) {
			return dialSOCKS5(ctx, socksProxyAddr, addr)
		},
		TLSHandshakeTimeout: pingTimeout,
		DisableKeepAlives:   true,
	}
	client := &http.Client{
		Transport: transport,
		Timeout:   pingTimeout,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	go func() {
		defer func() {
			client.CloseIdleConnections()
			cancel()
			p.mu.Lock()
			p.running = false
			p.mu.Unlock()
		}()

		var best time.Duration
		for i := 0; i < burstCount; i++ {
			select {
			case <-ctx.Done():
				return
			default:
			}

			latency, err := p.doPing(ctx, client)
			if err == nil {
				if best == 0 || latency < best {
					best = latency
					if p.onResult != nil {
						p.onResult(latency, nil)
					}
				}
			} else if i == burstCount-1 && best == 0 {
				if p.onResult != nil {
					p.onResult(0, fmt.Errorf("timeout"))
				}
			}

			if i < burstCount-1 {
				select {
				case <-ctx.Done():
					return
				case <-time.After(burstInterval):
				}
			}
		}
	}()
}

// Stop cancels an in-progress burst. Safe to call even after burst has finished.
func (p *Pinger) Stop() {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.cancel != nil {
		p.cancel()
	}
}

func (p *Pinger) doPing(ctx context.Context, client *http.Client) (time.Duration, error) {
	start := time.Now()
	req, err := http.NewRequestWithContext(ctx, http.MethodHead, p.pingURL, nil)
	if err != nil {
		return 0, err
	}
	resp, err := client.Do(req)
	if err != nil {
		if ctx.Err() != nil {
			return 0, ctx.Err()
		}
		return 0, err
	}
	resp.Body.Close()
	return time.Since(start), nil
}

// dialSOCKS5 opens a TCP connection through a SOCKS5 proxy without any
// external dependencies — implements just the no-auth subset of RFC 1928.
func dialSOCKS5(ctx context.Context, proxyAddr, targetAddr string) (net.Conn, error) {
	d := &net.Dialer{Timeout: pingTimeout}
	conn, err := d.DialContext(ctx, "tcp", proxyAddr)
	if err != nil {
		return nil, err
	}

	host, portStr, err := net.SplitHostPort(targetAddr)
	if err != nil {
		conn.Close()
		return nil, err
	}
	port, err := strconv.Atoi(portStr)
	if err != nil {
		conn.Close()
		return nil, err
	}

	// Greeting: version 5, one method, no-auth
	if _, err = conn.Write([]byte{0x05, 0x01, 0x00}); err != nil {
		conn.Close()
		return nil, err
	}
	reply := make([]byte, 2)
	if _, err = io.ReadFull(conn, reply); err != nil {
		conn.Close()
		return nil, err
	}
	if reply[0] != 0x05 || reply[1] != 0x00 {
		conn.Close()
		return nil, fmt.Errorf("socks5 auth rejected")
	}

	// CONNECT request using domain-name address type
	req := []byte{0x05, 0x01, 0x00, 0x03, byte(len(host))}
	req = append(req, []byte(host)...)
	req = append(req, byte(port>>8), byte(port&0xff))
	if _, err = conn.Write(req); err != nil {
		conn.Close()
		return nil, err
	}

	// Read response header
	hdr := make([]byte, 4)
	if _, err = io.ReadFull(conn, hdr); err != nil {
		conn.Close()
		return nil, err
	}
	if hdr[1] != 0x00 {
		conn.Close()
		return nil, fmt.Errorf("socks5 connect refused: code %d", hdr[1])
	}
	// Drain the bound-address field so conn is ready for TLS
	switch hdr[3] {
	case 0x01: // IPv4
		io.ReadFull(conn, make([]byte, 6))
	case 0x03: // domain
		buf := make([]byte, 1)
		io.ReadFull(conn, buf)
		io.ReadFull(conn, make([]byte, int(buf[0])+2))
	case 0x04: // IPv6
		io.ReadFull(conn, make([]byte, 18))
	}

	return conn, nil
}
