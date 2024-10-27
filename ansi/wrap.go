package ansi

import (
	"bytes"
	"unicode"
	"unicode/utf8"
)

// nbsp is a non-breaking space
const nbsp = 0xA0

// Hardwrap wraps a string or a block of text to a given line length, breaking
// word boundaries. This will preserve ANSI escape codes and will account for
// wide-characters in the string.
// When preserveSpace is true, spaces at the beginning of a line will be
// preserved.
func Hardwrap(s string, limit int, preserveSpace bool) string {
	if limit < 1 {
		return s
	}

	var (
		buf          bytes.Buffer
		curWidth     int
		forceNewline bool
	)

	addNewline := func() {
		buf.WriteByte('\n')
		curWidth = 0
	}

	for scanner := NewScanner(s, ScanRunes); scanner.Scan(); {
		cluster, width, control := scanner.Token()

		if control {
			buf.Write(cluster)
			continue
		}

		r, _ := utf8.DecodeRune(cluster)
		if r == '\n' {
			addNewline()
			forceNewline = false
			continue
		}
		if r == '\t' {
			width = 1
		}
		if curWidth+width > limit {
			addNewline()
			forceNewline = true
		}
		if curWidth == 0 {
			// Skip spaces at the beginning of a line
			if !preserveSpace && unicode.IsSpace(r) && len(cluster) <= 4 && r != utf8.RuneError && forceNewline {
				continue
			}
			forceNewline = false
		}

		buf.Write(cluster)
		curWidth += width
	}

	return buf.String()
}

// Wordwrap wraps a string or a block of text to a given line length, not
// breaking word boundaries. This will preserve ANSI escape codes and will
// account for wide-characters in the string.
// The breakpoints string is a list of characters that are considered
// breakpoints for word wrapping. A hyphen (-) is always considered a
// breakpoint.
//
// Note: breakpoints must be a string of 1-cell wide rune characters.
func Wordwrap(s string, limit int, breakpoints string) string {
	if limit < 1 {
		return s
	}

	var (
		buf      bytes.Buffer
		word     bytes.Buffer
		space    bytes.Buffer
		curWidth int
		wordLen  int
	)

	addSpace := func() {
		curWidth += space.Len()
		buf.Write(space.Bytes())
		space.Reset()
	}

	addWord := func() {
		if word.Len() == 0 {
			return
		}

		addSpace()
		curWidth += wordLen
		buf.Write(word.Bytes())
		word.Reset()
		wordLen = 0
	}

	addNewline := func() {
		buf.WriteByte('\n')
		curWidth = 0
		space.Reset()
	}

	for scanner := NewScanner(s, ScanRunes); scanner.Scan(); {
		cluster, width, isControl := scanner.Token()

		if isControl {
			word.Write(cluster)
			continue
		}

		switch r, _ := utf8.DecodeRune(cluster); {
		case r == '\n':
			if wordLen == 0 {
				if curWidth+space.Len() > limit {
					curWidth = 0
				} else {
					buf.Write(space.Bytes())
				}
				space.Reset()
			}
			addWord()
			addNewline()
		case r != utf8.RuneError && unicode.IsSpace(r) && r != nbsp:
			addWord()
			space.WriteRune(r)
		case r == '-':
			fallthrough
		case bytes.ContainsAny(cluster, breakpoints):
			addSpace()
			addWord()
			buf.Write(cluster)
			curWidth++
		default:
			word.Write(cluster)
			wordLen += width
			if curWidth+space.Len()+wordLen > limit &&
				wordLen < limit {
				addNewline()
			}
		}
	}

	addWord()

	return buf.String()
}

// Wrap wraps a string or a block of text to a given line length, breaking word
// boundaries if necessary. This will preserve ANSI escape codes and will
// account for wide-characters in the string. The breakpoints string is a list
// of characters that are considered breakpoints for word wrapping. A hyphen
// (-) is always considered a breakpoint.
//
// Note: breakpoints must be a string of 1-cell wide rune characters.
func Wrap(s string, limit int, breakpoints string) string {
	if limit < 1 {
		return s
	}

	var (
		buf      bytes.Buffer
		word     bytes.Buffer
		space    bytes.Buffer
		curWidth int // written width of the line
		wordLen  int // word buffer len without ANSI escape codes
	)

	addSpace := func() {
		curWidth += space.Len()
		buf.Write(space.Bytes())
		space.Reset()
	}

	addWord := func() {
		if word.Len() == 0 {
			return
		}

		addSpace()
		curWidth += wordLen
		buf.Write(word.Bytes())
		word.Reset()
		wordLen = 0
	}

	addNewline := func() {
		buf.WriteByte('\n')
		curWidth = 0
		space.Reset()
	}

	for scanner := NewScanner(s, ScanRunes); scanner.Scan(); {
		cluster, width, isControl := scanner.Token()

		if isControl {
			word.Write(cluster)
			continue
		}

		switch r, _ := utf8.DecodeRune(cluster); {
		case r == '\n':
			if wordLen == 0 {
				if curWidth+space.Len() > limit {
					curWidth = 0
				} else {
					// preserve whitespaces
					buf.Write(space.Bytes())
				}
				space.Reset()
			}
			addWord()
			addNewline()
		case r != utf8.RuneError && unicode.IsSpace(r) && r != nbsp: // nbsp is a non-breaking space
			addWord()
			space.WriteRune(r)
		case r == '-':
			fallthrough
		case bytes.ContainsAny(cluster, breakpoints):
			addSpace()
			if curWidth+wordLen+width > limit {
				word.Write(cluster)
				wordLen += width
			} else {
				addWord()
				buf.Write(cluster)
				curWidth += width
			}
		default:
			if curWidth == limit {
				addNewline()
			}
			if wordLen+width > limit {
				// Hardwrap the word if it's too long
				addWord()
			}

			word.Write(cluster)
			wordLen += width

			if wordLen == limit {
				// Hardwrap the word if it's too long
				addWord()
			}
			if curWidth+wordLen+space.Len() > limit {
				addNewline()
			}
		}
	}

	if word.Len() != 0 {
		// Preserve ANSI wrapped spaces at the end of string
		if curWidth+space.Len() > limit {
			buf.WriteByte('\n')
		}
		addSpace()
	}
	buf.Write(word.Bytes())

	return buf.String()
}
