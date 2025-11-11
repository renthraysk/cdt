package sf

import (
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
	{"true", "?1", true, "?1", ""},
	{"false", "?0", true, "?0", ""},
	{"boolean-invalid", "?x", false, "", "?x"},
	{"string-empty", `""`, true, `""`, ""},
	{"string", `"abc"`, true, `"abc"`, ""},
	{"string-invalid", `"abc`, false, "", `"abc`},
	{"byte-sequence-empty", "::", true, "::", ""},
	{"byte-sequence-invalid", ":", false, "", ":"},
	{"token-1", "aB*-", true, "aB*-", ""},
	{"token-2", "*abc", true, "*abc", ""},
	{"date", "@0", true, "@0", ""},
	{"date-invalid", "@A", false, "", "@A"},
	{"integer-0", "0", true, "0", ""},
	{"integer-1", "1", true, "1", ""},
	{"integer--1", "-1", true, "-1", ""},
	{"integer-max", nines[:integerDigits], true, nines[:15], ""},
	{"integer-min", "-" + nines[:integerDigits], true, "-" + nines[:integerDigits], ""},
	{"integer-over", nines[:integerDigits+1], false, "", nines[:integerDigits+1]},
	{"integer-under", "-" + nines[:integerDigits+1], false, "", "-" + nines[:integerDigits+1]},
	{"decimal", nines[:decimalDigits] + "." + nines[:decimalFracDigits], true, nines[:decimalDigits] + "." + nines[:decimalFracDigits], ""},
	{"decimal-over", nines[:decimalDigits+1] + ".1", false, "", nines[:decimalDigits+1] + ".1"},
	{"decimal-frac", "1." + nines[:decimalFracDigits+1], false, "", "1." + nines[:decimalFracDigits+1]},
	{"display-sequence", `%"%E2%82%AC"`, true, `%"%E2%82%AC"`, ""},
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
			if r != tt.remain {
				t.Errorf("expected remain %q, got %q", tt.remain, r)
			}
			if got != tt.item {
				t.Errorf("got %v, want %v", got, tt.item)
			}
		})
	}
}
