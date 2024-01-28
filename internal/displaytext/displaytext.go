package displaytext

import (
	"fmt"
	"github.com/atotto/clipboard"
	"internal/words"
	"time"
)

type DisplayText struct {
	Words []words.Word
	Index int
	WPM   int
}

// NewDisplayText Constructs a DisplayText object
func New(wpm int) *DisplayText {
	return &DisplayText{
		Words: []words.Word{},
		Index: 0,
		WPM:   wpm,
	}
}

// GetClipBoard populates the DisplayText object with the system clipboard content
func (t *DisplayText) GetClipBoard() {
	content, err := clipboard.ReadAll()
	if err != nil {
		fmt.Println("Error reading clipboard:", err)
		content = "ERROR_WHILE_READING_CLIPBOARD"
	}

	words := words.ProcessClipBoardText(content)
	// words = append(words, words.Word{}) // ???
	t.Words = words
}

// Step gets the next Word to be display and increments the object position
func (t *DisplayText) Step() (string, time.Duration) {
	defer t.IncIndex(+1)
	word := t.Words[t.Index]
	return word.String(), t.DisplayTime(word)
}

// DisplayTime calculates the display time in millisecond of a given Word
func (t *DisplayText) DisplayTime(word words.Word) time.Duration {
	wps := float64(t.WPM) / 60.0
	seconds := float64(word.Weight) / wps
	return time.Duration(seconds) * time.Millisecond
}

// HandleCmd given a cmd from the UI it updates the internal state of DisplayText
func (t *DisplayText) HandleCmd(cmd string) {
	switch cmd {
	case "inc":
		t.IncIndex(+5)
	case "dec":
		t.IncIndex(-5)
	case "restart":
		t.Index = 0
	case "inc WPM":
		t.WPM += 10
	case "dec WPM":
		t.WPM -= 10
	}
}

// IncIndex increments or decrements the index position by a given amount
func (t *DisplayText) IncIndex(inc int) {
	if inc > 0 {
		t.Index = min(t.Index+inc, len(t.Words))
	} else {
		t.Index = max(t.Index+inc, 0)
	}
}

// IsLastElement checks if the index is at the last element of a collection.
func (t *DisplayText) IsLastWord() bool {
	return t.Index >= len(t.Words)-1
}

// Percentage returns the relative position of the index
func (t *DisplayText) Percentage() int {
	return 100 * t.Index / len(t.Words)
}
