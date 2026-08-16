package displaytext

import (
	"github.com/neumann-mlucas/GoSpeedRead/internal/words"
	"testing"
)

// Helper function to create a DisplayText instance with predefined words
func setupDisplayText(wordStr string) *DisplayText {
	dt := New(300)
	dt.Words = words.ParseWords(wordStr)
	return dt
}

func TestORPIndex(t *testing.T) {
	cases := []struct{ n, want int }{
		{0, 0}, {1, 0}, {2, 1}, {5, 1}, {6, 2}, {9, 2}, {10, 3}, {13, 3}, {14, 4}, {30, 4},
	}
	for _, c := range cases {
		if got := ORPIndex(c.n); got != c.want {
			t.Errorf("ORPIndex(%d) = %d, want %d", c.n, got, c.want)
		}
	}
}

func TestFocalOnText(t *testing.T) {
	cases := []struct {
		in   string
		want int
	}{
		{"hello", 1},         // 5 letters → ORP=1 → 'e'
		{" “ hello ” ", 4},   // strip 3 leading non-core, ORP=1 in 5-letter core, 3+1=4
		{"a", 0},             // 1 letter → ORP=0
		{"", 0},              // empty
		{"foo.", 1},          // 3 letters → ORP=1
		{"hello,", 1},        // core = "hello", ORP=1
	}
	for _, c := range cases {
		if got := FocalOnText(c.in); got != c.want {
			t.Errorf("FocalOnText(%q) = %d, want %d", c.in, got, c.want)
		}
	}
}

func TestNew(t *testing.T) {
	dt := New(120)
	if dt.Index != 0 || dt.WPM != 120 {
		t.Errorf("NewDisplayText() failed to initialize correctly")
	}
}

func TestDisplayText_Step(t *testing.T) {
	dt := setupDisplayText("Hello, World")
	state := dt.Step()
	if state.Text != "Hello" || dt.Index != 1 {
		t.Errorf("Step() did not advance correctly")
	}
	state = dt.Step()
	if state.Text != "World" || dt.Index != 2 {
		t.Errorf("Step() did not advance correctly")
	}
}

func TestDisplayText_IncIndex(t *testing.T) {
	words := "Lorem ipsum dolor sit amet"
	tests := []struct {
		name string
		arg  int
		want int
	}{
		{"inc 1", 1, 1},
		{"inc 2", 2, 2},
		{"inc 100", 100, 5},
		{"dec 1", -1, 0},
		{"dec 100", -100, 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dt := setupDisplayText(words)
			dt.IncIndex(tt.arg)
			if dt.Index != tt.want {
				t.Errorf("IncIndex() did not increment correctly W: %d R: %d", tt.want, dt.Index)
			}
		})
	}
}

func TestDisplayText_GetPreviousWords(t *testing.T) {
	words := "Lorem ipsum dolor sit amet"
	tests := []struct {
		name string
		arg  int
		want string
	}{
		{"next 1", 1, "sit"},
		{"next 2", 2, "dolor sit"},
		{"next 3", 3, "ipsum dolor sit"},
		{"next 100", 100, "Lorem ipsum dolor sit"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dt := setupDisplayText(words)
			dt.Index = 4
			r := dt.GetPreviousWords(tt.arg)
			if r != tt.want {
				t.Errorf("GetNextWords() error W: %s R: %s", tt.want, r)
			}
		})
	}
	t.Run("at start", func(t *testing.T) {
		dt := setupDisplayText(words)
		r := dt.GetPreviousWords(5)
		if r != "" {
			t.Errorf("GetNextWords() error W: %s R: %s", "", r)
		}
	})
}

func TestDisplayText_GetNextWords(t *testing.T) {
	words := "Lorem ipsum dolor sit amet"
	tests := []struct {
		name string
		arg  int
		want string
	}{
		{"next 1", 1, "ipsum"},
		{"next 2", 2, "ipsum dolor"},
		{"next 3", 3, "ipsum dolor sit"},
		{"next 100", 100, "ipsum dolor sit amet"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dt := setupDisplayText(words)

			r := dt.GetNextWords(tt.arg)
			if r != tt.want {
				t.Errorf("GetNextWords() error W: %s R: %s", tt.want, r)
			}
		})
	}
	t.Run("at end", func(t *testing.T) {
		dt := setupDisplayText(words)
		dt.Index = 4
		r := dt.GetNextWords(5)
		if r != "" {
			t.Errorf("GetNextWords() error W: %s R: %s", "", r)
		}
	})
}
