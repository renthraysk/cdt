package sf

func itemCut(s string) (string, string, bool) {
	if len(s) <= 0 {
		return "", "", false
	}
	switch c := s[0]; {
	case isLower(c), isUpper(c):
		return tokenCut(s)
	case isDigit(c):
		return numericCut(s)
	default:
		switch c {
		case '"':
			return stringCut(s)
		case '*':
			return tokenCut(s)
		case '-':
			return numericCut(s)
		case ':':
			return byteSequenceCut(s)
		case '?':
			return boolCut(s)
		case '@':
			return dateCut(s)
		case '%':
			return displayStringCut(s)
		}
	}
	return "", s, false
}
