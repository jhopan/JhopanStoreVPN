//go:build darwin

package tun

import (
	"log"
)

// createDevice prepares TUN device name for macOS
// Note: go-tun2socks will handle actual device creation
func (t *TunDevice) createDevice() error {
	// Set default name for macOS (utun interfaces)
	if t.Name == "" {
		t.Name = "utun3" // Start from utun3 to avoid conflicts
	}
	
	log.Printf("[TUN] macOS TUN device name: %s", t.Name)
	log.Println("[TUN] Note: macOS requires root/sudo privileges")
	return nil
}
