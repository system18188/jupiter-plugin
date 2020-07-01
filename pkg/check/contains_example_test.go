package check_test

import (
	"fmt"
	"jupiter-plugin/pkg/check"

)

func ExampleContains() {
	fmt.Println(check.Contains("seafood", "foo"))
	fmt.Println(check.Contains("seafood", "bar"))
	fmt.Println(check.Contains("seafood", ""))
	fmt.Println(check.Contains("", ""))
	fmt.Println(check.Contains("", "foo"))

	// Output:
	// true
	// false
	// true
	// true
	// false
}

func ExampleContainsAny() {
	fmt.Println(check.ContainsAny("team", "i"))
	fmt.Println(check.ContainsAny("failure", "u & i"))
	fmt.Println(check.ContainsAny("foo", ""))
	fmt.Println(check.ContainsAny("", ""))

	// Output:
	// false
	// true
	// false
	// false
}

func ExampleContainsRune() {
	fmt.Println(check.ContainsRune("team", 'i'))
	fmt.Println(check.ContainsRune("failure", 'u'))

	// Output:
	// false
	// true
}

func ExampleContainsSpace() {
	fmt.Println(check.ContainsSpace("a"))
	fmt.Println(check.ContainsSpace(" a "))
	fmt.Println(check.ContainsSpace("ab c"))
	fmt.Println(check.ContainsSpace("ab\tc"))
	fmt.Println(check.ContainsSpace("ab\rc"))
	fmt.Println(check.ContainsSpace("ab\nc"))

	// Output:
	// false
	// true
	// true
	// true
	// true
	// true
}

func ExampleContainsOnly() {
	fmt.Println(check.ContainsOnly("abab", "abc"))
	fmt.Println(check.ContainsOnly("ab1", "abc"))
	fmt.Println(check.ContainsOnly("abz", "abc"))

	// Output:
	// true
	// false
	// false
}

func ExampleContainsNone() {
	fmt.Println(check.ContainsNone("abab", "xyz"))
	fmt.Println(check.ContainsNone("ab1", "xyz"))
	fmt.Println(check.ContainsNone("abz", "xyz"))

	// Output:
	// true
	// true
	// false
}

func ExampleContainsSlice() {
	fmt.Println(check.ContainsSlice("abcd", []string{"ab", "cd"}))
	fmt.Println(check.ContainsSlice("ab1", []string{"ab", "12"}))
	fmt.Println(check.ContainsSlice("abz", []string{"xy", "cd"}))
	fmt.Println(check.ContainsSlice("abz", []string{"", ""}))
	fmt.Println(check.ContainsSlice("", []string{"xy", "cd"}))

	// Output:
	// true
	// false
	// false
	// true
	// false
}

func ExampleContainsSliceAny() {
	fmt.Println(check.ContainsSliceAny("abab", []string{"ab", "cd"}))
	fmt.Println(check.ContainsSliceAny("ab1", []string{"abc", "12"}))
	fmt.Println(check.ContainsSliceAny("abz", []string{"xy", "cd"}))
	fmt.Println(check.ContainsSliceAny("abz", []string{"", ""}))
	fmt.Println(check.ContainsSliceAny("", []string{"xy", "cd"}))

	// Output:
	// true
	// false
	// false
	// true
	// false
}
