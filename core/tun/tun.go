package tun

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sync"
	"time"
)

// TunDevice represents a TUN network interface managed by external tun2socks binary
type TunDevice struct {
	Name      string
	IP        string
	Gateway   string
	DNS       []string
	MTU       int
	socksAddr string
	cmd       *exec.Cmd
	mu        sync.Mutex
	running   bool
}

// Config for TUN device creation
type Config struct {
	Name      string
	IP        string   // TUN interface IP, e.g., "10.0.0.2"
	Gateway   string   // TUN gateway, e.g., "10.0.0.1"
	DNS       []string // DNS servers
	MTU       int
	SocksAddr string   // SOCKS5 proxy address, e.g., "127.0.0.1:10809"
}

// NewTunDevice creates a new TUN device with the given configuration
func NewTunDevice(cfg Config) (*TunDevice, error) {
	if cfg.MTU == 0 {
		cfg.MTU = 1500
	}
	if cfg.IP == "" {
		cfg.IP = "10.0.0.2"
	}
	if cfg.Gateway == "" {
		cfg.Gateway = "10.0.0.1"
	}
	if len(cfg.DNS) == 0 {
		cfg.DNS = []string{"8.8.8.8", "1.1.1.1"}
	}
	if cfg.SocksAddr == "" {
		return nil, fmt.Errorf("SocksAddr is required")
	}

	tunDev := &TunDevice{
		Name:      cfg.Name,
		IP:        cfg.IP,
		Gateway:   cfg.Gateway,
		DNS:       cfg.DNS,
		MTU:       cfg.MTU,
		socksAddr: cfg.SocksAddr,
	}

	// Set platform-specific TUN device name
	if err := tunDev.createDevice(); err != nil {
		return nil, fmt.Errorf("failed to prepare TUN device: %w", err)
	}

	// Configure system routing first (before starting tun2socks)
	if err := tunDev.setupRouting(); err != nil {
		return nil, fmt.Errorf("failed to setup routing: %w", err)
	}

	log.Printf("[TUN] Device %s ready", tunDev.Name)
	return tunDev, nil
}

// Start begins tun2socks process
func (t *TunDevice) Start() error {
	t.mu.Lock()
	if t.running {
		t.mu.Unlock()
		return fmt.Errorf("TUN device already running")
	}

	// Find tun2socks binary (cross-platform)
	exePath, err := os.Executable()
	if err != nil {
		t.mu.Unlock()
		return fmt.Errorf("failed to get executable path: %w", err)
	}
	exeDir := filepath.Dir(exePath)
	tun2socksName := "tun2socks"
	if runtime.GOOS == "windows" {
		tun2socksName = "tun2socks.exe"
	}

	tun2socksPath := filepath.Join(exeDir, "bin", tun2socksName)
	if _, err := os.Stat(tun2socksPath); os.IsNotExist(err) {
		if wd, err := os.Getwd(); err == nil {
			tun2socksPath = filepath.Join(wd, "bin", tun2socksName)
			if _, err := os.Stat(tun2socksPath); os.IsNotExist(err) {
				t.mu.Unlock()
				return fmt.Errorf("%s not found in bin/ folder\nPlease download from: https://github.com/xjasonlyu/tun2socks/releases", tun2socksName)
			}
		}
	}

	log.Printf("[TUN] Using tun2socks binary: %s", tun2socksPath)

	args := []string{
		"-device", fmt.Sprintf("tun://%s", t.Name),
		"-proxy", fmt.Sprintf("socks5://%s", t.socksAddr),
		"-loglevel", "silent",
		"-tcp-no-delay",
		"-udp-timeout", "30s",
	}

	t.cmd = exec.Command(tun2socksPath, args...)
	t.cmd.Stdout = os.Stdout
	t.cmd.Stderr = os.Stderr

	if err := t.cmd.Start(); err != nil {
		t.mu.Unlock()
		return fmt.Errorf("failed to start tun2socks: %w\nNote: Requires admin/root privileges", err)
	}

	t.running = true
	log.Printf("[TUN] Started tun2socks: %s -> %s (PID: %d)", t.Name, t.socksAddr, t.cmd.Process.Pid)

	// processDone is closed when tun2socks exits (used for early-crash detection below)
	processDone := make(chan struct{})
	go func() {
		err := t.cmd.Wait()
		t.mu.Lock()
		wasRunning := t.running
		t.running = false
		t.mu.Unlock()
		close(processDone)
		if wasRunning {
			if err != nil {
				log.Printf("[TUN] tun2socks exited with error: %v", err)
			} else {
				log.Printf("[TUN] tun2socks exited normally")
			}
		}
	}()

	t.mu.Unlock()

	// Wait up to 1s for TUN device to initialize; bail out immediately on early crash
	select {
	case <-processDone:
		return fmt.Errorf("tun2socks exited immediately after start")
	case <-time.After(1 * time.Second):
		// still running — TUN device is ready
	}

	return nil
}

// Stop stops tun2socks process
func (t *TunDevice) Stop() {
	t.mu.Lock()
	defer t.mu.Unlock()

	if !t.running {
		return
	}

	t.running = false

	// Kill tun2socks process
	if t.cmd != nil && t.cmd.Process != nil {
		log.Printf("[TUN] Stopping tun2socks (PID: %d)", t.cmd.Process.Pid)
		if err := t.cmd.Process.Kill(); err != nil {
			log.Printf("[TUN] Failed to kill tun2socks: %v", err)
		}
		t.cmd = nil
	}

	log.Printf("[TUN] Stopped %s", t.Name)
}

// Close closes the TUN device and cleans up
func (t *TunDevice) Close() error {
	t.Stop()

	// Remove routing
	if err := t.removeRouting(); err != nil {
		log.Printf("[TUN] Warning: failed to remove routing: %v", err)
	}

	log.Printf("[TUN] Closed %s", t.Name)
	return nil
}

// IsRunning returns true if the TUN device is running
func (t *TunDevice) IsRunning() bool {
	t.mu.Lock()
	defer t.mu.Unlock()
	return t.running
}
