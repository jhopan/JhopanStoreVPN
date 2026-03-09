//go:build linux

package tun

import (
	"fmt"
	"log"
	"os/exec"
)

// setupRouting configures Linux routing to direct traffic through TUN
func (t *TunDevice) setupRouting() error {
	log.Println("[TUN] Setting up Linux routing...")

	// Bring interface up
	cmd := exec.Command("ip", "link", "set", "dev", t.Name, "up")
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to bring up interface: %w\nOutput: %s", err, output)
	}

	// Set MTU
	cmd = exec.Command("ip", "link", "set", "dev", t.Name, "mtu", fmt.Sprintf("%d", t.MTU))
	if output, err := cmd.CombinedOutput(); err != nil {
		log.Printf("[TUN] Warning: failed to set MTU: %v (%s)", err, output)
	}

	// Set interface IP address
	cmd = exec.Command("ip", "addr", "add", t.IP+"/24", "dev", t.Name)
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to set interface IP: %w\nOutput: %s", err, output)
	}
	log.Printf("[TUN] Set interface %s IP to %s", t.Name, t.IP)

	// Save current default gateway (for restoration)
	defaultGW, err := t.getDefaultGateway()
	if err != nil {
		log.Printf("[TUN] Warning: could not detect default gateway: %v", err)
	} else {
		log.Printf("[TUN] Current default gateway: %s", defaultGW)
	}

	// Add default route through TUN
	cmd = exec.Command("ip", "route", "add", "default", "via", t.Gateway, "dev", t.Name, "metric", "100")
	if output, err := cmd.CombinedOutput(); err != nil {
		log.Printf("[TUN] Warning: failed to add default route: %v (%s)", err, output)
	}
	log.Printf("[TUN] Added default route via %s", t.Gateway)

	// Configure DNS using resolvconf or systemd-resolved
	if err := t.setDNS(); err != nil {
		log.Printf("[TUN] Warning: failed to set DNS: %v", err)
	}

	log.Println("[TUN] Linux routing setup complete")
	return nil
}

// removeRouting removes Linux routing configuration
func (t *TunDevice) removeRouting() error {
	log.Println("[TUN] Removing Linux routing...")

	// Remove default route through TUN
	cmd := exec.Command("ip", "route", "del", "default", "via", t.Gateway, "dev", t.Name)
	if output, err := cmd.CombinedOutput(); err != nil {
		log.Printf("[TUN] Warning: failed to remove route: %v (%s)", err, output)
	}

	// Bring interface down
	cmd = exec.Command("ip", "link", "set", "dev", t.Name, "down")
	if output, err := cmd.CombinedOutput(); err != nil {
		log.Printf("[TUN] Warning: failed to bring down interface: %v (%s)", err, output)
	}

	log.Println("[TUN] Linux routing cleanup complete")
	return nil
}

// getDefaultGateway retrieves the current default gateway
func (t *TunDevice) getDefaultGateway() (string, error) {
	cmd := exec.Command("ip", "route", "show", "default")
	_, err := cmd.Output()
	if err != nil {
		return "", err
	}

	// Parse output: "default via 192.168.1.1 dev eth0 ..."
	// Extract gateway IP
	return "", fmt.Errorf("not implemented")
}

// setDNS configures DNS servers on Linux
func (t *TunDevice) setDNS() error {
	// Try systemd-resolved first
	for _, dns := range t.DNS {
		cmd := exec.Command("resolvectl", "dns", t.Name, dns)
		if err := cmd.Run(); err == nil {
			log.Printf("[TUN] Set DNS %s via systemd-resolved", dns)
			continue
		}
		
		// Fallback: try resolvconf
		cmd = exec.Command("resolvconf", "-a", t.Name)
		stdin, err := cmd.StdinPipe()
		if err != nil {
			continue
		}
		
		if err := cmd.Start(); err != nil {
			stdin.Close()
			continue
		}
		
		fmt.Fprintf(stdin, "nameserver %s\n", dns)
		stdin.Close()
		
		if err := cmd.Wait(); err == nil {
			log.Printf("[TUN] Set DNS %s via resolvconf", dns)
		}
	}
	
	return nil
}
