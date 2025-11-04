package sf

func boolCut(s string) (string, string, bool) {
	if len(s) >= len("?0") && s[0] == '?' {
		switch s[1] {
		case '0':
			return "?0", s[2:], true
		case '1':
			return "?1", s[2:], true
		}
	}
	return "", s, false
}

func Bool(v []string) (bool, bool) {
	if len(v) != 1 {
		return false, false
	}
	s, r, ok := boolCut(v[0])
	if !ok || len(r) != 0 {
		return false, false
	}
	return boolParse(s)
}

func boolParse(s string) (bool, bool) {
	if s, r, ok := boolCut(s); ok && len(r) == 0 {
		switch s {
		case "?0":
			return false, true
		case "?1":
			return true, true
		}
	}
	return false, false
}

func BoolAppend(p []byte, b bool) ([]byte, bool) {
	c := 0
	if b {
		c = 1
	}
	return append(p, '?', '0'+byte(c)), true
}
