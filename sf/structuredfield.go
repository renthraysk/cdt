package sf

type InputString string

func (is InputString) Items(yield func(i string) bool) {
	for s := string(is); len(s) > 0; {
		var i string
		var ok bool

		switch c := s[0]; {
		case isLower(c), isUpper(c):
			i, s, ok = tokenCut(s)
		case isDigit(c):
			i, s, ok = numericCut(s)
		default:
			switch c {
			case '"':
				i, s, ok = stringCut(s)
			case '*':
				i, s, ok = tokenCut(s)
			case '-':
				i, s, ok = numericCut(s)
			case ':':
				i, s, ok = byteSequenceCut(s)
			case '?':
				i, s, ok = boolCut(s)
			case '@':
				i, s, ok = dateCut(s)
			case '%':
				i, s, ok = displayStringCut(s)
			}
		}
		if !ok || !yield(i) {
			return
		}
	}
}
