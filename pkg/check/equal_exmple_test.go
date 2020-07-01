package check_test

import (
	"fmt"
	"jupiter-plugin/pkg/check"
)

func ExampleEqualFold() {
	fmt.Println(check.EqualFold("Go", "go"))

	// Output: true
}
