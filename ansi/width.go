package ansi

import (
	"bytes"
)

// Strip removes ANSI escape codes from a string.
func Strip(s string) string {
	var (
		buf bytes.Buffer // buffer for collecting printable characters
	)

	for scanner := NewScanner(s); scanner.Scan(); {
		if scanner.IsEscape() {
			continue
		}
		buf.Write(scanner.Bytes())
	}

	return buf.String()
}

// StringWidth returns the width of a string in cells. This is the number of
// cells that the string will occupy when printed in a terminal. ANSI escape
// codes are ignored and wide characters (such as East Asians and emojis) are
// accounted for.
func StringWidth(s string) int {
	if s == "" {
		return 0
	}

	var (
		width int
	)

	for scanner := NewScanner(s); scanner.Scan(); {
		width += scanner.Width()
	}

	return width
}
