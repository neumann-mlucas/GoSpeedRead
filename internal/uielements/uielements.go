package uielements

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type CustomProgressBar struct {
	widget.ProgressBar
	minHeight float32
}

func NewCustomProgressBar(data binding.Float, minHeight float32) *CustomProgressBar {
	progressBar := &CustomProgressBar{minHeight: minHeight}
	progressBar.ExtendBaseWidget(progressBar)
	progressBar.TextFormatter = func() string { return "" }
	progressBar.Bind(data)
	return progressBar
}

func (c *CustomProgressBar) MinSize() fyne.Size {
	minSize := c.ProgressBar.MinSize()
	if c.minHeight < minSize.Height {
		return fyne.NewSize(minSize.Width, c.minHeight)
	}
	return minSize
}

func BuildTopBar(left, right *canvas.Text) *fyne.Container {
	left.Text = fmt.Sprintf("  WPM:%8d  ", 300)
	right.Text = fmt.Sprintf("  Progress:%8d%%  ", 0)

	left.Alignment = fyne.TextAlignLeading
	right.Alignment = fyne.TextAlignTrailing

	return container.NewGridWithColumns(2, left, right)
}

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
