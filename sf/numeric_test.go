package sf

import (
	"bytes"
	"strconv"
	"testing"
)

func TestDecimalAppend(t *testing.T) {

	b := make([]byte, 0, 32)

	for _, f := range []float64{0, decimalMin, decimalMax} {

		t.Run(strconv.FormatFloat(f, 'f', 5, 64), func(t *testing.T) {
			a, ok := decimalAppend(b, f)

			i := bytes.IndexByte(a, '.')
			if i < 0 {
				t.Errorf("missing decimal")
			} else if (f >= 0 && i > decimalDigits) || i > decimalDigits+1 {
				t.Errorf("too many digits before the decimal: %v", string(a))
			} else if len(a)-i > 1+decimalFracDigits {
				t.Errorf("too many digits after the decimal: %v", string(a))
			}

			if ok {
				g, ok := decimalParse(string(a))
				if !ok {
					t.Errorf("errors")
				} else if g != f {
					t.Errorf("got %v, expected %v", g, f)
				}
			}
		})
	}
}
