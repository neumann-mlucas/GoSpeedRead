package internal

import (
	"reflect"
	"testing"
)

func TestCalcWeight(t *testing.T) {
	// Test cases for CalcWeight
	testCases := []struct {
		input    string
		expected int
	}{
		{"hello", 98},
		{"world", 98},
		{"programming", 215},
		{"Hello, World!", 304},
		{"", 0},
		{"!", 69},
	}

	// Iterate over test cases
	for _, testCase := range testCases {
		result := CalcWeight(testCase.input)
		if result != testCase.expected {
			t.Errorf("CalcWeight(%s) expected %d, but got %d", testCase.input, testCase.expected, result)
		}
	}
}

func TestStartsWithAny(t *testing.T) {
	// Test cases for StartsWithAny
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

	// Iterate over test cases
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

	// Iterate over test cases
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
				{Text: "Hello", Weight: 98, inQuote: false},
				{Text: "World!", Weight: 167, inQuote: false},
			},
		},
		{
			input: "This is a 'test' string.",
			expected: []Word{
				{Text: "This", Weight: 78, inQuote: false},
				{Text: "is", Weight: 39, inQuote: false},
				{Text: "a", Weight: 19, inQuote: false},
				{Text: "'test'", Weight: 167, inQuote: true},
				{Text: "string.", Weight: 187, inQuote: false},
			},
		},
		{
			input: "This is a 'quoted test' string.",
			expected: []Word{
				{Text: "This", Weight: 78, inQuote: false},
				{Text: "is", Weight: 39, inQuote: false},
				{Text: "a", Weight: 19, inQuote: false},
				{Text: "'quoted", Weight: 187, inQuote: true},
				{Text: "test'", Weight: 148, inQuote: true},
				{Text: "string.", Weight: 187, inQuote: false},
			},
		},
	}

	for _, test := range tests {
		result := ProcessString(test.input)
		if !reflect.DeepEqual(result, test.expected) {
			t.Errorf("Input: %s\nExpected: %v\nGot:      %v", test.input, test.expected, result)
		}
	}
}
