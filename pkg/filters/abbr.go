package filters


import (
	"bytes"
	"github.com/system18188/jupiter-plugin/pkg/check"
)

// Abbr abbreviate a String using ellipses.
// This will turn "goyy" into "g...".
func Abbr(s string, maxWidth int) string {
	if check.IsBlank(s) || maxWidth < 1 {
		return ""
	}
	s = TrimSpace(s)
	if len(s) <= maxWidth {
		return s
	}
	var b bytes.Buffer
	i := 0
	for _, r := range s {
		if i < maxWidth-3 {
			b.WriteRune(r)
		}
		i++
	}
	b.WriteString("...")
	return b.String()
}

// Anon anonymous a String using asterisk.
// This will turn "goyy" into "g***y".
func Anon(s string) string {
	return Anonymous(s, 1, 1, 3)
}

// Anonymity a String using asterisk.
// This will turn "15566668888" into "155****8888".
func Anonymity(s string) string {
	return Anonymous(s, 3, 4, 4)
}

// Anonymous a String using asterisk.
// This will turn "goyy" into "g**y".
func Anonymous(s string, leftLen, rightLen, starLen int) string {
	if check.IsBlank(s) {
		return ""
	}
	s = TrimSpace(s)
	var bl, br, bs bytes.Buffer

	rs := Runes(s)

	ls := len(rs)
	if ls < 2 {
		return s
	}

	if leftLen < 1 {
		leftLen = 1
	}
	if rightLen < 1 {
		rightLen = 1
	}
	if leftLen+rightLen > ls {
		return s
	}
	if starLen < 1 {
		starLen = ls - leftLen - rightLen
	}
	for i := 0; i < starLen; i++ {
		bs.WriteString("*")
	}
	for i, r := range rs {
		if i < leftLen {
			bl.WriteRune(r)
		}
		if i >= len(rs)-rightLen {
			br.WriteRune(r)
		}
	}
	return bl.String() + bs.String() + br.String()
}
