package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"sync"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"

	"jhovpn/assets"
	"jhovpn/core/ping"
	"jhovpn/core/singleinstance"
	"jhovpn/core/tun"
	"jhovpn/core/vless"
	"jhovpn/core/xray"
	appUI "jhovpn/ui"
	jhovpnTheme "jhovpn/ui/theme"
	"jhovpn/ui/tray"
)

func main() {
	// Single instance check - prevent multiple instances
	instance, err := singleinstance.New("JhopanStoreVPN")
	if err != nil {
		// Another instance is already running
		log.Printf("Another instance is already running. Exiting.")
		// Note: Cannot show dialog here as Fyne app is not initialized yet
		// User will see nothing (which is expected - instance already running in tray)
		os.Exit(0)
	}
	defer instance.Release()

	// Create Fyne app
	a := app.NewWithID("com.jhopanstorevpn.app")
	a.Settings().SetTheme(&jhovpnTheme.DarkTheme{})
	a.SetIcon(assets.IconData)

	w := a.NewWindow("JhopanStoreVPN")
	w.Resize(fyne.NewSize(420, 520))
	w.SetFixedSize(true)

	// Load background image
	bgResource := assets.LoadBackground()

	// State
	var (
		xrayProc  *xray.Process
		tunDevice *tun.TunDevice
		pinger    *ping.Pinger
		connected bool
		connectMu sync.Mutex
	)

	// Forward-declare UI
	var mainPage *appUI.MainPage
	settingsPage := appUI.NewSettingsPage()

	// ---- Handlers ----

	doDisconnect := func() {
		connectMu.Lock()
		connected = false
		connectMu.Unlock()

		if pinger != nil {
			pinger.Stop()
		}
		if xrayProc != nil {
			xrayProc.Stop()
		}
		// Close TUN device
		if tunDevice != nil {
			tunDevice.Close()
			tunDevice = nil
		}

		if mainPage != nil {
			mainPage.SetDisconnected()
		}
		log.Println("[JhopanStoreVPN] Disconnected")
	}

	var doConnect func()

	onXrayCrash := func() {
		log.Println("[JhopanStoreVPN] Xray crashed, cleaning up TUN device")
		if pinger != nil {
			pinger.Stop()
		}
		// Close TUN device
		if tunDevice != nil {
			tunDevice.Close()
			tunDevice = nil
		}

		connectMu.Lock()
		wasConnected := connected
		connected = false
		connectMu.Unlock()

		if mainPage != nil {
			mainPage.SetDisconnected()
			mainPage.SetStatus("Xray crashed!")
		}

		// Auto-reconnect if enabled
		if wasConnected && settingsPage.IsAutoReconnect() {
			go func() {
				log.Println("[JhopanStoreVPN] Auto-reconnecting in 3s...")
				time.Sleep(3 * time.Second)
				connectMu.Lock()
				if !connected {
					connectMu.Unlock()
					doConnect()
				} else {
					connectMu.Unlock()
				}
			}()
		}
	}

	doConnect = func() {
		if mainPage == nil {
			return
		}

		address := mainPage.AddressEntry.Text
		uuid := mainPage.UUIDEntry.Text

		if address == "" || uuid == "" {
			dialog.ShowError(fmt.Errorf("Address and UUID are required"), w)
			return
		}

		domain, _, err := vless.SplitAddress(address)
		if err != nil {
			dialog.ShowError(fmt.Errorf("invalid address: %w", err), w)
			return
		}

		path := settingsPage.GetPath()
		sni := settingsPage.GetSNI()
		if sni == "" {
			sni = domain
		}
		host := settingsPage.GetHost()
		if host == "" {
			host = sni
		}
		dns1, dns2 := settingsPage.GetDNS()
		allowInsecure := settingsPage.IsAllowInsecure()

		vc := vless.Config{
			Address: address,
			UUID:    uuid,
			Path:    path,
			SNI:     sni,
			Host:    host,
		}

		mainPage.SetConnecting()

		// Run connection in background so UI doesn't freeze
		go func() {
			// Generate xray config
			configJSON, err := xray.GenerateConfig(vc, dns1, dns2, allowInsecure)
			if err != nil {
				mainPage.SetDisconnected()
				log.Printf("[JhopanStoreVPN] Config error: %v", err)
				return
			}

			// Start xray
			xrayProc = xray.NewProcess(onXrayCrash)
			if err := xrayProc.Start(configJSON); err != nil {
				mainPage.SetDisconnected()
				log.Printf("[JhopanStoreVPN] Xray start error: %v", err)
				return
			}

			// Wait for xray SOCKS5 port to be ready (up to 10 seconds)
			mainPage.SetStatus("Waiting for Xray...")
			portReady := false
			for i := 0; i < 40; i++ {
				conn, dialErr := net.DialTimeout("tcp", "127.0.0.1:10808", 250*time.Millisecond)
				if dialErr == nil {
					conn.Close()
					portReady = true
					break
				}
				time.Sleep(250 * time.Millisecond)
				if !xrayProc.IsRunning() {
					mainPage.SetDisconnected()
					mainPage.SetStatus("Xray exited unexpectedly")
					return
				}
			}
			if !portReady {
				xrayProc.Stop()
				mainPage.SetDisconnected()
				mainPage.SetStatus("Xray port timeout")
				return
			}
			log.Println("[JhopanStoreVPN] Xray SOCKS5 port 10808 is ready")

			// Create and start TUN device
			mainPage.SetStatus("Creating VPN tunnel...")
			tunCfg := tun.Config{
				Name:      "",  // Auto-assign (tun0, utun3, etc.)
				IP:        "10.0.0.2",
				Gateway:   "10.0.0.1",
				DNS:       []string{"8.8.8.8", "1.1.1.1"},
				MTU:       1400,
				SocksAddr: "127.0.0.1:10808",
			}
			
			var tunErr error
			tunDevice, tunErr = tun.NewTunDevice(tunCfg)
			if tunErr != nil {
				log.Printf("[JhopanStoreVPN] ERROR: Failed to create TUN device: %v", tunErr)
				xrayProc.Stop()
				mainPage.SetDisconnected()
				mainPage.SetStatus("TUN creation failed!")
				dialog.ShowError(fmt.Errorf("Failed to create VPN tunnel:\n%v\n\nNote: VPN mode requires admin/root privileges", tunErr), w)
				return
			}
			
			if err := tunDevice.Start(); err != nil {
				log.Printf("[JhopanStoreVPN] ERROR: Failed to start TUN device: %v", err)
				tunDevice.Close()
				tunDevice = nil
				xrayProc.Stop()
				mainPage.SetDisconnected()
				mainPage.SetStatus("TUN start failed!")
				dialog.ShowError(fmt.Errorf("Failed to start VPN tunnel:\n%v", err), w)
				return
			}
			
			log.Printf("[JhopanStoreVPN] TUN device started: %s", tunDevice.Name)

			connectMu.Lock()
			connected = true
			connectMu.Unlock()

			mainPage.SetConnected()
			log.Println("[JhopanStoreVPN] Connected to", address)

			// Start ping loop
			pingURL := settingsPage.GetPingURL()
			pinger = ping.NewPinger(pingURL, func(latency time.Duration, pingErr error) {
				if pingErr != nil {
					mainPage.SetPing("timeout")
				} else {
					mainPage.SetPing(fmt.Sprintf("%d ms", latency.Milliseconds()))
				}
			})
			pinger.Start()
		}()
	}

	// ---- Clipboard Import ----
	importClipboard := func() {
		clip := w.Clipboard()
		if clip == nil {
			dialog.ShowError(fmt.Errorf("clipboard not available"), w)
			return
		}
		content := clip.Content()
		if content == "" {
			dialog.ShowError(fmt.Errorf("clipboard is empty"), w)
			return
		}

		vc, err := vless.ParseURI(content)
		if err != nil {
			dialog.ShowError(fmt.Errorf("clipboard parse error:\n%w", err), w)
			return
		}

		mainPage.AddressEntry.SetText(vc.Address)
		mainPage.UUIDEntry.SetText(vc.UUID)
		settingsPage.PathEntry.SetText(vc.Path)
		settingsPage.SNIEntry.SetText(vc.SNI)
		settingsPage.HostEntry.SetText(vc.Host)

		dialog.ShowInformation("Imported", "VLESS config imported from clipboard.", w)
	}

	// ---- Open settings as popup ----
	openSettings := func() {
		// Logo header in settings
		var settingsHeader fyne.CanvasObject
		if bgResource != nil {
			logoImg := canvas.NewImageFromResource(bgResource)
			logoImg.FillMode = canvas.ImageFillContain
			logoImg.SetMinSize(fyne.NewSize(120, 70))
			settingsHeader = container.NewVBox(
				container.NewHBox(layout.NewSpacer(), logoImg, layout.NewSpacer()),
				widget.NewSeparator(),
			)
		} else {
			settingsHeader = container.NewVBox(
				widget.NewLabelWithStyle("JhopanStoreVPN", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
				widget.NewSeparator(),
			)
		}

		// Scrollable settings content
		settingsScroll := container.NewVScroll(settingsPage.Container)
		settingsScroll.SetMinSize(fyne.NewSize(370, 320))

		// Combine header + scroll
		settingsContent := container.NewBorder(settingsHeader, nil, nil, nil, settingsScroll)

		d := dialog.NewCustom("Settings", "Close", settingsContent, w)
		d.Resize(fyne.NewSize(400, 450))
		d.Show()
	}

	// ---- Build Main Page ----
	mainPage = appUI.NewMainPage(doConnect, doDisconnect, openSettings, importClipboard, bgResource)

	// ---- Load saved preferences ----
	prefs := a.Preferences()
	if v := prefs.String("address"); v != "" {
		mainPage.AddressEntry.SetText(v)
	}
	if v := prefs.String("uuid"); v != "" {
		mainPage.UUIDEntry.SetText(v)
	}
	if v := prefs.String("path"); v != "" {
		settingsPage.PathEntry.SetText(v)
	}
	if v := prefs.String("sni"); v != "" {
		settingsPage.SNIEntry.SetText(v)
	}
	if v := prefs.String("host"); v != "" {
		settingsPage.HostEntry.SetText(v)
	}
	if v := prefs.String("dns1"); v != "" {
		settingsPage.DNS1Entry.SetText(v)
	}
	if v := prefs.String("dns2"); v != "" {
		settingsPage.DNS2Entry.SetText(v)
	}
	if v := prefs.String("ping_url"); v != "" {
		settingsPage.PingURLEntry.SetText(v)
	}
	if prefs.StringWithFallback("system_tray_set", "no") == "yes" {
		settingsPage.SystemTrayCheck.SetChecked(prefs.Bool("system_tray"))
	}
	if prefs.StringWithFallback("auto_reconnect_set", "no") == "yes" {
		settingsPage.AutoReconnectCheck.SetChecked(prefs.Bool("auto_reconnect"))
	}
	if prefs.StringWithFallback("allow_insecure_set", "no") == "yes" {
		settingsPage.AllowInsecureCheck.SetChecked(prefs.Bool("allow_insecure"))
	}

	// ---- Auto-save on every field change ----
	savePrefs := func() {
		prefs.SetString("address", mainPage.AddressEntry.Text)
		prefs.SetString("uuid", mainPage.UUIDEntry.Text)
		prefs.SetString("path", settingsPage.PathEntry.Text)
		prefs.SetString("sni", settingsPage.SNIEntry.Text)
		prefs.SetString("host", settingsPage.HostEntry.Text)
		prefs.SetString("dns1", settingsPage.DNS1Entry.Text)
		prefs.SetString("dns2", settingsPage.DNS2Entry.Text)
		prefs.SetString("ping_url", settingsPage.PingURLEntry.Text)
		prefs.SetBool("system_tray", settingsPage.SystemTrayCheck.Checked)
		prefs.SetString("system_tray_set", "yes")
		prefs.SetBool("auto_reconnect", settingsPage.AutoReconnectCheck.Checked)
		prefs.SetString("auto_reconnect_set", "yes")
		prefs.SetBool("allow_insecure", settingsPage.AllowInsecureCheck.Checked)
		prefs.SetString("allow_insecure_set", "yes")
		
		// Force sync to disk to prevent data loss
		log.Println("[JhopanStoreVPN] Preferences saved")
	}

	mainPage.AddressEntry.OnChanged = func(_ string) { savePrefs() }
	mainPage.UUIDEntry.OnChanged = func(_ string) { savePrefs() }
	settingsPage.PathEntry.OnChanged = func(_ string) { savePrefs() }
	settingsPage.SNIEntry.OnChanged = func(_ string) { savePrefs() }
	settingsPage.HostEntry.OnChanged = func(_ string) { savePrefs() }
	settingsPage.DNS1Entry.OnChanged = func(_ string) { savePrefs() }
	settingsPage.DNS2Entry.OnChanged = func(_ string) { savePrefs() }
	settingsPage.PingURLEntry.OnChanged = func(_ string) { savePrefs() }
	settingsPage.SystemTrayCheck.OnChanged = func(_ bool) { savePrefs() }
	settingsPage.AutoReconnectCheck.OnChanged = func(_ bool) { savePrefs() }
	settingsPage.AllowInsecureCheck.OnChanged = func(_ bool) { savePrefs() }

	w.SetContent(mainPage.Container)

	// ---- System Tray ----
	visible := true
	tray.SetupTray(a, tray.Callbacks{
		OnShowHide: func() {
			if visible {
				w.Hide()
			} else {
				w.Show()
			}
			visible = !visible
		},
		OnConnect:    doConnect,
		OnDisconnect: doDisconnect,
		OnExit: func() {
			log.Println("[JhopanStoreVPN] Exit requested from tray, saving preferences...")
			savePrefs() // Save before quitting
			time.Sleep(100 * time.Millisecond) // Give time for preferences to flush
			doDisconnect()
			a.Quit()
		},
	})

	// Window close behavior depends on system tray toggle
	w.SetCloseIntercept(func() {
		log.Println("[JhopanStoreVPN] Window close intercepted, saving preferences...")
		savePrefs()
		time.Sleep(50 * time.Millisecond) // Give time for preferences to flush
		if settingsPage.IsSystemTray() {
			w.Hide()
			visible = false
		} else {
			doDisconnect()
			a.Quit()
		}
	})

	w.SetOnClosed(func() {
		savePrefs()
		doDisconnect()
	})

	_ = widget.NewLabel("v1.0.0")

	w.ShowAndRun()

	// Final cleanup
	doDisconnect()
}
