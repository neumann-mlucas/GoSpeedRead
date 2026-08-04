package uielements

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type CustomProgressBar struct {
	widget.ProgressBar
	minHeight float32
}

// NewCustomProgressBar creates a fyne progress bar with a custom height
func NewCustomProgressBar(data binding.Float, minHeight float32) *CustomProgressBar {
	progressBar := &CustomProgressBar{minHeight: minHeight}
	progressBar.ExtendBaseWidget(progressBar)
	progressBar.TextFormatter = func() string { return "" }
	progressBar.Bind(data)
	return progressBar
}

// MinSize overrides original method
func (c *CustomProgressBar) MinSize() fyne.Size {
	minSize := c.ProgressBar.MinSize()
	if c.minHeight < minSize.Height {
		return fyne.NewSize(minSize.Width, c.minHeight)
	}
	return minSize
}

// BuildTopBar creates the top display and maps buttons to commands
func BuildTopBar(initialWPM int, wpm *canvas.Text, incWPM, decWPM, restart func()) *fyne.Container {
	wpm.Text = fmt.Sprintf("  WPM:%4d  ", initialWPM)
	wpm.TextStyle.Bold = true

	incButton := widget.NewButtonWithIcon("", theme.ContentAddIcon(), incWPM)
	decButton := widget.NewButtonWithIcon("", theme.ContentRemoveIcon(), decWPM)
	resetButton := widget.NewButtonWithIcon("", theme.MediaReplayIcon(), restart)

	left := container.NewHBox(decButton, wpm, incButton, layout.NewSpacer())
	right := container.NewGridWithColumns(3, layout.NewSpacer(), layout.NewSpacer(), resetButton)

	return container.NewGridWithColumns(3, left, layout.NewSpacer(), right)
}

// BuildBottomBar creates the bottom display and returns the container + play button
// so caller can toggle the play/pause icon.
func BuildBottomBar(progress binding.Float, position *canvas.Text, play, forward, back func()) (*fyne.Container, *widget.Button) {
	playButton := widget.NewButtonWithIcon("Play", theme.MediaPlayIcon(), play)
	forwardButton := widget.NewButtonWithIcon("", theme.MediaFastForwardIcon(), forward)
	rewindButton := widget.NewButtonWithIcon("", theme.MediaFastRewindIcon(), back)

	progress.Set(0)
	progressBar := NewCustomProgressBar(progress, 8)

	progressRow := container.NewBorder(nil, nil, nil, position, progressBar)
	buttons := container.NewGridWithColumns(3, rewindButton, playButton, forwardButton)
	return container.NewVBox(progressRow, buttons), playButton
}

// BuildCenterBox builds the main display. `center` can be any CanvasObject
// (e.g. a container of three text spans for ORP highlighting).
func BuildCenterBox(left *canvas.Text, center fyne.CanvasObject, right *canvas.Text) *fyne.Container {
	left.Alignment = fyne.TextAlignTrailing
	right.Alignment = fyne.TextAlignLeading
	return container.NewGridWithColumns(3, left, center, right)
}
