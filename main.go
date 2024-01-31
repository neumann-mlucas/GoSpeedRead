package main

import (
	"fmt"
	"time"

	"internal/displaytext"
	"internal/uielements"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/theme"
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

	progress binding.Float
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

		progress: binding.NewFloat(),
	}

	app.Top = uielements.BuildTopBar(app.newLabel("WPM", ""), app.newLabel("PCT", ""))

	app.Center = uielements.BuildCenterBox(app.newLabel("WordPrevious", ""), app.newLabel("WordCurrent", " READY "), app.newLabel("WordNext", ""))

	app.Bottom = uielements.BuildBottomBar(app.in, app.bind)

	app.Window = container.NewBorder(app.Top, app.Bottom, nil, nil, app.Center)
	return app
}

func (s *SpeedRead) newLabel(name string, value string) *canvas.Text {
	w := canvas.NewText(value, theme.ForegroundColor())
	s.labels[name] = w
	return w
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

				sr.labels["WordCurrent"].Text = state.Text
				sr.labels["WordCurrent"].Refresh()

				sr.labels["PCT"].Text = fmt.Sprintf("  Progress:%8d%%  ", int(100*state.Prct))
				sr.labels["PCT"].Refresh()

				sr.labels["WPM"].Text = fmt.Sprintf("  WPM:%8d  ", sr.text.WPM)
				sr.labels["WPM"].Refresh()

				sr.progress.Set(state.Prct)

				time.Sleep(state.Time)
			}
		}
	}
}
