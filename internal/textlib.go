package internal

import (
	"fmt"
	"strings"
	"time"
	"unicode"

	"github.com/atotto/clipboard"
)

const (
	QUOTES = "'\"‘’“”»«"
)

type DisplayText struct {
	Words []Word
	Index int
	WPM   int
}

type Word struct {
	Text    string
	Weight  int
	inQuote bool
}

func NewDisplayText(wpm int) *DisplayText {
	rawText := GetClipBoard()
	text := ProcessClipBoardText(rawText)

	text = append(text, Word{})

	// inpChannel := make(chan string)
	// outChannel := make(chan string)

	return &DisplayText{
		Words: text,
		Index: 0,
		WPM:   wpm,
		// inp:   inpChannel,
		// out:   outChannel,
	}
}

func (t *DisplayText) GetClipBoard() {
	content, err := clipboard.ReadAll()
	if err != nil {
		fmt.Println("Error reading clipboard:", err)
		content = ""
	}

	words := ProcessClipBoardText(content)
	words = append(words, Word{})
	t.Words = words
}

func (t *DisplayText) Step(cmd string) Word {
	word := t.Words[t.Index]
	t.HandleCmd(cmd)
	time.Sleep(t.DisplayTime())
	t.IncIndex(+1)
	return word
}

func (t *DisplayText) DisplayTime() time.Duration {
	wps := float64(t.WPM) / 60.0
	currWord := t.Words[t.Index]
	seconds := float64(currWord.Weight) / wps
	return time.Duration(seconds) * time.Millisecond
}

func (t *DisplayText) IncIndex(inc int) {
	if inc > 0 {
		t.Index = min(t.Index+inc, len(t.Words))
	} else {
		t.Index = max(t.Index+inc, 0)
	}
}

func (t *DisplayText) End() bool {
	return t.Index >= len(t.Words)
}

func (t *DisplayText) Percentage() int {
	return 100 * t.Index / len(t.Words)
}

func (t *DisplayText) HandleCmd(cmd string) {
	switch cmd {
	case "inc":
		t.IncIndex(+5)
	case "dec":
		t.IncIndex(-5)
	case "restart":
		t.Index = 0
	}
}

// ProcessClipBoardText splits a string into words and calculates their weights
func ProcessClipBoardText(s string) []Word {
	words := []Word{}
	inQuote := false

	for _, w := range strings.Fields(s) {
		w := strings.TrimSpace(w)

		if StartsWithAny(w, QUOTES) {
			inQuote = true
		}

		if len(w) > 0 {
			word := Word{Text: w, Weight: CalcWeight(w), inQuote: inQuote}
			words = append(words, word)
		}

		// check if quotation starts and ends in this word. e.g 'foo'
		if EndsWithAny(w, QUOTES) {
			inQuote = false
		}
	}

	return words
}

// StartsWithAny checks if the input string 's' starts with any of the characters in 'chars'.
func StartsWithAny(s string, chars string) bool {
	for _, c := range chars {
		if strings.HasPrefix(s, string(c)) {
			return true
		}
	}
	return false
}

// HasQuote checks if the input string contains quotes
func EndsWithAny(s string, chars string) bool {
	for _, c := range chars {
		if strings.HasSuffix(s, string(c)) {
			return true
		}
	}
	return false
}

// CalcWeight calculates the weight of a word based on its length and presence of punctuation.
func CalcWeight(s string) int {
	// the mean length of an English word is 5.1 characters
	weight := 1000 * float64(len(s)) / 5.10
	if weight < 800 {
		weight = 800
	}
	if strings.ContainsFunc(s, unicode.IsPunct) {
		weight += 500
	}
	return int(weight)
}

func GetClipBoard() string {
	// Get the content of the system clipboard
	content, err := clipboard.ReadAll()
	if err != nil {
		fmt.Println("Error reading clipboard:", err)
		return ""
	}

	// Print the clipboard content
	fmt.Println("Clipboard content:", content)
	return content
}
