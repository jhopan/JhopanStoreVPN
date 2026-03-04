package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

// SettingsData holds all settings values.
type SettingsData struct {
	// Connection
	Path string
	SNI  string
	Host string

	// DNS
	DNS1 string
	DNS2 string

	// Behavior
	SystemTray    bool
	AutoReconnect bool
	AllowInsecure bool

	// Ping
	PingURL string
}

// DefaultSettings returns the default settings.
func DefaultSettings() SettingsData {
	return SettingsData{
		Path:          "/vless",
		SNI:           "",
		Host:          "biznet.vpnstore28.my.id",
		DNS1:          "8.8.8.8",
		DNS2:          "8.8.4.4",
		SystemTray:    true,
		AutoReconnect: true,
		AllowInsecure: true,
		PingURL:       "https://dns.google",
	}
}

// SettingsPage holds the settings UI and current values.
type SettingsPage struct {
	PathEntry  *widget.Entry
	SNIEntry   *widget.Entry
	HostEntry  *widget.Entry
	DNS1Entry  *widget.Entry
	DNS2Entry  *widget.Entry
	PingURLEntry *widget.Entry

	SystemTrayCheck    *widget.Check
	AutoReconnectCheck *widget.Check
	AllowInsecureCheck *widget.Check

	Container   *fyne.Container
	LogoResource fyne.Resource // set by caller
}

// NewSettingsPage creates the settings UI with default values.
func NewSettingsPage() *SettingsPage {
	sp := &SettingsPage{}
	defaults := DefaultSettings()

	// -- Connection section --
	sp.PathEntry = widget.NewEntry()
	sp.PathEntry.SetText(defaults.Path)
	sp.PathEntry.SetPlaceHolder("/vless")

	sp.SNIEntry = widget.NewEntry()
	sp.SNIEntry.SetText("biznet.vpnstore28.my.id")
	sp.SNIEntry.SetPlaceHolder("biznet.vpnstore28.my.id")

	sp.HostEntry = widget.NewEntry()
	sp.HostEntry.SetText(defaults.Host)
	sp.HostEntry.SetPlaceHolder("biznet.vpnstore28.my.id")

	// -- DNS section --
	sp.DNS1Entry = widget.NewEntry()
	sp.DNS1Entry.SetText(defaults.DNS1)
	sp.DNS1Entry.SetPlaceHolder("8.8.8.8")

	sp.DNS2Entry = widget.NewEntry()
	sp.DNS2Entry.SetText(defaults.DNS2)
	sp.DNS2Entry.SetPlaceHolder("8.8.4.4")

	// -- Ping section --
	sp.PingURLEntry = widget.NewEntry()
	sp.PingURLEntry.SetText(defaults.PingURL)
	sp.PingURLEntry.SetPlaceHolder("https://dns.google")

	// -- Behavior toggles --
	sp.SystemTrayCheck = widget.NewCheck("Minimize to System Tray", nil)
	sp.SystemTrayCheck.SetChecked(defaults.SystemTray)

	sp.AutoReconnectCheck = widget.NewCheck("Auto Reconnect", nil)
	sp.AutoReconnectCheck.SetChecked(defaults.AutoReconnect)

	sp.AllowInsecureCheck = widget.NewCheck("Allow Insecure TLS (skip verify)", nil)
	sp.AllowInsecureCheck.SetChecked(defaults.AllowInsecure)

	// Layout
	connHeader := widget.NewLabelWithStyle("— Connection —", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})
	connForm := container.New(layout.NewFormLayout(),
		widget.NewLabel("Path:"), sp.PathEntry,
		widget.NewLabel("SNI:"), sp.SNIEntry,
		widget.NewLabel("Host:"), sp.HostEntry,
	)

	dnsHeader := widget.NewLabelWithStyle("— DNS —", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})
	dnsForm := container.New(layout.NewFormLayout(),
		widget.NewLabel("DNS 1:"), sp.DNS1Entry,
		widget.NewLabel("DNS 2:"), sp.DNS2Entry,
	)

	pingHeader := widget.NewLabelWithStyle("— HTTP Ping —", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})
	pingForm := container.New(layout.NewFormLayout(),
		widget.NewLabel("Ping URL:"), sp.PingURLEntry,
	)

	behaviorHeader := widget.NewLabelWithStyle("— Behavior —", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})
	behaviorBox := container.NewVBox(
		sp.SystemTrayCheck,
		sp.AutoReconnectCheck,
		sp.AllowInsecureCheck,
	)

	sp.Container = container.NewVBox(
		connHeader,
		connForm,
		widget.NewSeparator(),
		dnsHeader,
		dnsForm,
		widget.NewSeparator(),
		pingHeader,
		pingForm,
		widget.NewSeparator(),
		behaviorHeader,
		behaviorBox,
		widget.NewSeparator(),
		layout.NewSpacer(),
	)

	return sp
}

// GetPath returns path value (with default).
func (sp *SettingsPage) GetPath() string {
	if sp.PathEntry.Text == "" {
		return "/vless"
	}
	return sp.PathEntry.Text
}

// GetSNI returns the SNI value. Defaults to biznet.vpnstore28.my.id if empty.
func (sp *SettingsPage) GetSNI() string {
	if sp.SNIEntry.Text == "" {
		return "biznet.vpnstore28.my.id"
	}
	return sp.SNIEntry.Text
}

// GetHost returns the host header value.
func (sp *SettingsPage) GetHost() string {
	if sp.HostEntry.Text == "" {
		return "biznet.vpnstore28.my.id"
	}
	return sp.HostEntry.Text
}

// GetDNS returns DNS server values.
func (sp *SettingsPage) GetDNS() (string, string) {
	d1 := sp.DNS1Entry.Text
	d2 := sp.DNS2Entry.Text
	if d1 == "" {
		d1 = "8.8.8.8"
	}
	if d2 == "" {
		d2 = "8.8.4.4"
	}
	return d1, d2
}

// GetPingURL returns the ping URL.
func (sp *SettingsPage) GetPingURL() string {
	if sp.PingURLEntry.Text == "" {
		return "https://dns.google"
	}
	return sp.PingURLEntry.Text
}

// IsSystemTray returns whether system tray is enabled.
func (sp *SettingsPage) IsSystemTray() bool {
	return sp.SystemTrayCheck.Checked
}

// IsAutoReconnect returns whether auto-reconnect is enabled.
func (sp *SettingsPage) IsAutoReconnect() bool {
	return sp.AutoReconnectCheck.Checked
}

// IsAllowInsecure returns whether to skip TLS certificate verification.
func (sp *SettingsPage) IsAllowInsecure() bool {
	return sp.AllowInsecureCheck.Checked
}
