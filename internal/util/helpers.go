package util

import (
	"golang.org/x/exp/constraints"
	"math"
)

// Min returns the lesser of the two comparable values.
func Min[T constraints.Ordered](a, b T) T {
	if a < b {
		return a
	}

	return b
}

// WrapText inserts newlines so that the given text contains at most width characters per line. The newly wrapped
// string and the number of lines are returned. Empty strings are considered to be on one line. This function does not
// consider words; it strictly wraps on characters.
func WrapText(text string, width int) (string, int) {
	if text == "" || len(text) <= width {
		return text, 1
	}

	lines := int(math.Ceil(float64(len(text)) / float64(width)))
	wrapped := ""

	for i := 0; i < lines; i++ {
		left := i * width
		right := Min(left+width, len(text))
		wrapped += text[left:right]

		// do not add a newline after the last line
		if i < lines-1 {
			wrapped += "\n"
		}
	}

	return wrapped, lines
}
