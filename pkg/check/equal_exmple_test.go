package check_test

import (
	"fmt"
	"github.com/system18188/jupiter-plugin/pkg/check"
)

func ExampleEqualFold() {
	fmt.Println(check.EqualFold("Go", "go"))

	// Output: true
}
