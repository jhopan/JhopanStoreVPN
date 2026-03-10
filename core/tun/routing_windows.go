//go:build windows

package tun

import (
	"fmt"
	"log"
	"os/exec"
	"strings"
	"time"
)

// setupRouting configures Windows routing to direct traffic through TUN
// Note: tun2socks will create the TUN interface, we just configure routes
func (t *TunDevice) setupRouting() error {
	log.Println("[TUN] Setting up Windows routing...")

	// Note: tun2socks will handle TUN interface creation and IP configuration
	// We don't need to manually set interface IP here

	// Get default gateway interface (for restoration later if needed)
	defaultGW, err := t.getDefaultGateway()
	if err != nil {
		log.Printf("[TUN] Warning: could not detect default gateway: %v", err)
	} else {
		log.Printf("[TUN] Default gateway: %s", defaultGW)
	}

	log.Println("[TUN] Windows routing setup complete (routes will be added after TUN device is created)")
	return nil
}

// removeRouting removes Windows routing configuration
func (t *TunDevice) removeRouting() error {
	log.Println("[TUN] Removing Windows routing...")

	// Try to remove default route through TUN gateway
	// This ensures no zombie routes remain even if tun2socks didn't clean up properly
	cmd := exec.Command("route", "delete", "0.0.0.0", "mask", "0.0.0.0", t.Gateway)
	if output, err := cmd.CombinedOutput(); err != nil {
		// Route might not exist (which is fine - normal cleanup already happened)
		if !strings.Contains(string(output), "not found") && !strings.Contains(string(output), "element not found") {
			log.Printf("[TUN] Warning: failed to remove route: %v (%s)", err, output)
		} else {
			log.Println("[TUN] No TUN routes to remove (already cleaned)")
		}
	} else {
		log.Printf("[TUN] Removed default route via %s", t.Gateway)
	}

	log.Println("[TUN] Windows routing cleanup complete")
	return nil
}

// getDefaultGateway retrieves the current default gateway (before VPN)
func (t *TunDevice) getDefaultGateway() (string, error) {
	cmd := exec.Command("route", "print", "0.0.0.0")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", err
	}

	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		if strings.Contains(line, "0.0.0.0") {
			fields := strings.Fields(line)
			if len(fields) >= 3 {
				return fields[2], nil // Gateway IP
			}
		}
	}

	return "", fmt.Errorf("default gateway not found")
}

// addDefaultRoute adds a default route through TUN (if needed manually)
func (t *TunDevice) addDefaultRoute() error {
	// Wait for TUN interface to be ready
	time.Sleep(2 * time.Second)

	log.Println("[TUN] Adding default route through TUN...")
	
	// Add route for all traffic (0.0.0.0/0) through TUN
	cmd := exec.Command("route", "add", "0.0.0.0", "mask", "0.0.0.0", t.Gateway, "metric", "1")
	if output, err := cmd.CombinedOutput(); err != nil {
		// Route might already exist
		if !strings.Contains(string(output), "already exists") {
			log.Printf("[TUN] Warning: failed to add default route: %v (%s)", err, output)
		}
	} else {
		log.Printf("[TUN] Added default route via %s", t.Gateway)
	}

	return nil
}
