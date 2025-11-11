package sf

type InputString string

func (is InputString) parse(yield func(string) bool) {
	s := string(is)
	v, s, ok := itemCut(s)
	if !ok {
		return
	}
	i := 0
	for i < len(s) && s[i] == ' ' {
		i++
	}
	if i >= len(s) {
		yield(v)
		return
	}
	if s[i] != ',' {
		return
	}
	i++
	for {
		for i < len(s) && s[i] == ' ' {
			i++
		}
		if i >= len(s) {
			return
		}
		v, s, ok = itemCut(s[i:])

	}
}
