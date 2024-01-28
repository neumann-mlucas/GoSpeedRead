package main

import (
	// "fmt"
	"time"

	"internal/displaytext"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	// "fyne.io/fyne/v2/layout"
	// "fyne.io/fyne/v2/data/binding"
)

const WPM = 300

type SpeedRead struct {
	Top    *fyne.Container
	Center *fyne.Container
	Bottom *fyne.Container
	Window *fyne.Container

	labels map[string]*canvas.Text
	text   *displaytext.DisplayText
}

func main() {
	myApp := app.New()
	myWindow := myApp.NewWindow("SpeedRead")

	app := NewSpeedRead()
	content := app.Window

	myWindow.SetContent(content)
	myWindow.Resize(fyne.NewSize(600, 200))
	myWindow.ShowAndRun()
}

func NewSpeedRead() *SpeedRead {
	app := &SpeedRead{
		labels: make(map[string]*canvas.Text),
		text:   displaytext.New(WPM),
	}

	app.Top = app.BuildTopBar()
	app.Center = app.BuildCenterBox()
	app.Bottom = app.BuildBottomBar()

	app.Window = container.NewBorder(app.Top, app.Bottom, nil, nil, app.Center)

	return app
}

func (s *SpeedRead) newLabel(name string, value string) *canvas.Text {
	w := canvas.NewText(value, theme.ForegroundColor())
	s.labels[name] = w
	return w
}

func (s *SpeedRead) BuildTopBar() *fyne.Container {
	wpm := s.newLabel("WPM", "WPM:  l00")
	pct := s.newLabel("PCT", "PCT:   0%")
	pct.Alignment = fyne.TextAlignLeading

	return container.NewGridWithColumns(2, wpm, pct)
}

func (s *SpeedRead) BuildBottomBar() *fyne.Container {
	playButton := widget.NewButtonWithIcon("Play", theme.MediaPlayIcon(), func() {
		go s.Play()
	})
	fowardButton := widget.NewButtonWithIcon("", theme.MediaFastForwardIcon(), func() {
		s.text.IncIndex(+5)
	})
	rewindButton := widget.NewButtonWithIcon("", theme.MediaFastRewindIcon(), func() {
		s.text.IncIndex(-5)
	})

	return container.NewGridWithColumns(3, rewindButton, playButton, fowardButton)
}

func (s *SpeedRead) BuildCenterBox() *fyne.Container {
	centerLabel := s.newLabel("W_CURRENT", "PLAY")
	leftLabel := s.newLabel("W_PREVIOUS", "")
	rightLabel := s.newLabel("W_NEXT", "")

	centerLabel.Alignment = fyne.TextAlignCenter
	centerLabel.TextStyle = fyne.TextStyle{Bold: true}
	centerLabel.TextSize = 32

	centerContainer := container.NewCenter(container.NewStack(centerLabel))

	rightLabel.Alignment = fyne.TextAlignLeading
	leftLabel.Alignment = fyne.TextAlignTrailing

	return container.NewGridWithColumns(3, leftLabel, centerContainer, rightLabel)
}

func (s *SpeedRead) Play() {
	s.text.GetClipBoard()
	for !s.text.IsLastWord() {
		w, t := s.text.Step()
		s.labels["W_CURRENT"].Text = w
		s.labels["W_CURRENT"].Refresh()
		time.Sleep(t)
	}
	for !s.text.IsLastWord() {
	}
}
