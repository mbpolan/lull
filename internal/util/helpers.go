package util

import (
	"fmt"
	"golang.org/x/exp/constraints"
	"math"
	"strings"
	"time"
)

// ConsoleBell emits the default terminal bell sound.
func ConsoleBell() {
	fmt.Print("\a")
}

// Min returns the lesser of the two comparable values.
func Min[T constraints.Ordered](a, b T) T {
	if a < b {
		return a
	}

	return b
}

// FormatDuration returns a human friendly string representing the given duration (ie: 1.23 s).
func FormatDuration(t time.Duration) string {
	if t < time.Second {
		return fmt.Sprintf("%d ms", t.Milliseconds())
	} else if t < time.Minute {
		return fmt.Sprintf("%.2f s", t.Seconds())
	} else if t < time.Hour {
		return fmt.Sprintf("%.2f m", t.Minutes())
	} else {
		return fmt.Sprintf("%.2f h", t.Hours())
	}
}

// WrapText inserts newlines so that the given text contains at most width characters per line. The newly wrapped
// string and the number of lines are returned. Empty strings are considered to be on one line. If the text contains
// newlines, each line will be considered separately. This function does not consider words; it strictly wraps
// only on characters.
func WrapText(text string, width int) (string, int) {
	if text == "" || len(text) <= width {
		return text, 1
	}

	lines := strings.Split(text, "\n")
	totalLines := 0
	wrapped := ""

	for i, line := range lines {
		wrappedLine, count := WrapLine(line, width)
		wrapped += wrappedLine
		totalLines += count

		if i < len(lines)-1 {
			wrapped += "\n"
		}
	}

	return wrapped, totalLines
}

// WrapLine inserts newlines so that the given text contains at most width characters per line. The newly wrapped
// string and the number of lines are returned. This function does not consider existing newline characters in the
// given text; wse WrapText instead if that's the case.
func WrapLine(text string, width int) (string, int) {
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
