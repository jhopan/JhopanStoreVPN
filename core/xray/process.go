package xray

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sync"
)

// Process manages the Xray-core subprocess lifecycle.
type Process struct {
	mu      sync.Mutex
	cmd     *exec.Cmd
	running bool
	onCrash func() // callback when xray crashes unexpectedly
}

// NewProcess creates a new Xray process manager.
func NewProcess(onCrash func()) *Process {
	return &Process{onCrash: onCrash}
}

// IsRunning returns whether xray is currently running.
func (p *Process) IsRunning() bool {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.running
}

// Start launches xray with the given config JSON bytes.
// It writes config to a temp file and runs xray-core.
func (p *Process) Start(configJSON []byte) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.running {
		return fmt.Errorf("xray is already running")
	}

	// Write config to temp file
	configDir, err := os.MkdirTemp("", "jhovpn")
	if err != nil {
		return fmt.Errorf("failed to create temp dir: %w", err)
	}
	configPath := filepath.Join(configDir, "config.json")
	if err := os.WriteFile(configPath, configJSON, 0644); err != nil {
		return fmt.Errorf("failed to write config: %w", err)
	}

	// Find xray binary
	xrayBin := findXrayBinary()
	if xrayBin == "" {
		return fmt.Errorf("xray binary not found. Place xray/xray.exe next to the application")
	}

	// Create log file for xray output (next to executable or cwd)
	logPath := filepath.Join(filepath.Dir(configPath), "xray.log")
	logFile, err := os.Create(logPath)
	if err != nil {
		logFile = nil
	}

	p.cmd = exec.Command(xrayBin, "run", "-c", configPath)
	setProcAttr(p.cmd)
	if logFile != nil {
		p.cmd.Stdout = logFile
		p.cmd.Stderr = logFile
	}

	if err := p.cmd.Start(); err != nil {
		if logFile != nil {
			logFile.Close()
		}
		return fmt.Errorf("failed to start xray: %w", err)
	}

	p.running = true

	// Monitor process in background
	go func() {
		err := p.cmd.Wait()
		if logFile != nil {
			logFile.Close()
		}
		p.mu.Lock()
		wasRunning := p.running
		p.running = false
		p.mu.Unlock()

		// Clean up temp config
		os.RemoveAll(configDir)

		// If it was supposed to be running (crash), invoke callback
		if wasRunning && err != nil && p.onCrash != nil {
			p.onCrash()
		}
	}()

	return nil
}

// Stop kills the xray process.
func (p *Process) Stop() error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if !p.running || p.cmd == nil || p.cmd.Process == nil {
		p.running = false
		return nil
	}

	p.running = false

	if runtime.GOOS == "windows" {
		return p.cmd.Process.Kill()
	}
	return p.cmd.Process.Signal(os.Interrupt)
}

// findXrayBinary looks for xray binary next to the executable, in PATH, or in current dir.
func findXrayBinary() string {
	names := []string{"xray"}
	if runtime.GOOS == "windows" {
		names = []string{"xray.exe"}
	}

	// Check next to executable
	if exePath, err := os.Executable(); err == nil {
		dir := filepath.Dir(exePath)
		for _, name := range names {
			p := filepath.Join(dir, name)
			if _, err := os.Stat(p); err == nil {
				return p
			}
		}
	}

	// Check current working directory
	if cwd, err := os.Getwd(); err == nil {
		for _, name := range names {
			p := filepath.Join(cwd, name)
			if _, err := os.Stat(p); err == nil {
				return p
			}
		}
	}

	// Check PATH
	for _, name := range names {
		if p, err := exec.LookPath(name); err == nil {
			return p
		}
	}

	return ""
}
