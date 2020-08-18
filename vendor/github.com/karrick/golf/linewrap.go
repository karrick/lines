package golf

// https://en.wikipedia.org/wiki/Line_wrap_and_word_wrap
//
// 0x2028 LINE SEPARATOR
//   * may be used to represent this semantic unambiguously
// 0x2029 PARAGRAPH SEPARATOR
//   * may be used to represent this semantic unambiguously

// TODO: Hyphens are candidates for line breaking, and non-breaking hyphens are
// not.

import (
	"fmt"
	"unicode"
)

// Wrap returns input string with leading and trailing whitespace removed, with
// the exception of exactly one trailing newline character, and intraline
// whitespace normalized to single spaces, except in cases where a word would
// exceed a line length of 80 columns, in which case a newline will be used
// rather than a space.
func Wrap(input string) string {
	return defaultWrapper.Wrap(input)
}

var defaultWrapper = LineWrapper{Max: 80} // 80 columns and no prefix string

type wrapStates uint

const (
	wrapStatesBegin wrapStates = iota
	wrapStatesSpace
	wrapStatesWord
)

// LineWrapper wraps input strings to specified maximum number of columns, or
// when zero, 80 columns. When Prefix is not empty string, it prefixes each line
// with the prefix string.
type LineWrapper struct {
	// Prefix is the optional string to be places at the start of every output
	// line.
	Prefix string

	// Max is the optional maximum number of columns to wrap the text to. When
	// 0, the LineWrapper will wrap lines fit on in 80 columns.
	Max int
}

func (w LineWrapper) Wrap(input string) string {
	var wordStart int        // index of the loop is the starting position of the current rune, measured in bytes
	var lineLen, wordLen int // number of runes in each
	var state wrapStates     // state machine's state

	// if input == "" {
	// 	return input
	// }

	// Subtract one from the max column number because when printing in a
	// terminal of 80 columns, when we print into the 80th column, the terminal
	// continues the line and places the character on the following line.
	maxLineLength := w.Max - 1
	if maxLineLength == -1 {
		maxLineLength = 79 // When w.Max was 0, then this acts like it was 80
	}

	// Create a buffer large enough to hold entire line.
	output := make([]byte, 0, len(input))

	prefix := []byte{'\n'} // when no prefix, want to use newline as the prefix
	prefixLen := len(w.Prefix)
	if prefixLen > 0 {
		output = append(output, w.Prefix...) // copy in the prefix lines
		prefix = append(prefix, w.Prefix...) // newline followed by the prefix
	}

	for i, r := range input {
		switch state {
		case wrapStatesBegin:
			if isSpace(r) {
				state = wrapStatesSpace
				continue
			}
			state = wrapStatesWord
			wordLen = 1
			wordStart = i
		case wrapStatesSpace:
			if isSpace(r) {
				continue
			}
			state = wrapStatesWord
			wordLen = 1
			wordStart = i
		case wrapStatesWord:
			if !isSpace(r) {
				wordLen++
				continue
			}
			state = wrapStatesSpace
			if lineLen+wordLen >= maxLineLength {
				output = append(output, prefix...)
				lineLen = prefixLen
			} else if lineLen > 0 {
				output = append(output, ' ')
				lineLen++
			}
			output = append(output, input[wordStart:i]...)
			lineLen += wordLen
			wordLen = 0
		default:
			panic(fmt.Errorf("invalid state %v: template; %q; index: %d", state, input, i))
		}
	}

	switch state {
	case wrapStatesBegin:
		// empty input string: no-op
	case wrapStatesSpace:
		// final character was space: no-op
	case wrapStatesWord:
		// final character was word
		// if word not empty, append to line
		if lineLen+wordLen >= maxLineLength {
			output = append(output, prefix...)
		} else if lineLen > 0 {
			output = append(output, ' ')
		}
		output = append(output, input[wordStart:]...)
	default:
		panic(fmt.Errorf("invalid final state %v: template; %q", state, input))
	}

	return string(append(output, '\n'))
}

// isSpace returns true when r is a unicode space but not the non-breaking
// space, in which case it returns false.
func isSpace(r rune) bool {
	return r != '\u00a0' && unicode.IsSpace(r) // non-breaking space
}
