package main

import (
	"fmt"
	"log"
	"net"
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
	"jhovpn/core/proxy"
	"jhovpn/core/vless"
	"jhovpn/core/xray"
	appUI "jhovpn/ui"
	jhovpnTheme "jhovpn/ui/theme"
	"jhovpn/ui/tray"
)

func main() {
	// Safety: reset proxy on startup
	proxy.ResetSystemProxy()

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
		xrayProc    *xray.Process
		pinger      *ping.Pinger
		connected   bool
		connectMu   sync.Mutex
		autoReconWg sync.WaitGroup
		stopAutoRec chan struct{}
	)

	// Forward-declare UI
	var mainPage *appUI.MainPage
	settingsPage := appUI.NewSettingsPage()

	// ---- Handlers ----

	doDisconnect := func() {
		connectMu.Lock()
		connected = false
		// Stop auto-reconnect
		if stopAutoRec != nil {
			close(stopAutoRec)
			stopAutoRec = nil
		}
		connectMu.Unlock()

		if pinger != nil {
			pinger.Stop()
		}
		if xrayProc != nil {
			xrayProc.Stop()
		}
		proxy.ResetSystemProxy()

		if mainPage != nil {
			mainPage.SetDisconnected()
		}
		log.Println("[JhopanStoreVPN] Disconnected")
	}

	var doConnect func()

	onXrayCrash := func() {
		log.Println("[JhopanStoreVPN] Xray crashed, resetting proxy")
		if pinger != nil {
			pinger.Stop()
		}
		proxy.ResetSystemProxy()

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

			// Wait for xray HTTP proxy port to be ready (up to 10 seconds)
			mainPage.SetStatus("Waiting for Xray...")
			portReady := false
			for i := 0; i < 40; i++ {
				conn, dialErr := net.DialTimeout("tcp", "127.0.0.1:10809", 250*time.Millisecond)
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
			log.Println("[JhopanStoreVPN] Xray port 10809 is ready")

			// Set system proxy
			mainPage.SetStatus("Setting proxy...")
			if err := proxy.SetSystemProxy(); err != nil {
				log.Printf("[JhopanStoreVPN] WARNING: failed to set system proxy: %v", err)
				mainPage.SetStatus("Proxy set failed!")
			}

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

			// Auto-reconnect monitor
			if settingsPage.IsAutoReconnect() {
				connectMu.Lock()
				stopAutoRec = make(chan struct{})
				ch := stopAutoRec
				connectMu.Unlock()

				autoReconWg.Add(1)
				go func() {
					defer autoReconWg.Done()
					ticker := time.NewTicker(10 * time.Second)
					defer ticker.Stop()
					for {
						select {
						case <-ch:
							return
						case <-ticker.C:
							if xrayProc != nil && !xrayProc.IsRunning() {
								connectMu.Lock()
								if connected {
									connected = false
									connectMu.Unlock()
									log.Println("[JhopanStoreVPN] Detected xray stopped, reconnecting...")
									proxy.ResetSystemProxy()
									mainPage.SetStatus("Reconnecting...")
									time.Sleep(2 * time.Second)
									doConnect()
									return
								}
								connectMu.Unlock()
							}
						}
					}
				}()
			}
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
			doDisconnect()
			a.Quit()
		},
	})

	// Window close behavior depends on system tray toggle
	w.SetCloseIntercept(func() {
		savePrefs()
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
