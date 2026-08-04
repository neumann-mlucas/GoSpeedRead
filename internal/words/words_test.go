package words

import (
	"reflect"
	"testing"
)

func TestCalcWeight(t *testing.T) {
	testCases := []struct {
		input    string
		expected int
	}{
		{"hello ", 1176},                     // length impacts weight
		{"world!", 1176 + PunctuationWeight}, // punctuation inpacts weight
		{"foo", MinimumWeight},               // minimal weight
	}

	for _, testCase := range testCases {
		result := CalcWeight(testCase.input)
		if result != testCase.expected {
			t.Errorf("CalcWeight(%s) expected %d, but got %d", testCase.input, testCase.expected, result)
		}
	}
}

func TestProcessString(t *testing.T) {
	tests := []struct {
		input    string
		expected []Word
	}{
		{
			input: "Hello World!",
			expected: []Word{
				{Text: "Hello", Weight: 0, inQuote: false},
				{Text: "World!", Weight: 0, inQuote: false},
			},
		},
		{
			input: "This is a 'test' string.",
			expected: []Word{
				{Text: "This", Weight: 0, inQuote: false},
				{Text: "is", Weight: 0, inQuote: false},
				{Text: "a", Weight: 0, inQuote: false},
				{Text: "'test'", Weight: 0, inQuote: true},
				{Text: "string.", Weight: 0, inQuote: false},
			},
		},
		{
			input: "This is a 'quoted test' string.",
			expected: []Word{
				{Text: "This", Weight: 0, inQuote: false},
				{Text: "is", Weight: 0, inQuote: false},
				{Text: "a", Weight: 0, inQuote: false},
				{Text: "'quoted", Weight: 0, inQuote: true},
				{Text: "test'", Weight: 0, inQuote: true},
				{Text: "string.", Weight: 0, inQuote: false},
			},
		},
	}

	for _, test := range tests {
		result := ParseWords(test.input)
		for i := range result {
			result[i].Weight = 0 // this text ignores the weights
		}
		if !reflect.DeepEqual(result, test.expected) {
			t.Errorf("Input: %s\nExpected: %v\nGot:      %v", test.input, test.expected, result)
		}
	}
}

func TestParseWords_ParagraphPause(t *testing.T) {
	got := ParseWords("foo bar\n\nbaz qux")
	if len(got) != 4 {
		t.Fatalf("expected 4 words, got %d", len(got))
	}
	if got[1].Weight-CalcWeight("bar") != ParagraphWeight {
		t.Errorf("last word of paragraph 1 missing ParagraphWeight bonus (got %d, base %d)",
			got[1].Weight, CalcWeight("bar"))
	}
	if got[3].Weight != CalcWeight("qux") {
		t.Errorf("last word overall should not get ParagraphWeight (got %d)", got[3].Weight)
	}
	if got[0].Weight != CalcWeight("foo") {
		t.Errorf("mid-paragraph word should have plain weight (got %d)", got[0].Weight)
	}
}
