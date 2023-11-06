package internal

import (
	"fmt"
	"strings"
	"unicode"

	"github.com/atotto/clipboard"
)

const (
	QUOTES = "'\"‘’“”»«"
)

type Word struct {
	Text    string
	Weight  int
	inQuote bool
}

// ProcessString splits a string into words and calculates their weights
func ProcessString(s string) []Word {
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
	weight := float64(100*len(s)) / 5.10
	// if weight < 100 {
	// 	weight = 100
	// }
	if strings.ContainsFunc(s, unicode.IsPunct) {
		weight += 50
	}
	return int(weight)
}

func GetClipBoard() []Word {
	// Get the content of the system clipboard
	content, err := clipboard.ReadAll()
	if err != nil {
		fmt.Println("Error reading clipboard:", err)
		return []Word{}
	}

	// Print the clipboard content
	fmt.Println("Clipboard content:", content)

	return ProcessString(content)
}
