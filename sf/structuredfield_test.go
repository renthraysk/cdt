package sf

import (
	"slices"
	"strings"
	"testing"
)

var nines = strings.Repeat("9", 24)

var parseTests = []struct {
	name string // description of this test case
	// Named input parameters for target function.
	s    string
	want []string
}{
	{"empty", "", []string{}},
	{"true", "?1", []string{"?1"}},
	{"false", "?0", []string{"?0"}},
	{"boolean-invalid", "?x", []string{}},
	{"string-empty", `""`, []string{`""`}},
	{"string", `"abc"`, []string{`"abc"`}},
	{"string-invalid", `"abc`, []string{}},
	{"byte-sequence-empty", "::", []string{"::"}},
	{"byte-sequence-invalid", ":", []string{}},

	{"token-1", "aB*-", []string{"aB*-"}},
	{"token-2", "*abc", []string{"*abc"}},

	{"date", "@0", []string{"@0"}},
	{"date-invalid", "@A", []string{}},

	{"integer-0", "0", []string{"0"}},
	{"integer-1", "1", []string{"1"}},
	{"integer--1", "-1", []string{"-1"}},

	{"integer-max", nines[:integerDigits], []string{nines[:15]}},
	{"integer-min", "-" + nines[:integerDigits], []string{"-" + nines[:integerDigits]}},
	{"integer-over", nines[:integerDigits+1], []string{}},
	{"integer-under", "-" + nines[:integerDigits+1], []string{}},

	{"decimal", nines[:decimalDigits] + "." + nines[:decimalFracDigits], []string{nines[:decimalDigits] + "." + nines[:decimalFracDigits]}},
	{"decimal-over", nines[:decimalDigits+1] + ".1", []string{}},
	{"decimal-frac", "1." + nines[:decimalFracDigits+1], []string{}},

	{"display-sequence", `%"%E2%82%AC"`, []string{`%"%E2%82%AC"`}},
}

func TestParse(t *testing.T) {

	buf := make([]string, 0, 16)

	for _, tt := range parseTests {
		t.Run(tt.name, func(t *testing.T) {
			got := slices.AppendSeq(buf[:0], InputString(tt.s).Items)
			if !slices.Equal(got, tt.want) {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}

func BenchmarkParse(b *testing.B) {

	buf := make([]string, 0, 16)

	b.ReportAllocs()
	for b.Loop() {
		for _, tt := range parseTests {
			b := buf[:0]
			for v := range InputString(tt.s).Items {
				b = append(b, v)
			}
		}
	}
}
