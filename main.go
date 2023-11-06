package main

import (
	"fmt"
	"image/color"
	"speedread/internal"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func RenderWords(words []internal.Word, content *widget.Label) {
	for _, w := range words {
		time.Sleep(time.Second)
		content.Text = w.Text
		content.TextStyle.Bold = !content.TextStyle.Bold
		content.Refresh()
	}
}

func BuildBottonBar(displayText *widget.Label) *fyne.Container {
	playButton := widget.NewButton("Play", func() {
		words := internal.GetClipBoard()
		RenderWords(words, displayText)
	})
	playButton.Icon = theme.MediaPlayIcon()

	fowardButton := widget.NewButton("", func() {
		fmt.Println("foward")
	})
	fowardButton.Icon = theme.MediaFastForwardIcon()

	rewindButton := widget.NewButton("", func() {
		fmt.Println("Rewind")
	})
	rewindButton.Icon = theme.MediaFastRewindIcon()

	return container.NewGridWithColumns(3, rewindButton, playButton, fowardButton)
}

// Clipboard content: In Fyne a Canvas is the area within which an
func main() {
	myApp := app.New()
	myWindow := myApp.NewWindow("SpeedRead")

	top := canvas.NewText("top bar", color.White)

	label := widget.NewLabel("")
	label.Alignment = fyne.TextAlignCenter
	label.TextStyle = fyne.TextStyle{Bold: true}

	bottonBar := BuildBottonBar(label)

	content := container.NewBorder(top, bottonBar, nil, nil, label)
	myWindow.SetContent(content)

	myWindow.Resize(fyne.NewSize(300, 100))
	myWindow.ShowAndRun()

}
