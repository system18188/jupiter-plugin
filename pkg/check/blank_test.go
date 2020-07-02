package check_test

import (
	"github.com/system18188/jupiter-plugin/pkg/check"
	"testing"
)

func TestIsBlank(t *testing.T) {
	s := []struct {
		in  string
		out bool
	}{
		{"", true},
		{"  ", true},
		{" \t\r\n ", true},
		{"a", false},
		{" a ", false},
		{" \t\r\n a \t\r\n ", false},
	}
	for _, v := range s {
		d := check.IsBlank(v.in)
		if d != v.out {
			t.Errorf(`IsBlank("%v") = %v; want %v`, v.in, d, v.out)
		}
	}
}

func TestIsAnyBlank(t *testing.T) {
	s := []struct {
		in1, in2 string
		out      bool
	}{
		{"", " ", true},
		{"  ", " a ", true},
		{" \t\r\n ", " a ", true},
		{"a", " a ", false},
		{" a ", " a ", false},
		{" \t\r\n a \t\r\n ", " a ", false},
	}
	for _, v := range s {
		d := check.IsAnyBlank(v.in1, v.in2)
		if d != v.out {
			t.Errorf(`IsAnyBlank("%v", "%v") = %v; want %v`, v.in1, v.in2, d, v.out)
		}
	}
}

func TestIsNoneBlank(t *testing.T) {
	s := []struct {
		in1, in2 string
		out      bool
	}{
		{"", " ", false},
		{"  ", " a ", false},
		{" \t\r\n ", " a ", false},
		{"a", " a ", true},
		{" a ", " a ", true},
		{" \t\r\n a \t\r\n ", " a ", true},
	}
	for _, v := range s {
		d := check.IsNoneBlank(v.in1, v.in2)
		if d != v.out {
			t.Errorf(`IsNoneBlank("%v", "%v") = %v; want %v`, v.in1, v.in2, d, v.out)
		}
	}
}

func TestIsNotBlank(t *testing.T) {
	s := []struct {
		in  string
		out bool
	}{
		{"", false},
		{"  ", false},
		{" \t\r\n ", false},
		{"a", true},
		{" a ", true},
		{" \t\r\n a \t\r\n ", true},
	}
	for _, v := range s {
		d := check.IsNotBlank(v.in)
		if d != v.out {
			t.Errorf(`IsNotBlank("%v") = %v; want %v`, v.in, d, v.out)
		}
	}
}
