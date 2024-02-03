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
func BuildTopBar(cmdChan chan string, wpm *canvas.Text) *fyne.Container {
	wpm.Text = fmt.Sprintf("  WPM:%4d  ", 300)
	wpm.TextStyle.Bold = true

	incButton := widget.NewButtonWithIcon("", theme.ContentAddIcon(), func() {
		cmdChan <- "inc wpm"
	})
	decButton := widget.NewButtonWithIcon("", theme.ContentRemoveIcon(), func() {
		cmdChan <- "dec wpm"
	})
	resetButton := widget.NewButtonWithIcon("", theme.MediaReplayIcon(), func() {
		cmdChan <- "restart"
	})

	left := container.NewHBox(decButton, wpm, incButton, layout.NewSpacer())
	right := container.NewGridWithColumns(3, layout.NewSpacer(), layout.NewSpacer(), resetButton)

	return container.NewGridWithColumns(3, left, layout.NewSpacer(), right)
}

// BuildBottomBar creates the botton display and maps buttons to commands
func BuildBottomBar(cmdChan chan string, progress binding.Float) *fyne.Container {
	playButton := widget.NewButtonWithIcon("Play", theme.MediaPlayIcon(), func() {
		cmdChan <- "play"
	})
	fowardButton := widget.NewButtonWithIcon("", theme.MediaFastForwardIcon(), func() {
		cmdChan <- "inc"
	})
	rewindButton := widget.NewButtonWithIcon("", theme.MediaFastRewindIcon(), func() {
		cmdChan <- "dec"
	})

	progress.Set(0)
	progressBar := NewCustomProgressBar(progress, 8)

	return container.NewVBox(progressBar, container.NewGridWithColumns(3, rewindButton, playButton, fowardButton))
}

// BuildCenterBox builds main display where the Clipboard text will be displayed
func BuildCenterBox(left, center, right *canvas.Text) *fyne.Container {
	center.Text = " Ready "

	left.Text = " >>> "
	right.Text = " <<< "

	center.Alignment = fyne.TextAlignCenter
	center.TextStyle = fyne.TextStyle{Bold: true}
	center.TextSize = 32

	centerContainer := container.NewCenter(container.NewStack(center))

	right.Alignment = fyne.TextAlignLeading
	left.Alignment = fyne.TextAlignTrailing

	return container.NewGridWithColumns(3, left, centerContainer, right)
}
