//go:build windows

package proxy

import (
	"fmt"
	"log"
	"syscall"
	"time"
	"unsafe"

	"golang.org/x/sys/windows/registry"
)

const proxyAddr = "127.0.0.1:10809"

var (
	wininet                = syscall.NewLazyDLL("wininet.dll")
	procInternetSetOptionW = wininet.NewProc("InternetSetOptionW")
)

const (
	internetOptionSettingsChanged uintptr = 39
	internetOptionRefresh         uintptr = 37
)

// notifyProxyChange tells Windows to re-read proxy settings immediately.
func notifyProxyChange() {
	procInternetSetOptionW.Call(0, internetOptionSettingsChanged, 0, 0)
	procInternetSetOptionW.Call(0, internetOptionRefresh, 0, 0)
}

const regPath = `Software\Microsoft\Windows\CurrentVersion\Internet Settings`

// SetSystemProxy enables the Windows system proxy via direct registry API.
// It also disables "Automatically detect settings" which can override manual proxy.
func SetSystemProxy() error {
	log.Println("[Proxy] Opening registry key for write...")
	key, err := registry.OpenKey(registry.CURRENT_USER, regPath, registry.SET_VALUE|registry.QUERY_VALUE)
	if err != nil {
		return fmt.Errorf("failed to open registry key: %w", err)
	}
	defer key.Close()

	// Disable "Automatically detect settings" — this interferes with manual proxy
	log.Println("[Proxy] Disabling AutoDetect...")
	_ = key.SetDWordValue("AutoDetect", 0)

	log.Println("[Proxy] Setting ProxyEnable=1")
	if err := key.SetDWordValue("ProxyEnable", 1); err != nil {
		return fmt.Errorf("failed to set ProxyEnable: %w", err)
	}

	log.Printf("[Proxy] Setting ProxyServer=%s", proxyAddr)
	if err := key.SetStringValue("ProxyServer", proxyAddr); err != nil {
		return fmt.Errorf("failed to set ProxyServer: %w", err)
	}

	bypass := "localhost;127.*;10.*;192.168.*;<local>"
	log.Printf("[Proxy] Setting ProxyOverride=%s", bypass)
	if err := key.SetStringValue("ProxyOverride", bypass); err != nil {
		return fmt.Errorf("failed to set ProxyOverride: %w", err)
	}

	notifyProxyChange()

	// Verify the setting was applied
	time.Sleep(100 * time.Millisecond)
	val, _, verr := key.GetIntegerValue("ProxyEnable")
	if verr != nil || val != 1 {
		log.Printf("[Proxy] WARNING: Verification failed (val=%d err=%v), retrying...", val, verr)
		_ = key.SetDWordValue("ProxyEnable", 1)
		_ = key.SetStringValue("ProxyServer", proxyAddr)
		notifyProxyChange()
		time.Sleep(100 * time.Millisecond)
	}

	log.Println("[Proxy] System proxy ENABLED")
	return nil
}

// ResetSystemProxy disables the Windows system proxy and re-enables auto detect.
func ResetSystemProxy() error {
	key, err := registry.OpenKey(registry.CURRENT_USER, regPath, registry.SET_VALUE)
	if err != nil {
		log.Printf("[Proxy] Could not open registry key to reset: %v", err)
		return nil
	}
	defer key.Close()

	log.Println("[Proxy] Setting ProxyEnable=0")
	_ = key.SetDWordValue("ProxyEnable", 0)
	// Restore auto-detect
	_ = key.SetDWordValue("AutoDetect", 1)

	notifyProxyChange()
	log.Println("[Proxy] System proxy DISABLED")
	return nil
}

// Ensure unsafe import for InternetSetOption pointer parameters
var _ = unsafe.Pointer(nil)
