//go:build darwin

package tun

import (
	"fmt"
	"log"
	"os/exec"
)

// setupRouting configures macOS routing to direct traffic through TUN
func (t *TunDevice) setupRouting() error {
	log.Println("[TUN] Setting up macOS routing...")

	// Set interface IP address and destination (point-to-point)
	cmd := exec.Command("ifconfig", t.Name, t.IP, t.Gateway, "up")
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to configure interface: %w\nOutput: %s", err, output)
	}
	log.Printf("[TUN] Set interface %s IP to %s -> %s", t.Name, t.IP, t.Gateway)

	// Set MTU
	cmd = exec.Command("ifconfig", t.Name, "mtu", fmt.Sprintf("%d", t.MTU))
	if output, err := cmd.CombinedOutput(); err != nil {
		log.Printf("[TUN] Warning: failed to set MTU: %v (%s)", err, output)
	}

	// Save current default gateway
	defaultGW, err := t.getDefaultGateway()
	if err != nil {
		log.Printf("[TUN] Warning: could not detect default gateway: %v", err)
	} else {
		log.Printf("[TUN] Current default gateway: %s", defaultGW)
	}

	// Add default route through TUN
	cmd = exec.Command("route", "add", "default", t.Gateway)
	if output, err := cmd.CombinedOutput(); err != nil {
		log.Printf("[TUN] Warning: failed to add default route: %v (%s)", err, output)
	}
	log.Printf("[TUN] Added default route via %s", t.Gateway)

	// Set DNS servers using networksetup or scutil
	if err := t.setDNS(); err != nil {
		log.Printf("[TUN] Warning: failed to set DNS: %v", err)
	}

	log.Println("[TUN] macOS routing setup complete")
	return nil
}

// removeRouting removes macOS routing configuration
func (t *TunDevice) removeRouting() error {
	log.Println("[TUN] Removing macOS routing...")

	// Remove default route through TUN
	cmd := exec.Command("route", "delete", "default", t.Gateway)
	if output, err := cmd.CombinedOutput(); err != nil {
		log.Printf("[TUN] Warning: failed to remove route: %v (%s)", err, output)
	}

	// Bring interface down
	cmd = exec.Command("ifconfig", t.Name, "down")
	if output, err := cmd.CombinedOutput(); err != nil {
		log.Printf("[TUN] Warning: failed to bring down interface: %v (%s)", err, output)
	}

	log.Println("[TUN] macOS routing cleanup complete")
	return nil
}

// getDefaultGateway retrieves the current default gateway
func (t *TunDevice) getDefaultGateway() (string, error) {
	cmd := exec.Command("route", "-n", "get", "default")
	_, err := cmd.Output()
	if err != nil {
		return "", err
	}

	// Parse output to extract gateway
	// Format: "gateway: 192.168.1.1"
	return "", fmt.Errorf("not implemented")
}

// setDNS configures DNS servers on macOS
func (t *TunDevice) setDNS() error {
	// macOS: Use scutil to set DNS
	// This is complex and requires writing to dynamic store
	// For now, log DNS servers (user can set manually)
	log.Printf("[TUN] DNS servers to use: %v", t.DNS)
	log.Println("[TUN] Note: macOS DNS configuration may require manual setup in System Preferences")
	
	// TODO: Implement scutil-based DNS configuration
	// Or use networksetup command for specific interface
	
	return nil
}
