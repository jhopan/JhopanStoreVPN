//go:build linux

package tun

import (
	"log"
)

// createDevice prepares TUN device name for Linux
// Note: go-tun2socks will handle actual device creation
func (t *TunDevice) createDevice() error {
	// Set default name for Linux
	if t.Name == "" {
		t.Name = "tun0"
	}
	
	log.Printf("[TUN] Linux TUN device name: %s", t.Name)
	log.Println("[TUN] Note: Linux requires root/sudo or CAP_NET_ADMIN")
	return nil
}
