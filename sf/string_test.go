package sf

import "testing"

var tests = []struct {
	name string // description of this test case
	// Named input parameters for target function.
	p         []byte
	escaped   string
	unescaped string
	ok        bool
}{
	{"empty", nil, "", "", true},
	{"simple", nil, "abc", "abc", true},

	{"esc1", nil, `\"abc\"`, `"abc"`, true},
	{"esc2", nil, `\\\"abc\\\"`, `\"abc\"`, true},

	{"non-print", nil, "ab\x80c", "", false},
	{"esc-with-non-print", nil, `a\\b` + "\x80" + `\\c`, "", false},

	{"non-print", nil, "", "ab\x80c", false},
	{"esc-with-non-print", nil, "", `a\\b` + "\x80" + `\\c`, false},
}

func Test_appendStringEscape(t *testing.T) {
	for _, tt := range tests {
		if !tt.ok && tt.unescaped == "" {
			continue
		}
		t.Run(tt.name, func(t *testing.T) {
			got, ok := appendStringEscape(tt.p, tt.unescaped)
			if ok != tt.ok {
				t.Errorf("appendStringEscape() = %v, want %v", ok, tt.ok)
			} else if string(got) != tt.escaped {
				t.Errorf("appendStringEscape() = %v, want %v", got, tt.escaped)
			}
		})
	}
}

func Test_appendStringUnescape(t *testing.T) {
	for _, tt := range tests {
		if !tt.ok && tt.escaped == "" {
			continue
		}
		t.Run(tt.name, func(t *testing.T) {

			got, ok := appendStringUnescape(tt.p, tt.escaped)
			if ok != tt.ok {
				t.Errorf("appendStringUnescape() = %v, want %v", ok, tt.ok)
			} else if string(got) != tt.unescaped {
				t.Errorf("appendStringUnescape() = %v, want %v", got, tt.unescaped)
			}
		})
	}
}
