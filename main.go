package main

import (
	"fmt"
	// "image/color"
	"speedread/internal"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	// "fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	// "fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type SpeedRead struct {
	Top    *fyne.Container
	Center *fyne.Container
	Bottom *fyne.Container

	labels map[string]*widget.Label
}

func NewSpeedRead() *SpeedRead {
	app := &SpeedRead{
		labels: make(map[string]*widget.Label),
	}
	app.Top = app.BuildTopBar()
	app.Center = app.BuildCenterBox()
	app.Bottom = app.BuildBottomBar()

	return app
}

func (s *SpeedRead) newLabel(name string) *widget.Label {
	w := widget.NewLabel("")
	s.labels[name] = w
	return w
}

func (s *SpeedRead) Play() {
	// TODO: update other bars
	// TODO: handle errors from step
	// TODO: pass arguments to NewDisplayText (WPM etc)
	go func() {
		text := internal.NewDisplayText(300)
		for {
			w := text.Step("nothing")
			fmt.Println(w)
			if w == "" {
				break
			}
			s.labels["W_CURRENT"].SetText(w)
			// s.labels["W_CURRENT"].Refresh()
		}
	}()

}

func (s *SpeedRead) BuildTopBar() *fyne.Container {
	wpm := s.newLabel("WPM")
	pct := s.newLabel("PCT")
	return container.NewGridWithColumns(2, wpm, pct)
}

func (s *SpeedRead) BuildBottomBar() *fyne.Container {
	playButton := widget.NewButtonWithIcon("Play", theme.MediaPlayIcon(), func() {
		s.Play()
	})
	fowardButton := widget.NewButtonWithIcon("", theme.MediaFastForwardIcon(), func() {
		fmt.Println("Foward")
	})
	rewindButton := widget.NewButtonWithIcon("", theme.MediaFastRewindIcon(), func() {
		fmt.Println("Rewind")
	})

	return container.NewGridWithColumns(3, rewindButton, playButton, fowardButton)
}

func (s *SpeedRead) BuildCenterBox() *fyne.Container {
	centerLabel := s.newLabel("W_CURRENT")
	left := s.newLabel("W_PREVIOUS")
	right := s.newLabel("W_NEXT")

	centerLabel.Alignment = fyne.TextAlignCenter
	centerLabel.TextStyle = fyne.TextStyle{Bold: true}
	// Use 'NewCenter' to place the label in the center of the window.
	centerContainer := container.NewCenter(container.NewStack(centerLabel))

	return container.NewGridWithColumns(3, left, centerContainer, right)
}

func main() {
	myApp := app.New()
	myWindow := myApp.NewWindow("SpeedRead")

	app := NewSpeedRead()
	content := container.NewBorder(app.Top, app.Bottom, nil, nil, app.Center)

	myWindow.SetContent(content)
	myWindow.Resize(fyne.NewSize(600, 200))
	myWindow.ShowAndRun()
}
