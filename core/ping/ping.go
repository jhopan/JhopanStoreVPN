package ping

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"sync"
	"time"
)

const (
	DefaultPingURL = "https://dns.google"
	pingInterval   = 5 * time.Second
	pingTimeout    = 5 * time.Second
)

// Pinger periodically measures HTTP latency through the SOCKS proxy.
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

// Start begins the ping loop through the local SOCKS proxy.
func (p *Pinger) Start() {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.running {
		return
	}

	ctx, cancel := context.WithCancel(context.Background())
	p.cancel = cancel
	p.running = true

	proxyURL, _ := url.Parse("http://127.0.0.1:10809")
	transport := &http.Transport{
		Proxy: http.ProxyURL(proxyURL),
		DialContext: (&net.Dialer{
			Timeout: pingTimeout,
		}).DialContext,
		TLSHandshakeTimeout: pingTimeout,
	}
	client := &http.Client{
		Transport: transport,
		Timeout:   pingTimeout,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	go func() {
		defer client.CloseIdleConnections()
		ticker := time.NewTicker(pingInterval)
		defer ticker.Stop()
		p.doPing(ctx, client)
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				p.doPing(ctx, client)
			}
		}
	}()
}

// Stop halts the ping loop.
func (p *Pinger) Stop() {
	p.mu.Lock()
	defer p.mu.Unlock()
	if !p.running {
		return
	}
	p.running = false
	if p.cancel != nil {
		p.cancel()
	}
}

func (p *Pinger) doPing(ctx context.Context, client *http.Client) {
	start := time.Now()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, p.pingURL, nil)
	if err != nil {
		if p.onResult != nil {
			p.onResult(0, err)
		}
		return
	}
	resp, err := client.Do(req)
	if err != nil {
		if ctx.Err() != nil {
			return
		}
		if p.onResult != nil {
			p.onResult(0, err)
		}
		return
	}
	resp.Body.Close()
	latency := time.Since(start)
	if resp.StatusCode >= 400 {
		if p.onResult != nil {
			p.onResult(latency, fmt.Errorf("status %d", resp.StatusCode))
		}
		return
	}
	if p.onResult != nil {
		p.onResult(latency, nil)
	}
}
