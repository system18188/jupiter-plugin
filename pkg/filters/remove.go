package filters

import (
	"jupiter-plugin/pkg/check"
	"strings"
)

// Remove removes all occurrences of a substring from within the source string.
func Remove(s, remove string) string {
	if check.IsBlank(s) || check.IsBlank(remove) {
		return s
	}
	return strings.Replace(s, remove, "", -1)
}

// RemoveSpace removes all whitespaces from a string.
func RemoveSpace(s string) string {
	if check.IsBlank(s) {
		return s
	}
	sps := [...]string{"\t", "\n", "\v", "\f", "\r", " "}
	for _, sp := range sps {
		s = strings.Replace(s, sp, "", -1)
	}
	return s
}

// RemoveBlank removes all whitespaces from a string.
func RemoveBlank(s string) string {
	if check.IsBlank(s) {
		return s
	}

	s = TrimSpace(s)
	sps := [...]string{"\t", "\n", "\v", "\f", "\r", " "}
	for _, sp := range sps {
		s = strings.Replace(s, sp, "", -1)
	}
	return s
}


// RemoveStart removes a substring only if it is at the beginning of a source string,
// otherwise returns the source string.
func RemoveStart(s, remove string) string {
	if check.IsBlank(s) || check.IsBlank(remove) {
		return s
	}
	if strings.HasPrefix(s, remove) {
		return s[len(remove):]
	}
	return s
}

// RemoveEnd removes a substring only if it is at the end of a source string,
// otherwise returns the source string.
func RemoveEnd(s, remove string) string {
	if check.IsBlank(s) || check.IsBlank(remove) {
		return s
	}
	if strings.HasSuffix(s, remove) {
		return s[0 : len(s)-len(remove)]
	}
	return s
}

// 删除 STYLE
func RemoveURL(s string) string {
	return CompileUrl.ReplaceAllString(s, "")
}

// 删除HTML
func RemoveHtml(s string) string {
	return CompileHtml.ReplaceAllString(s, "")
}

// 删除 STYLE
func RemoveStyle(s string) string {
	return CompileStyle.ReplaceAllString(s, "")
}
