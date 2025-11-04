package sf

import (
	"encoding/hex"
	"testing"
	"unicode/utf8"
)

func BenchmarkPrecent(b *testing.B) {

	b.ReportAllocs()
	for b.Loop() {
		percentDecode("%E2%82%AC")
	}
}

func testRune(t *testing.T, r rune) {
	var p [4]byte
	var h [12]byte

	t.Helper()
	b := utf8.AppendRune(p[:0], r)
	x := hex.AppendEncode(h[4:4], b)
	if len(x) == 0 {
		return
	}
	in := append(h[:0], '%', x[0], x[1])
	x = x[2:]
	for len(x) >= 2 {
		in = append(in, '%', x[0], x[1])
		x = x[2:]
	}

	gotR, gotN, _ := percentDecode(string(in))
	expectedR, expectedN := utf8.DecodeRune(b)

	if gotR != expectedR {
		t.Errorf("rune(0x%4X) rune expected 0x%4X, got 0x%4X", r, expectedR, gotR)
	}
	if gotN != 3*expectedN {
		t.Errorf("rune(0x%4X) length expected %d, got %d", r, expectedN, gotN)
	}
}

func TestPercentDecode(t *testing.T) {

	t.Run("0x0000-0xD7FF", func(t *testing.T) {
		for r := rune(0); r <= 0xD7FF; r++ {
			testRune(t, r)
		}
	})
	t.Run("0xD800-0xDFFF", func(t *testing.T) {
		for r := rune(0xD800); r <= 0xDFFF; r++ {
			testRune(t, r)
		}
	})
	t.Run("0xE000-0x10FFFF", func(t *testing.T) {
		for r := rune(0xE000); r <= utf8.MaxRune; r++ {
			testRune(t, r)
		}
	})
}

func TestPercentCut(t *testing.T) {
	tests := []struct {
		name   string
		in     string
		wantN  int
		wantOk bool
	}{
		{"empty", "", 0, false},
		{"no-percent", "abc", 0, false},
		{"too-short", "%A", 0, false},
		{"invalid-nib1", "%G1", 0, false},
		{"invalid-nib2", "%1G", 0, false},
		{"ascii", "%41rest", 3, true},
		{"ascii_upper_bound", "%7A_", 3, true},
		{"2-byte", "%C3%A9xyz", 6, true},
		{"2-byte_lower", "%c3%a9", 6, true},
		{"bad_continuation", "%C3%078", 0, false},
		{"3-byte", "%E2%82%ACend", 9, true},
		{"4-byte", "%F0%9F%98%80etc", 12, true},
		{"4-byte_invalid_second_nibble", "%F8%90%80%80", 0, false},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			_, n, ok := percentDecode(tc.in)
			if ok != tc.wantOk {
				t.Errorf("percentCut(%q) ok = (%v), want (%v)", tc.in, ok, tc.wantOk)
			}
			if !tc.wantOk {
				return
			}
			if tc.wantN != n {
				t.Errorf("percentCut(%q) n = (%d), want (%d)", tc.in, n, tc.wantN)
			}
		})
	}
}

func TestParseDisplayString(t *testing.T) {
	tests := []struct {
		name   string
		in     string
		want   string
		wantOk bool
	}{
		{
			name:   "valid ascii string",
			in:     "%\" !%22#$%25&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_`abcdefghijklmnopqrstuvwxyz{|}~\"",
			want:   " !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_`abcdefghijklmnopqrstuvwxyz{|}~",
			wantOk: true,
		},
		{
			name:   "valid lowercase non-ascii string",
			in:     `%"f%c3%bc%c3%bc"`,
			want:   "füü",
			wantOk: true,
		},
		{
			name: "invalid unqouted string",
			in:   "%foo",
		},
		{
			name: "invalid string missing initial quote",
			in:   `%foo"`,
		},
		{
			name: "invalid string missing closing quote",
			in:   `%"foo`,
		},
		{
			name: "invalid tab in string",
			in:   "%\"\t\"",
		},
		{
			name: "invalid newline in string",
			in:   "%\"\n\"",
		},
		{
			name: "invalid single quoted string",
			in:   `%'foo'`,
		},
		{
			name: "invalid string bad escaping",
			in:   `%\"foo %a"`,
		},
		{
			name:   "valid string with escaped quotes",
			in:     "%\"foo %22bar%22 \\ baz\"",
			want:   "foo \"bar\" \\ baz",
			wantOk: true,
		},
		{
			name: "invalid sequence id utf-8 string",
			in:   `%"%a0%a1"`,
		},
		{
			name: "invalid 2 bytes sequence utf-8 string",
			in:   `%"%c3%28"`,
		},
		{
			name: "invalid 3 bytes sequence utf-8 string",
			in:   `%"%e2%28%a1"`,
		},
		{
			name: "invalid 4 bytes sequence utf-8 string",
			in:   `%"%f0%28%8c%28"`,
		},
		{
			name: "invalid hex utf-8 string",
			in:   `%"%g0%1w"`,
		},
		{
			name:   "valid byte order mark in display string",
			in:     `%"BOM: %ef%bb%bf"`,
			want:   "BOM: \uFEFF",
			wantOk: true,
		},
		{
			name: "invalid unfinished 4 bytes rune",
			in:   `%"%f0%9f%98"`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			_, _, ok := displayStringCut(tc.in)
			if ok != tc.wantOk {
				t.Fatalf("test %q: want ok to be %v, got: %v", tc.name, tc.wantOk, ok)
			}
		})
	}
}
