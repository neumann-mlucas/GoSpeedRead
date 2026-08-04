package displaytext

import (
	"fmt"
	"github.com/atotto/clipboard"
	"internal/words"
	"strings"
	"time"
)

type DisplayText struct {
	Words []words.Word
	Index int
	WPM   int
}

type DisplayState struct {
	Text  string
	Focal int // rune index of ORP (Optimal Recognition Point) letter within Text
	Time  time.Duration
	Prct  float64
	Next  string
	Prev  string
	Index int
	Total int
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
	words := words.ParseWords(content)
	t.Words = words
}

// Step gets the next Word to be display and increments the index position
func (t *DisplayText) Step() DisplayState {
	defer t.IncIndex(+1)
	word := t.Words[t.Index]
	text := word.String()
	// ponytail: ORP index computed on rendered text incl. quote padding; may misalign
	// focal letter on quoted words. Compute against Word.Text raw core if it matters.
	return DisplayState{
		Text:  text,
		Focal: ORPIndex(len([]rune(text))),
		Time:  t.DisplayTime(word),
		Prct:  t.Percentage(),
		Next:  t.GetNextWords(5),
		Prev:  t.GetPreviousWords(5),
		Index: t.Index,
		Total: len(t.Words),
	}
}

// ORPIndex returns the Optimal Recognition Point rune index for a word of length n.
// Standard Spritz-style table: focus letter shifts right as words grow.
func ORPIndex(n int) int {
	switch {
	case n <= 1:
		return 0
	case n <= 5:
		return 1
	case n <= 9:
		return 2
	case n <= 13:
		return 3
	default:
		return 4
	}
}

// DisplayTime calculates the display time in millisecond of a given Word
func (t *DisplayText) DisplayTime(word words.Word) time.Duration {
	wps := float64(t.WPM) / 60.0
	ms := float64(word.Weight) / wps
	return time.Duration(ms) * time.Millisecond
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
	return t.Index >= len(t.Words)
}

// IsFirstOrLastWord returns true if Index it at the first or last position
func (t *DisplayText) IsFirstOrLastWord() bool {
	return len(t.Words) == 0 || t.Index == 0 || t.Index >= len(t.Words)-1
}

// Percentage returns the relative position of the index
func (t *DisplayText) Percentage() float64 {
	return float64(t.Index) / float64(len(t.Words))
}

// GetPreviousWords returns a string with N words come before the current word
func (t *DisplayText) GetPreviousWords(i int) string {
	if t.Index <= 0 {
		return ""
	}
	var words []string
	for _, w := range t.Words[max(0, t.Index-i):t.Index] {
		words = append(words, w.Text)

	}
	return strings.Join(words, " ")
}

// GetNextWords returns a string with N words come after the current word
func (t *DisplayText) GetNextWords(i int) string {
	if t.IsLastWord() {
		return ""
	}
	var words []string
	for _, w := range t.Words[t.Index+1 : min(len(t.Words), t.Index+i+1)] {
		words = append(words, w.Text)

	}
	return strings.Join(words, " ")
}
