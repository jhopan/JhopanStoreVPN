package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// MainPage holds the main connection UI - vertical, clean layout.
type MainPage struct {
	AddressEntry *widget.Entry
	UUIDEntry    *widget.Entry

	ConnectBtn    *widget.Button
	DisconnectBtn *widget.Button

	StatusLabel *widget.Label
	PingLabel   *widget.Label

	Container fyne.CanvasObject
}

// NewMainPage creates the main vertical UI.
func NewMainPage(onConnect, onDisconnect, onSettings, onClipboard func(), logoResource fyne.Resource) *MainPage {
	mp := &MainPage{}

	mp.AddressEntry = widget.NewEntry()
	mp.AddressEntry.SetPlaceHolder("example.com:443")

	mp.UUIDEntry = widget.NewEntry()
	mp.UUIDEntry.SetPlaceHolder("xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx")

	mp.ConnectBtn = widget.NewButton("  CONNECT  ", onConnect)
	mp.ConnectBtn.Importance = widget.HighImportance

	mp.DisconnectBtn = widget.NewButton("DISCONNECT", onDisconnect)
	mp.DisconnectBtn.Importance = widget.DangerImportance
	mp.DisconnectBtn.Disable()

	mp.StatusLabel = widget.NewLabel("Disconnected")
	mp.StatusLabel.Alignment = fyne.TextAlignCenter

	mp.PingLabel = widget.NewLabel("Ping: -")
	mp.PingLabel.Alignment = fyne.TextAlignCenter

	// Top-right toolbar: clipboard icon + gear icon
	clipboardBtn := widget.NewButtonWithIcon("", theme.ContentPasteIcon(), onClipboard)
	clipboardBtn.Importance = widget.LowImportance

	settingsBtn := widget.NewButtonWithIcon("", theme.SettingsIcon(), onSettings)
	settingsBtn.Importance = widget.LowImportance

	// Title
	title := widget.NewLabelWithStyle("JhopanStoreVPN", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})

	// Top bar: title centered, clipboard + gear at right
	topBar := container.NewHBox(
		layout.NewSpacer(),
		title,
		layout.NewSpacer(),
		clipboardBtn,
		settingsBtn,
	)

	// Logo/banner image at top
	var logoBox fyne.CanvasObject
	if logoResource != nil {
		logoImage := canvas.NewImageFromResource(logoResource)
		logoImage.FillMode = canvas.ImageFillContain
		logoImage.SetMinSize(fyne.NewSize(180, 110))
		logoBox = container.NewHBox(layout.NewSpacer(), logoImage, layout.NewSpacer())
	} else {
		logoBox = widget.NewLabel("")
	}

	// Form - use padded container so entries fill available width
	addressLabel := widget.NewLabelWithStyle("Address", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
	uuidLabel := widget.NewLabelWithStyle("UUID", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})

	// Pad form content with margins so entries don't touch edges
	formContent := container.NewVBox(
		addressLabel,
		mp.AddressEntry,
		uuidLabel,
		mp.UUIDEntry,
	)
	paddedForm := container.New(layout.NewPaddedLayout(), formContent)

	// Buttons row
	buttons := container.NewHBox(
		layout.NewSpacer(),
		mp.ConnectBtn,
		mp.DisconnectBtn,
		layout.NewSpacer(),
	)

	// Status bar at bottom
	statusBar := container.NewVBox(
		widget.NewSeparator(),
		mp.StatusLabel,
		mp.PingLabel,
	)

	// Main content with border: top=header, bottom=status, center=scrollable content
	centerContent := container.NewVBox(
		logoBox,
		paddedForm,
		layout.NewSpacer(),
		buttons,
	)

	mp.Container = container.NewBorder(
		container.NewVBox(topBar, widget.NewSeparator()), // top
		statusBar,    // bottom
		nil,          // left
		nil,          // right
		centerContent, // center (fills remaining space)
	)

	return mp
}

// SetConnected updates UI for connected state.
func (mp *MainPage) SetConnected() {
	mp.ConnectBtn.Disable()
	mp.DisconnectBtn.Enable()
	mp.StatusLabel.SetText("Connected")
}

// SetDisconnected updates UI for disconnected state.
func (mp *MainPage) SetDisconnected() {
	mp.ConnectBtn.Enable()
	mp.DisconnectBtn.Disable()
	mp.StatusLabel.SetText("Disconnected")
	mp.PingLabel.SetText("Ping: -")
}

// SetConnecting updates UI for connecting state.
func (mp *MainPage) SetConnecting() {
	mp.ConnectBtn.Disable()
	mp.DisconnectBtn.Disable()
	mp.StatusLabel.SetText("Connecting...")
	mp.PingLabel.SetText("Ping: -")
}

// SetStatus sets a status message.
func (mp *MainPage) SetStatus(text string) {
	mp.StatusLabel.SetText(text)
}

// SetPing sets the ping display.
func (mp *MainPage) SetPing(text string) {
	mp.PingLabel.SetText("Ping: " + text)
}
