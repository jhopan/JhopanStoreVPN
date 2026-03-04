//go:build !windows

package xray

import "os/exec"

// setProcAttr is a no-op on non-Windows platforms.
func setProcAttr(cmd *exec.Cmd) {
	// No special attributes needed on Unix/macOS
}
