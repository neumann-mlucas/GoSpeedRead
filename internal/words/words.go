package words

import (
	"strings"
	"unicode"
)

const (
	Quotes            = "'\"‘’“”»«"
	MeanWordLenght    = 5.1
	MinimumWeight     = 800
	PunctuationWeight = 500
)

type Word struct {
	Text    string
	Weight  int
	inQuote bool
}

func (w Word) String() string {
	trimFn := func(r rune) bool {
		return unicode.IsSpace(r) || unicode.IsPunct(r) || strings.ContainsRune(Quotes, r)
	}
	word := strings.TrimFunc(w.Text, trimFn)

	if w.inQuote {
		return " “ " + word + " ” "
	}
	return word
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

// EndsWithAny checks if the input string 's' ends with any of the characters in 'chars'.
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
	weight := 1000 * float64(len(s)) / MeanWordLenght
	if weight < MinimumWeight {
		weight = MinimumWeight
	}
	if strings.ContainsFunc(s, unicode.IsPunct) {
		weight += PunctuationWeight
	}
	return int(weight)
}

// TrimPunct trims all punctuation characters but the quote characters
func TrimPunct(s string) string {
	fn := func(r rune) bool { return !strings.ContainsRune(Quotes, r) && unicode.IsPunct(r) }
	return strings.TrimFunc(s, fn)
}

// ParseWords splits a string into words and calculates their weights
func ParseWords(s string) []Word {
	words := []Word{}
	inQuote := false
	for _, w := range strings.Fields(s) {
		w := strings.TrimSpace(w)

		// strip punctuation to avoid corner case of: foo ("bar ..."
		if StartsWithAny(TrimPunct(w), Quotes) {
			inQuote = true
		}

		if len(w) > 0 {
			word := Word{Text: w, Weight: CalcWeight(w), inQuote: inQuote}
			words = append(words, word)
		}

		// strip punctuation to avoid corner case of: "... foo". bar...
		if EndsWithAny(TrimPunct(w), Quotes) {
			inQuote = false
		}
	}
	return words
}
