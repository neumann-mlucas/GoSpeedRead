package main

import (
	"fmt"
	// "image/color"
	"github.com/neumann-mlucas/GoSpeedRead/internal"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	// "fyne.io/fyne/v2/layout"
	// "fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

const WPM = 300

type SpeedRead struct {
	Top    *fyne.Container
	Center *fyne.Container
	Bottom *fyne.Container
	Window *fyne.Container

	labels map[string]*canvas.Text
	text   *internal.DisplayText
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
		text:   internal.NewDisplayText(WPM),
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
		s.Play()
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
	centerLabel := s.newCanvas("W_CURRENT", "PLAY")
	leftLabel := s.newCanvas("W_PREVIOUS", "")
	rightLabel := s.newCanvas("W_NEXT", "")

	centerLabel.Alignment = fyne.TextAlignCenter
	centerLabel.TextStyle = fyne.TextStyle{Bold: true}
	centerLabel.TextSize = 32

	// Use 'NewCenter' to place the label in the center of the window.
	centerContainer := container.NewCenter(container.NewStack(centerLabel))

	rightLabel.Alignment = fyne.TextAlignLeading
	leftLabel.Alignment = fyne.TextAlignTrailing

	return container.NewGridWithColumns(3, leftLabel, centerContainer, rightLabel)
}

func (s *SpeedRead) Play() {
	// TODO: update other bars
	// TODO: handle errors from step
	// TODO: pass arguments to NewDisplayText (WPM etc)
	go func() {
		s.text.GetClipBoard()
		for !s.text.End() {
			w := s.text.Step("nothing")
			fmt.Println(w, w.Weight)
			if w.Text == "" {
				break
			}
			s.canvas["W_CURRENT"].Text = w.Text
			s.canvas["W_CURRENT"].Refresh()
			s.labels["PCT"].SetText(fmt.Sprintf("PCT: %3d %%", s.text.Percentage()))
			s.labels["WPM"].SetText(fmt.Sprintf("WPM: %d", s.text.WPM))
		}
	}()

}
