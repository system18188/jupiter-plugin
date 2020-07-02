package check_test

import (
	"fmt"
	"github.com/system18188/jupiter-plugin/pkg/check"
)

func ExampleIsBlank() {
	fmt.Println(check.IsBlank(""))
	fmt.Println(check.IsBlank("  "))
	fmt.Println(check.IsBlank(" \t\r\n "))
	fmt.Println(check.IsBlank("a"))
	fmt.Println(check.IsBlank(" a  "))
	fmt.Println(check.IsBlank(" \t\r\n a \t\r\n "))

	// Output:
	// true
	// true
	// true
	// false
	// false
	// false
}

func ExampleIsAnyBlank() {
	fmt.Println(check.IsAnyBlank("", " "))
	fmt.Println(check.IsAnyBlank("  ", " a "))
	fmt.Println(check.IsAnyBlank(" \t\r\n ", " a "))
	fmt.Println(check.IsAnyBlank("a", " a "))
	fmt.Println(check.IsAnyBlank(" a ", " a "))
	fmt.Println(check.IsAnyBlank(" \t\r\n a \t\r\n ", " a "))

	// Output:
	// true
	// true
	// true
	// false
	// false
	// false
}

func ExampleIsNoneBlank() {
	fmt.Println(check.IsNoneBlank("", " "))
	fmt.Println(check.IsNoneBlank("  ", " a "))
	fmt.Println(check.IsNoneBlank(" \t\r\n ", " a "))
	fmt.Println(check.IsNoneBlank("a", " a "))
	fmt.Println(check.IsNoneBlank(" a ", " a "))
	fmt.Println(check.IsNoneBlank(" \t\r\n a \t\r\n ", " a "))

	// Output:
	// false
	// false
	// false
	// true
	// true
	// true
}

func ExampleIsNotBlank() {
	fmt.Println(check.IsNotBlank(""))
	fmt.Println(check.IsNotBlank("  "))
	fmt.Println(check.IsNotBlank(" \t\r\n "))
	fmt.Println(check.IsNotBlank("a"))
	fmt.Println(check.IsNotBlank(" a  "))
	fmt.Println(check.IsNotBlank(" \t\r\n a \t\r\n "))

	// Output:
	// false
	// false
	// false
	// true
	// true
	// true
}
