package ansi

import (
	"bytes"
)

// Truncate truncates a string to a given length, adding a tail to the
// end if the string is longer than the given length.
// This function is aware of ANSI escape codes and will not break them, and
// accounts for wide-characters (such as East Asians and emojis).
func Truncate(s string, length int, tail string) string {
	if sw := StringWidth(s); sw <= length {
		return s
	}

	tw := StringWidth(tail)
	length -= tw
	if length < 0 {
		return ""
	}

	var buf bytes.Buffer
	curWidth := 0
	ignoring := false

	// Here we iterate over the bytes of the string and collect printable
	// characters and runes. We also keep track of the width of the string
	// in cells.
	// Once we reach the given length, we start ignoring characters and only
	// collect ANSI escape codes until we reach the end of string.
	for scanner := NewScanner(s, ScanRunes); scanner.Scan(); {
		cluster, width, isControl := scanner.Token()
		if isControl {
			buf.Write(cluster)
			continue
		}
		if ignoring {
			continue
		}

		// Is this gonna be too wide?
		// If so write the tail and stop collecting.
		if curWidth+width > length && !ignoring {
			ignoring = true
			buf.WriteString(tail)
			continue
		}

		curWidth += width
		buf.Write(cluster)
	}

	return buf.String()
}
