package sf

import (
	"strconv"
	"strings"
	"testing"
)

var nines = strings.Repeat("9", 24)

var parseTests = []struct {
	name string // description of this test case

	inputString string
	ok          bool
	item        string
	remain      string
}{
	{"empty", "", false, "", ""},

	{"boolean-true", "?1", true, "?1", ""},
	{"boolean-false", "?0", true, "?0", ""},
	{"boolean-long", "?10", true, "?1", "0"}, // RFC only cares about first 2 characters
	{"boolean-neg", "?-1", false, "", "?-1"},
	{"boolean-plus", "?+1", false, "", "?+1"},
	{"boolean-invalid", "?x", false, "", "?x"},

	{"string-empty", `""`, true, `""`, ""},
	{"string", `"abc"`, true, `"abc"`, ""},
	{"string-invalid", `"abc`, false, "", `"abc`},
	{"string-escaped-quote", `"a\"b"`, true, `"a\"b"`, ""},
	{"string-trailing-backslash", `"abc\\"`, true, `"abc\\"`, ""},
	{"string-multi-line", `"abc\ndef"`, false, "", `"abc\ndef"`},

	{"byte-sequence-empty", "::", true, "::", ""},
	{"byte-sequence-invalid", ":", false, "", ":"},
	{"byte-sequence-single-colon", ":", false, "", ":"},
	{"byte-sequence-with-content", ":abc:", true, ":abc:", ""},
	{"byte-sequence-percent-empty", `:%`, false, "", `:%`},
	{"byte-sequence-percent-invalid", `:%G`, false, "", `:%G`},
	{"byte-sequence-mixed", `:%41:abc:`, false, `:%41:abc:`, ""},

	{"token-1", "aB*-", true, "aB*-", ""},
	{"token-2", "*abc", true, "*abc", ""},
	{"token-empty", "", false, "", ""},
	{"token-special-start", "!abc", false, "", "!abc"},
	{"token-special-end", "abc!", true, "abc!", ""},
	{"token-space", "a b", true, "a", " b"},
	{"token-long", strings.Repeat("a", 1000), true, strings.Repeat("a", 1000), ""}, // Assuming reasonable length limit

	{"integer-0", "0", true, "0", ""},
	{"integer-1", "1", true, "1", ""},
	{"integer--1", "-1", true, "-1", ""},
	{"integer-max", nines[:integerDigits], true, nines[:15], ""},
	{"integer-min", "-" + nines[:integerDigits], true, "-" + nines[:integerDigits], ""},
	{"integer-over", nines[:integerDigits+1], false, "", nines[:integerDigits+1]},
	{"integer-under", "-" + nines[:integerDigits+1], false, "", "-" + nines[:integerDigits+1]},
	{"integer-decimal-mixed", "1.0", true, "1.0", ""},

	{"date", "@0", true, "@0", ""},
	{"date-invalid", "@A", false, "", "@A"},
	{"date-neg", "@-1", true, "@-1", ""},
	{"date-overflow", "@" + strconv.FormatInt(integerMax+1, 10), false, "", "@" + strconv.FormatInt(integerMax+1, 10)},
	{"date-leading-zero", "@000", true, "@000", ""},

	{"decimal", nines[:decimalDigits] + "." + nines[:decimalFracDigits], true, nines[:decimalDigits] + "." + nines[:decimalFracDigits], ""},
	{"decimal-over", nines[:decimalDigits+1] + ".1", false, "", nines[:decimalDigits+1] + ".1"},
	{"decimal-frac", "1." + nines[:decimalFracDigits+1], false, "", "1." + nines[:decimalFracDigits+1]},
	{"decimal-neg", "-1.5", true, "-1.5", ""},
	{"decimal-leading-dot", ".5", false, "", ".5"}, // No leading dot support
	{"decimal-trailing-dot", "1.", true, "1.", ""},
	{"decimal-plus", "+1.5", false, "", "+1.5"},
	{"decimal-leading-zero-frac", "1.05", true, "1.05", ""},
	{"decimal-no-digits", ".", false, "", "."},
	{"decimal-frac-leading-zero", "1.005", true, "1.005", ""},

	{"display-sequence", `%"%E2%82%AC"`, true, `%"%E2%82%AC"`, ""},
	{"display-sequence-empty-percent", "%", false, "", "%"},
	{"display-sequence-invalid-percent", "%G", false, "", "%G"},
	{"display-sequence-mixed", `%"abc"%E2%82%AC`, true, `%"abc"`, "%E2%82%AC"},

	{"partial-parse-remaining", "abc def", true, "abc", " def"}, // Assuming token parse leaves remainder
	{"multi-item", "?1 \"abc\"", true, "?1", " \"abc\""},        // Full input not single item
	{"whitespace-prefix", " ?1", false, "", " ?1"},
	{"whitespace-suffix", "?1 ", true, "?1", " "}, // Assuming trim not automatic
}

func TestItemCut(t *testing.T) {
	for _, tt := range parseTests {
		t.Run(tt.name, func(t *testing.T) {
			got, r, ok := itemCut(tt.inputString)

			if ok != tt.ok {
				t.Errorf("expected ok %v, got %v", tt.ok, ok)
			}
			if !tt.ok {
				return
			}
			if got != tt.item {
				t.Errorf("expected item %v, got %v", tt.item, got)
			}
			if r != tt.remain {
				t.Errorf("expected remain %q, got %q", tt.remain, r)
			}
		})
	}
}
