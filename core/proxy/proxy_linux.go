//go:build linux

package proxy

import (
	"os/exec"
)

const proxyURL = "http://127.0.0.1:10809"
const socksURL = "socks5://127.0.0.1:10808"

// SetSystemProxy sets GNOME system proxy settings.
// For KDE/other DEs, users should configure manually or use env vars.
func SetSystemProxy() error {
	commands := [][]string{
		{"gsettings", "set", "org.gnome.system.proxy", "mode", "manual"},
		{"gsettings", "set", "org.gnome.system.proxy.http", "host", "127.0.0.1"},
		{"gsettings", "set", "org.gnome.system.proxy.http", "port", "10809"},
		{"gsettings", "set", "org.gnome.system.proxy.https", "host", "127.0.0.1"},
		{"gsettings", "set", "org.gnome.system.proxy.https", "port", "10809"},
		{"gsettings", "set", "org.gnome.system.proxy.socks", "host", "127.0.0.1"},
		{"gsettings", "set", "org.gnome.system.proxy.socks", "port", "10808"},
	}
	for _, args := range commands {
		exec.Command(args[0], args[1:]...).Run() // best-effort
	}
	return nil
}

// ResetSystemProxy resets GNOME proxy to 'none'.
func ResetSystemProxy() error {
	exec.Command("gsettings", "set", "org.gnome.system.proxy", "mode", "none").Run()
	return nil
}
