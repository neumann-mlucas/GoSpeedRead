package displaytext

import (
	"fmt"
	"github.com/atotto/clipboard"
	"internal/words"
	"strings"
	"time"
)

const (
	LOREM = `
Lorem ipsum dolor sit amet, "consectetur adipiscing elit. Nunc rutrum tincidunt massa". Pellentesque habitant morbi tristique senectus et netus et malesuada fames ac turpis egestas. Vestibulum vel elit ut nisi fermentum lacinia et et dolor. Aliquam lacinia et lacus eget ultrices. Phasellus at arcu ex. Donec interdum venenatis mi eu tempus. Proin convallis est libero, ut luctus dolor lobortis eu. Etiam a egestas odio. Donec vitae nisl mauris. Praesent interdum lectus quis odio maximus, a posuere turpis rutrum.
Donec quis libero ipsum. Nulla faucibus sapien pulvinar aliquet pretium. Aliquam quis mattis nulla. Donec sodales viverra convallis. Integer mollis luctus orci id vulputate. In eget sem lobortis, congue arcu ac, scelerisque velit. Quisque nec consequat massa. Etiam mattis accumsan porta. Nulla rutrum sapien felis, egestas fermentum elit molestie ut. Suspendisse viverra convallis risus, non rutrum metus lacinia sed.
Etiam a malesuada odio, et luctus nunc. Cras sit amet scelerisque justo, a tincidunt nibh. Sed aliquam pretium imperdiet. Etiam quis egestas lectus, eu volutpat dui. Vivamus congue tincidunt nulla, sit amet vulputate risus mattis ac. Ut dapibus odio nisi, eu posuere augue luctus eget. Donec eu arcu eu ante ultricies sollicitudin. Pellentesque habitant morbi tristique senectus et netus et malesuada fames ac turpis egestas. Donec nec arcu ut eros rhoncus vulputate nec eu magna. Cras hendrerit gravida libero vel condimentum. Vivamus augue massa, volutpat non interdum sit amet, ultricies nec turpis. Phasellus eget finibus enim, a tempus ligula. Maecenas metus velit, tincidunt in libero sed, consequat consectetur urna. Proin a interdum quam. Integer ultrices viverra tristique.
Nullam lobortis, tellus vitae finibus tincidunt, nunc lorem fermentum ante, sit amet aliquam elit felis quis massa. Cras venenatis tortor eget urna malesuada lacinia. Vestibulum molestie aliquam quam, non molestie felis venenatis vitae. Aenean lobortis tortor quis suscipit imperdiet. Proin mollis sit amet nulla a pretium. Aliquam id nunc libero. Curabitur ac malesuada mauris. Integer volutpat lorem sit amet sapien lobortis, id elementum quam sollicitudin. Cras dignissim sem eu urna varius, ac dignissim purus ornare. Pellentesque at ultrices felis. Morbi finibus eleifend lectus vitae dignissim. Maecenas et tortor at leo semper aliquam. Morbi quis dolor id urna aliquet consequat quis vel libero. Ut ac ante ac enim finibus consectetur pretium pharetra mi. Vivamus nibh magna, condimentum sit amet dolor id, malesuada commodo diam. Suspendisse potenti.
Donec in vulputate ipsum. Sed nisl dui, porttitor in ante non, pellentesque porta augue. Ut tincidunt porta tincidunt. Suspendisse sed nulla nec lacus porta pharetra sit amet aliquet lacus. Nam sed augue facilisis, feugiat ex eu, iaculis ex. Mauris dignissim felis ac turpis porttitor tincidunt. Vestibulum sit amet iaculis massa. Lorem ipsum dolor sit amet, consectetur adipiscing elit. Nulla ultrices auctor lacinia. Integer accumsan porta tortor, id maximus tortor efficitur at. Curabitur fermentum ante id libero lobortis, eu efficitur tortor convallis. Phasellus posuere egestas tortor sit amet venenatis. Vestibulum ut commodo nulla, ac tristique ante. Vestibulum placerat ipsum vitae enim dapibus vehicula. Integer at dui eleifend, scelerisque justo eget, efficitur risus. In sed odio at tortor venenatis viverra quis at ante.`
)

type DisplayText struct {
	Words []words.Word
	Index int
	WPM   int
}

type DisplayState struct {
	Text string
	Time time.Duration
	Prct float64
	Next string
	Prev string
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
	return DisplayState{
		Text: word.String(),
		Time: t.DisplayTime(word),
		Prct: t.Percentage(),
		Next: t.GetNextWords(5),
		Prev: t.GetPreviusWords(5),
	}
}

// DisplayTime calculates the display time in millisecond of a given Word
func (t *DisplayText) DisplayTime(word words.Word) time.Duration {
	wps := float64(t.WPM) / 60.0
	seconds := float64(word.Weight) / wps
	return time.Duration(seconds) * time.Millisecond
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
	return t.Index == 0 || t.Index >= len(t.Words)-1
}

// Percentage returns the relative position of the index
func (t *DisplayText) Percentage() float64 {
	return float64(t.Index) / float64(len(t.Words))
}

// GetPreviusWords returns a string with N words come before the current word
func (t *DisplayText) GetPreviusWords(i int) string {
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
