package ansi

import (
	"strings"
)

// Transpose breaks a line into individual lines for each rune preserving ANSI
// escape codes which are distributed to each new line.
//
// todo minimize the ANSI codes. Currently if there are multiple codes they are
// just concatenated, which may lead to redundancy in some cases.
func Transpose(s string) string {
	var (
		sb      strings.Builder
		prefix  strings.Builder
		lines   = make([]strings.Builder, 0, len(s))
		scanner = NewScanner(s, ScanRunes)
	)

	for scanner.Scan() {
		p, _, isEscape := scanner.Token()
		if isEscape {
			prefix.Write(p)
			for i := range lines {
				lines[i].Write(p)
			}
			continue
		}
		n := len(lines)
		lines = append(lines, strings.Builder{})
		lines[n].WriteString(prefix.String())
		lines[n].Write(p)
	}

	for i, l := range lines {
		if i > 0 {
			sb.WriteString("\n")
		}
		sb.WriteString(l.String())
	}

	return sb.String()
}
