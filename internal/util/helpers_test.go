package util

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func Test_FormatDuration_Milliseconds(t *testing.T) {
	d := time.Millisecond * 350

	text := FormatDuration(d)

	assert.Equal(t, "350 ms", text)
}

func Test_FormatDuration_Seconds(t *testing.T) {
	d := time.Second * 4

	text := FormatDuration(d)

	assert.Equal(t, "4.00 s", text)
}

func Test_FormatDuration_Minutes(t *testing.T) {
	d := time.Minute * 2

	text := FormatDuration(d)

	assert.Equal(t, "2.00 m", text)
}

func Test_FormatDuration_Hours(t *testing.T) {
	d := time.Hour * 7

	text := FormatDuration(d)

	assert.Equal(t, "7.00 h", text)
}

func Test_WrapText_LessThanWidth(t *testing.T) {
	text := "cat dog foo"
	width := 15 // few characters longer than text

	wrapped, lines := WrapText(text, width)

	assert.Equal(t, 1, lines)
	assert.Equal(t, "cat dog foo", wrapped)
}

func Test_WrapText_EqualToWidth(t *testing.T) {
	text := "cat dog foo"
	width := 11 // exact the same as text

	wrapped, lines := WrapText(text, width)

	assert.Equal(t, 1, lines)
	assert.Equal(t, "cat dog foo", wrapped)
}

func Test_WrapText_LengthEvenlyDivisibleByWidth(t *testing.T) {
	text := "lorem ipsum lol" // 15 characters
	width := 3

	wrapped, lines := WrapText(text, width)

	assert.Equal(t, 5, lines)
	assert.Equal(t, "lor\nem \nips\num \nlol", wrapped)
}

func Test_WrapText_LengthNotEvenlyDivisibleByWidth(t *testing.T) {
	text := "lorem ipsum rofl" // 16 characters
	width := 3

	wrapped, lines := WrapText(text, width)

	assert.Equal(t, 6, lines)
	assert.Equal(t, "lor\nem \nips\num \nrof\nl", wrapped)
}

func Test_WrapText_EmbeddedNewlines(t *testing.T) {
	text := "lorem ipsum\nwhat"
	width := 3

	wrapped, lines := WrapText(text, width)

	assert.Equal(t, 6, lines)
	assert.Equal(t, "lor\nem \nips\num\nwha\nt", wrapped)
}
