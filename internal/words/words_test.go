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

func TestStartsWithAny(t *testing.T) {
	testCases := []struct {
		input    string
		chars    string
		expected bool
	}{
		{"hello", "helo", true},
		{"world", "abc", false},
		{"", "abc", false},
		{"apple", "ac", true},
	}

	for _, testCase := range testCases {
		result := StartsWithAny(testCase.input, testCase.chars)
		if result != testCase.expected {
			t.Errorf("StartsWithAny(%s, %s) expected %v, but got %v", testCase.input, testCase.chars, testCase.expected, result)
		}
	}
}

func TestEndsWithAny(t *testing.T) {
	// Test cases for EndsWithAny
	testCases := []struct {
		input    string
		chars    string
		expected bool
	}{
		{"hello", "lo", true},
		{"world", "abc", false},
		{"", "abc", false},
		{"apple", "ep", true},
	}

	for _, testCase := range testCases {
		result := EndsWithAny(testCase.input, testCase.chars)
		if result != testCase.expected {
			t.Errorf("EndsWithAny(%s, %s) expected %v, but got %v", testCase.input, testCase.chars, testCase.expected, result)
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
		for i := 0; i < len(result); i++ {
			result[i].Weight = 0 // this text ignores the weights
		}
		if !reflect.DeepEqual(result, test.expected) {
			t.Errorf("Input: %s\nExpected: %v\nGot:      %v", test.input, test.expected, result)
		}
	}
}
