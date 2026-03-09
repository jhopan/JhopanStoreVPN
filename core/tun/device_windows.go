//go:build windows

package tun

import (
	"log"
)

// createDevice prepares TUN device name for Windows
// Note: go-tun2socks will handle actual device creation
func (t *TunDevice) createDevice() error {
	// Set default name for Windows
	if t.Name == "" {
		t.Name = "tun0"
	}
	
	log.Printf("[TUN] Windows TUN device name: %s", t.Name)
	log.Println("[TUN] Note: Windows requires Wintun/TAP driver and admin privileges")
	return nil
}
