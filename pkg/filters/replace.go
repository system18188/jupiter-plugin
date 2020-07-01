package filters

import (
	"bytes"
	"strings"
)

// Replace returns a copy of the string s with the first n
// non-overlapping instances of old replaced by new.
// If n < 0, there is no limit on the number of replacements.
func Replace(s, old, new string, n int) string { return strings.Replace(s, old, new, n) }

// 将HTML标签全转换成小写
func ReplaceHtmlLower(s string) string {
	return CompileHtml.ReplaceAllStringFunc(s, strings.ToLower)
}

// 将HTML标签全转换成小写
func ReplaceBytesHtmlLower(b []byte) []byte {
	return CompileHtml.ReplaceAllFunc(b, bytes.ToLower)
}