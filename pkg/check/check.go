package check

import (
	"net"
	"strings"
)

// IsStrEmpty returns true if strings is empty otherwise false
func IsStrEmpty(v string) bool {
	return len(strings.TrimSpace(v)) == 0
}

// IsSliceContainsString method checks given string in the slice if found returns
// true otherwise false.
func IsSliceContainsString(search string, strSlice ...string) bool {
	for _, str := range strSlice {
		if strings.EqualFold(str, search) {
			return true
		}
	}
	return false
}

// IsDomain checks the specified string is domain name.
func IsDomain(s string) bool {
	return !IsIP(s) && "localhost" != s
}

// IsIP checks the specified string is IP.
func IsIP(s string) bool {
	return nil != net.ParseIP(s)
}
