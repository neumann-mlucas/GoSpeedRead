package main

import (
	"fmt"
	"time"

	"internal/displaytext"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	// "fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	// "fyne.io/fyne/v2/layout"
	// "fyne.io/fyne/v2/data/binding"
)

const (
	WPM    = 300
	Width  = 800
	Height = 200
)

type SpeedRead struct {
	Top    *fyne.Container
	Center *fyne.Container
	Bottom *fyne.Container
	Window *fyne.Container

	labels map[string]*canvas.Text
	text   *displaytext.DisplayText

	in     chan string
	out    chan displaytext.DisplayState
	signal chan struct{}
}

func main() {
	myApp := app.New()
	myWindow := myApp.NewWindow("SpeedRead")

	app := NewSpeedRead()
	content := app.Window

	myWindow.SetContent(content)
	myWindow.Resize(fyne.NewSize(Width, Height))

	go app.HandleCMD()
	go app.Play()

	myWindow.ShowAndRun()
}

func NewSpeedRead() *SpeedRead {
	app := &SpeedRead{
		labels: make(map[string]*canvas.Text),
		text:   displaytext.New(WPM),

		in:     make(chan string),
		out:    make(chan displaytext.DisplayState),
		signal: make(chan struct{}),
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
	wpm := s.newLabel("WPM", " WPM:  l00")
	pct := s.newLabel("PCT", " Progress:       0%")

	wpm.Alignment = fyne.TextAlignLeading
	pct.Alignment = fyne.TextAlignTrailing

	return container.NewGridWithColumns(2, wpm, pct)
}

func (s *SpeedRead) BuildBottomBar() *fyne.Container {
	playButton := widget.NewButtonWithIcon("Play", theme.MediaPlayIcon(), func() {
		s.in <- "play"
	})
	fowardButton := widget.NewButtonWithIcon("", theme.MediaFastForwardIcon(), func() {
		s.in <- "inc"
	})
	rewindButton := widget.NewButtonWithIcon("", theme.MediaFastRewindIcon(), func() {
		s.in <- "dec"
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

func (sr *SpeedRead) HandleCMD() {
	for cmd := range sr.in {
		switch cmd {
		case "step":
			resp := sr.text.Step()
			sr.out <- resp
		case "restart":
			sr.text.Index = 0
		case "play":
			sr.text.GetClipBoard()
			sr.signal <- struct{}{}
		case "inc":
			sr.text.IncIndex(+5)
		case "dec":
			sr.text.IncIndex(-5)
		case "inc wpm":
			sr.text.WPM += 10
		case "dec wpm":
			sr.text.WPM += 10
		}

	}
}

func (sr *SpeedRead) Play() {
	for {
		<-sr.signal
		for !sr.text.IsLastWord() {
			select {
			case <-sr.signal:
				<-sr.signal

			default:
				sr.in <- "step"
				state := <-sr.out

				sr.labels["W_CURRENT"].Text = state.Text
				sr.labels["W_CURRENT"].Refresh()

				sr.labels["PCT"].Text = fmt.Sprintf(" Progress:%8d%%", state.Prct)
				sr.labels["PCT"].Refresh()

				time.Sleep(state.Time)
			}
		}
	}
}
