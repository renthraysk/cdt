package sf

func isTokenStart(c byte) bool {
	const ts = upper | lower | 1<<'*'

	return isASCII(c, ts%(1<<64), ts>>64)
}

func isToken(c byte) bool {
	const tchar = digit | upper | lower |
		1<<'!' | 1<<'#' | 1<<'$' | 1<<'%' | 1<<'&' | 1<<'\'' | 1<<'*' |
		1<<'+' | 1<<'-' | 1<<'.' | 1<<'^' | 1<<'_' | 1<<'`' | 1<<'|' | 1<<'~' |
		1<<':' | 1<<'/'
	const token = tchar | 1<<':' | 1<<'/'

	return isASCII(c, token%(1<<64), token>>64)
}

func tokenCut(s string) (string, string, bool) {
	if len(s) <= 0 || !isTokenStart(s[0]) {
		return "", s, false
	}
	i := 1
	for i < len(s) && isToken(s[i]) {
		i++
	}
	if i >= len(s) {
		return s, "", true
	}
	return s[:i], s[:i], true
}
