package main

import (
	"flag"
	"fmt"
	"image/color"
	"sync"
	"time"

	"internal/displaytext"
	"internal/uielements"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

const (
	previewLen = 16
	helpText   = " space play/pause · h/l ← → · j/k ± wpm · r reset · v reload clipboard · ? help "
)

var focalColor = color.NRGBA{R: 220, G: 60, B: 60, A: 255}

type SpeedRead struct {
	Window *fyne.Container

	wpmLabel      *canvas.Text
	positionLabel *canvas.Text
	helpLabel     *canvas.Text

	wordCurrentPrefix *canvas.Text
	wordCurrentFocal  *canvas.Text
	wordCurrentSuffix *canvas.Text
	wordNext          *canvas.Text
	wordPrevious      *canvas.Text

	playButton *widget.Button

	text     *displaytext.DisplayText
	progress binding.Float

	mu      sync.Mutex
	playing bool
}

func main() {
	width := flag.Int("width", 800, "The width of the window")
	height := flag.Int("height", 200, "The height of the window")
	wpm := flag.Int("WPM", 300, "Word per minute")
	themeName := flag.String("theme", "dark", "Theme: dark or light")
	// ponytail: font size derived from --height default; fyne lacks a clean
	// window-resize callback. Upgrade to a custom widget with Resize() override
	// if runtime scaling matters.
	fontSize := flag.Int("fontsize", 0, "Center font size (0 = auto from height)")
	flag.Parse()

	if *fontSize == 0 {
		*fontSize = int(float32(*height) * 0.16)
	}

	myApp := app.New()
	if *themeName == "light" {
		myApp.Settings().SetTheme(theme.LightTheme())
	} else {
		myApp.Settings().SetTheme(theme.DarkTheme())
	}

	myWindow := myApp.NewWindow("SpeedRead")

	sr := NewSpeedRead(*wpm, float32(*fontSize))
	myWindow.SetContent(sr.Window)
	myWindow.Resize(fyne.NewSize(float32(*width), float32(*height)))
	myWindow.Canvas().SetOnTypedKey(sr.HandleKey)
	myWindow.Canvas().SetOnTypedRune(sr.HandleRune)

	go sr.Run()

	myWindow.ShowAndRun()
}

// NewSpeedRead creates a new app
func NewSpeedRead(wpm int, fontSize float32) *SpeedRead {
	fg := theme.ForegroundColor()
	previewSize := fontSize * 0.5

	sr := &SpeedRead{
		text:              displaytext.New(wpm),
		progress:          binding.NewFloat(),
		wpmLabel:          newLabel(fmt.Sprintf("  WPM:%4d  ", wpm), 0, fg),
		positionLabel:     newLabel(" 0/0 ", 0, fg),
		helpLabel:         newLabel(helpText, 0, fg),
		wordCurrentPrefix: newLabel("", fontSize, fg),
		wordCurrentFocal:  newLabel("", fontSize, focalColor),
		wordCurrentSuffix: newLabel("", fontSize, fg),
		wordNext:          newLabel("", previewSize, fg),
		wordPrevious:      newLabel("", previewSize, fg),
	}

	boldMono := fyne.TextStyle{Bold: true, Monospace: true}
	sr.wpmLabel.TextStyle = boldMono
	sr.positionLabel.TextStyle = fyne.TextStyle{Monospace: true}
	sr.helpLabel.TextStyle = fyne.TextStyle{Italic: true, Monospace: true}
	sr.helpLabel.Hide()
	sr.wordCurrentPrefix.TextStyle = boldMono
	sr.wordCurrentFocal.TextStyle = boldMono
	sr.wordCurrentSuffix.TextStyle = boldMono
	sr.wordNext.TextStyle = fyne.TextStyle{Monospace: true}
	sr.wordPrevious.TextStyle = fyne.TextStyle{Monospace: true}

	sr.setCurrentWord(" READY ", -1)

	centerGroup := container.NewCenter(container.NewHBox(
		sr.wordCurrentPrefix, sr.wordCurrentFocal, sr.wordCurrentSuffix,
	))

	top := uielements.BuildTopBar(wpm, sr.wpmLabel, sr.incWPM, sr.decWPM, sr.restart)
	center := uielements.BuildCenterBox(sr.wordPrevious, centerGroup, sr.wordNext)
	bottom, playBtn := uielements.BuildBottomBar(sr.progress, sr.positionLabel, sr.togglePlay, sr.stepForward, sr.stepBack)
	sr.playButton = playBtn

	sr.Window = container.NewBorder(top, container.NewVBox(bottom, sr.helpLabel), nil, nil, center)
	return sr
}

func newLabel(text string, size float32, c color.Color) *canvas.Text {
	t := canvas.NewText(text, c)
	if size > 0 {
		t.TextSize = size
	}
	return t
}

// setCurrentWord updates the 3-span center display with an ORP-highlighted focal letter.
// focal < 0 disables highlight (all text in prefix).
func (sr *SpeedRead) setCurrentWord(text string, focal int) {
	runes := []rune(text)
	if focal < 0 || focal >= len(runes) {
		sr.wordCurrentPrefix.Text = text
		sr.wordCurrentFocal.Text = ""
		sr.wordCurrentSuffix.Text = ""
	} else {
		sr.wordCurrentPrefix.Text = string(runes[:focal])
		sr.wordCurrentFocal.Text = string(runes[focal : focal+1])
		sr.wordCurrentSuffix.Text = string(runes[focal+1:])
	}
	sr.wordCurrentPrefix.Refresh()
	sr.wordCurrentFocal.Refresh()
	sr.wordCurrentSuffix.Refresh()
}

// HandleKey handles keyboard user input for control keys
func (sr *SpeedRead) HandleKey(k *fyne.KeyEvent) {
	switch k.Name {
	case fyne.KeySpace:
		sr.togglePlay()
	case fyne.KeyR:
		sr.restart()
	case fyne.KeyV:
		sr.reloadClipboard()
	case fyne.KeyH, fyne.KeyLeft:
		sr.stepBack()
	case fyne.KeyJ, fyne.KeyDown:
		sr.decWPM()
	case fyne.KeyK, fyne.KeyUp:
		sr.incWPM()
	case fyne.KeyL, fyne.KeyRight:
		sr.stepForward()
	}
}

// HandleRune catches printable chars that don't map to a KeyName (e.g. ?)
func (sr *SpeedRead) HandleRune(r rune) {
	if r == '?' {
		sr.toggleHelp()
	}
}

func (sr *SpeedRead) togglePlay() {
	sr.mu.Lock()
	if !sr.playing && sr.text.IsFirstOrLastWord() {
		sr.text.GetClipBoard()
		sr.text.Index = 0
	}
	sr.playing = !sr.playing
	playing := sr.playing
	sr.mu.Unlock()

	if playing {
		sr.playButton.SetIcon(theme.MediaPauseIcon())
		sr.playButton.SetText("Pause")
	} else {
		sr.playButton.SetIcon(theme.MediaPlayIcon())
		sr.playButton.SetText("Play")
	}
}

func (sr *SpeedRead) reloadClipboard() {
	sr.mu.Lock()
	defer sr.mu.Unlock()
	sr.text.GetClipBoard()
	sr.text.Index = 0
}

func (sr *SpeedRead) toggleHelp() {
	if sr.helpLabel.Visible() {
		sr.helpLabel.Hide()
	} else {
		sr.helpLabel.Show()
	}
}

func (sr *SpeedRead) restart()     { sr.mu.Lock(); sr.text.Index = 0; sr.mu.Unlock() }
func (sr *SpeedRead) stepForward() { sr.mu.Lock(); sr.text.IncIndex(+5); sr.mu.Unlock() }
func (sr *SpeedRead) stepBack()    { sr.mu.Lock(); sr.text.IncIndex(-5); sr.mu.Unlock() }
func (sr *SpeedRead) incWPM()      { sr.mu.Lock(); sr.text.WPM += 10; sr.mu.Unlock() }
func (sr *SpeedRead) decWPM()      { sr.mu.Lock(); sr.text.WPM -= 10; sr.mu.Unlock() }

// Run drives the RSVP loop. Sleeps when paused or at end.
func (sr *SpeedRead) Run() {
	for {
		sr.mu.Lock()
		if !sr.playing || sr.text.IsLastWord() {
			sr.mu.Unlock()
			// ponytail: 50ms idle poll; upgrade to sync.Cond if first-frame latency after play matters
			time.Sleep(50 * time.Millisecond)
			continue
		}
		state := sr.text.Step()
		wpm := sr.text.WPM
		sr.mu.Unlock()

		sr.render(state, wpm)
		time.Sleep(state.Time)
	}
}

func (sr *SpeedRead) render(state displaytext.DisplayState, wpm int) {
	sr.setCurrentWord(state.Text, state.Focal)

	sr.wordNext.Text = headRunes(state.Next, previewLen)
	sr.wordNext.Refresh()

	sr.wordPrevious.Text = tailRunes(state.Prev, previewLen)
	sr.wordPrevious.Refresh()

	sr.wpmLabel.Text = fmt.Sprintf("  WPM:%4d  ", wpm)
	sr.wpmLabel.Refresh()

	sr.positionLabel.Text = fmt.Sprintf(" %d/%d ", state.Index+1, state.Total)
	sr.positionLabel.Refresh()

	sr.progress.Set(state.Prct)
}

// headRunes returns first n runes of s (UTF-8 safe).
func headRunes(s string, n int) string {
	r := []rune(s)
	if len(r) <= n {
		return s
	}
	return string(r[:n])
}

// tailRunes returns last n runes of s (UTF-8 safe).
func tailRunes(s string, n int) string {
	r := []rune(s)
	if len(r) <= n {
		return s
	}
	return string(r[len(r)-n:])
}
