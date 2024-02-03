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

	inp   chan string
	out   chan displaytext.DisplayState
	pause chan struct{}

	progress binding.Float
}

func main() {
	myApp := app.New()
	myWindow := myApp.NewWindow("SpeedRead")

	app := NewSpeedRead()
	content := app.Window

	myWindow.SetContent(content)
	myWindow.Resize(fyne.NewSize(Width, Height))

	myWindow.Canvas().SetOnTypedKey(app.HandleKey)

	go app.HandleCMD()
	go app.Run()

	myWindow.ShowAndRun()
}

// NewSpeedRead creates a new app
func NewSpeedRead() *SpeedRead {
	app := &SpeedRead{
		labels: make(map[string]*canvas.Text),
		text:   displaytext.New(WPM),

		inp:   make(chan string),
		out:   make(chan displaytext.DisplayState),
		pause: make(chan struct{}),

		progress: binding.NewFloat(),
	}

	app.Top = uielements.BuildTopBar(app.inp, app.newLabel("WPM", ""))

	app.Center = uielements.BuildCenterBox(app.newLabel("WordPrevious", ""), app.newLabel("WordCurrent", " READY "), app.newLabel("WordNext", ""))

	app.Bottom = uielements.BuildBottomBar(app.inp, app.progress)

	app.Window = container.NewBorder(app.Top, app.Bottom, nil, nil, app.Center)

	return app
}

// newLabel registers a new text label the the app `labels` map
func (sr *SpeedRead) newLabel(name string, value string) *canvas.Text {
	w := canvas.NewText(value, theme.ForegroundColor())
	sr.labels[name] = w
	return w
}

// HandleKey handles keyboard user input
func (sr *SpeedRead) HandleKey(k *fyne.KeyEvent) {
	switch k.Name {
	case fyne.KeySpace:
		sr.inp <- "play"
	case fyne.KeyR:
		sr.inp <- "restart"
	case fyne.KeyH, fyne.KeyLeft:
		sr.inp <- "dec"
	case fyne.KeyJ, fyne.KeyDown:
		sr.inp <- "dec wpm"
	case fyne.KeyK, fyne.KeyUp:
		sr.inp <- "inc wpm"
	case fyne.KeyL, fyne.KeyRight:
		sr.inp <- "inc"
	}
}

// HandleCMD alters the state of DisplayText base on a command
func (sr *SpeedRead) HandleCMD() {
	for cmd := range sr.inp {
		switch cmd {
		case "step":
			resp := sr.text.Step()
			sr.out <- resp
		case "restart":
			sr.text.Index = 0
		case "play":
			if sr.text.AtStartOrEnd() {
				sr.text.GetClipBoard()
			}
			sr.pause <- struct{}{}
		case "inc":
			sr.text.IncIndex(+5)
		case "dec":
			sr.text.IncIndex(-5)
		case "inc wpm":
			sr.text.WPM += 10
		case "dec wpm":
			sr.text.WPM -= 10
		}

	}
}

// Run updates the UI
func (sr *SpeedRead) Run() {
	for {
		<-sr.pause
		for !sr.text.IsLastWord() {
			select {
			case <-sr.pause:
				<-sr.pause

			default:
				sr.inp <- "step"
				state := <-sr.out

				sr.labels["WordCurrent"].Text = state.Text
				sr.labels["WordCurrent"].Refresh()

				sr.labels["WPM"].Text = fmt.Sprintf("  WPM:%4d  ", sr.text.WPM)
				sr.labels["WPM"].Refresh()

				sr.progress.Set(state.Prct)

				time.Sleep(state.Time)
			}
		}
	}
}
