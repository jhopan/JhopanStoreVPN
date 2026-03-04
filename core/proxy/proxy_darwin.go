//go:build darwin

package proxy

import (
	"fmt"
	"os/exec"
	"strings"
)

const proxyHost = "127.0.0.1"
const proxyPort = "10809"

// getNetworkServices returns active network service names.
func getNetworkServices() ([]string, error) {
	out, err := exec.Command("networksetup", "-listallnetworkservices").CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("failed to list network services: %w", err)
	}
	lines := strings.Split(strings.TrimSpace(string(out)), "\n")
	var services []string
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "An asterisk") {
			continue
		}
		services = append(services, line)
	}
	return services, nil
}

// SetSystemProxy enables HTTP/HTTPS proxy on all network services.
func SetSystemProxy() error {
	services, err := getNetworkServices()
	if err != nil {
		return err
	}
	for _, svc := range services {
		exec.Command("networksetup", "-setwebproxy", svc, proxyHost, proxyPort).Run()
		exec.Command("networksetup", "-setsecurewebproxy", svc, proxyHost, proxyPort).Run()
		exec.Command("networksetup", "-setsocksfirewallproxy", svc, proxyHost, "10808").Run()
	}
	return nil
}

// ResetSystemProxy disables HTTP/HTTPS/SOCKS proxy on all network services.
func ResetSystemProxy() error {
	services, err := getNetworkServices()
	if err != nil {
		return err
	}
	for _, svc := range services {
		exec.Command("networksetup", "-setwebproxystate", svc, "off").Run()
		exec.Command("networksetup", "-setsecurewebproxystate", svc, "off").Run()
		exec.Command("networksetup", "-setsocksfirewallproxystate", svc, "off").Run()
	}
	return nil
}
