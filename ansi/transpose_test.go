package ansi

import (
	"testing"
)

func Test(t *testing.T) {
	var tt = []struct {
		name   string
		input  string
		expect string
	}{
		{"normal string",
			"Vertical",
			"V\ne\nr\nt\ni\nc\na\nl",
		},
		{"string with color",
			"\x1b[21mVertical\x1b[0m",
			"\x1b[21mV\x1b[0m\n\x1b[21me\x1b[0m\n\x1b[21mr\x1b[0m\n\x1b[21mt\x1b[0m\n\x1b[21mi\x1b[0m\n\x1b[21mc\x1b[0m\n\x1b[21ma\x1b[0m\n\x1b[21ml\x1b[0m",
		},
	}

	for i, c := range tt {
		t.Run(c.name, func(t *testing.T) {
			if result := Transpose(c.input); result != c.expect {
				t.Errorf("test case %d failed: expected %q, got %q", i+1, c.expect, result)
			}
		})
	}
}
