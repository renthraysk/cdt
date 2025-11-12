package sf

// https://www.rfc-editor.org/rfc/rfc9651#name-tokens
// https://www.rfc-editor.org/rfc/rfc9651#name-parsing-a-token
func tokenCut(s string) (string, string, bool) {
	const (
		// tchar https://www.rfc-editor.org/rfc/rfc9110.html#name-tokens
		tchar = 1<<'!' | 1<<'#' | 1<<'$' | 1<<'%' | 1<<'&' | 1<<'\'' | 1<<'*' |
			1<<'+' | 1<<'-' | 1<<'.' | 1<<'^' | 1<<'_' | 1<<'`' | 1<<'|' | 1<<'~' |
			digit | alpha

		tokenStart = alpha | 1<<'*'
		token      = tchar | 1<<':' | 1<<'/'
	)

	if len(s) <= 0 || !isASCII(s[0], tokenStart%(1<<64), tokenStart>>64) {
		return "", s, false
	}
	i := 1
	for i < len(s) && isASCII(s[i], token%(1<<64), token>>64) {
		i++
	}
	if i >= len(s) {
		return s, "", true
	}
	return s[:i], s[i:], true
}
