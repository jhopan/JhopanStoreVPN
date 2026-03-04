package proxy

// SetSystemProxy sets the OS-level HTTP proxy to the given address.
// ResetSystemProxy removes/disables the system proxy.
// These are implemented per-platform in proxy_windows.go, proxy_darwin.go, proxy_linux.go.
