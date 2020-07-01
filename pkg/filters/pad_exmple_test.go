package filters_test

import (
	"fmt"
	"jupiter-plugin/pkg/filters"
)

func ExampleRepeat() {
	fmt.Println(filters.Repeat("-", 10))

	// Output: ----------
}

func ExamplePadStart() {
	fmt.Printf("%#q\n", filters.PadStart("bat", 5))

	// Output: `  bat`
}

func ExamplePadEnd() {
	fmt.Printf("%#q\n", filters.PadEnd("bat", 5))

	// Output: `bat  `
}

func ExamplePadLeft() {
	fmt.Println(filters.PadLeft("bat", 5, "*"))

	// Output: **bat
}

func ExamplePadRight() {
	fmt.Println(filters.PadRight("bat", 5, "*"))

	// Output: bat**
}

func ExamplePad() {
	fmt.Printf("%#q\n", filters.Pad("bat", 5))

	// Output: ` bat `
}

func ExampleCenter() {
	fmt.Println(filters.Center("bat", 8, "tag"))

	// Output: tabattag
}
