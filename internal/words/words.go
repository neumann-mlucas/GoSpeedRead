package words

import (
	"strings"
	"unicode"
	"unicode/utf8"
)

const (
	Quotes            = "'\"‘’“”»«"
	MeanWordLenght    = 5.1
	MinimumWeight     = 800
	PunctuationWeight = 500
	ParagraphWeight   = 1500
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

// ParseWords splits a string into words and calculates their weights.
// Words at paragraph boundaries (\n\n) get an extra ParagraphWeight pause.
func ParseWords(s string) []Word {
	out := []Word{}
	paragraphs := strings.Split(s, "\n\n")
	for pi, p := range paragraphs {
		inQuote := false
		pieces := strings.Fields(p)
		for i, raw := range pieces {
			w := strings.TrimSpace(raw)
			trimmed := TrimPunct(w)
			first, _ := utf8.DecodeRuneInString(trimmed)
			last, _ := utf8.DecodeLastRuneInString(trimmed)
			if strings.ContainsRune(Quotes, first) {
				inQuote = true
			}
			if len(w) > 0 {
				weight := CalcWeight(w)
				if i == len(pieces)-1 && pi < len(paragraphs)-1 {
					weight += ParagraphWeight
				}
				out = append(out, Word{Text: w, Weight: weight, inQuote: inQuote})
			}
			if strings.ContainsRune(Quotes, last) {
				inQuote = false
			}
		}
	}
	return out
}
