package check_test

import (
	"github.com/system18188/jupiter-plugin/pkg/check"
	"testing"
)

func TestEqualFold(t *testing.T) {
	s := []struct {
		s, t string
		out  bool
	}{
		{"abc", "abc", true},
		{"ABcd", "ABcd", true},
		{"123abc", "123ABC", true},
		{"αβδ", "ΑΒΔ", true},
		{"abc", "xyz", false},
		{"abc", "XYZ", false},
		{"abcdefghijk", "abcdefghijX", false},
		{"abcdefghijk", "abcdefghij\u212A", true},
		{"abcdefghijK", "abcdefghij\u212A", true},
		{"abcdefghijkz", "abcdefghij\u212Ay", false},
		{"abcdefghijKz", "abcdefghij\u212Ay", false},
	}
	for _, v := range s {
		if out := check.EqualFold(v.s, v.t); out != v.out {
			t.Errorf("EqualFold(%#q, %#q) = %v, want %v", v.s, v.t, out, v.out)
		}
		if out := check.EqualFold(v.t, v.s); out != v.out {
			t.Errorf("EqualFold(%#q, %#q) = %v, want %v", v.t, v.s, out, v.out)
		}
	}
}
