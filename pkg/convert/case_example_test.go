package convert_test

import (
	"fmt"
	"jupiter-plugin/pkg/convert"
)

func ExampleToLower() {
	s := []string{"",
		"abc",
		"AbC123",
		"azAZ09_"}
	for _, v := range s {
		fmt.Println(convert.ToLower(v))
	}

	// Output:
	//
	// abc
	// abc123
	// azaz09_
}

func ExampleToUpper() {
	s := []string{"",
		"abc",
		"AbC123",
		"azAZ09_"}
	for _, v := range s {
		fmt.Println(convert.ToUpper(v))
	}

	// Output:
	//
	// ABC
	// ABC123
	// AZAZ09_
}

func ExampleTitle() {
	s := []string{"",
		"a",
		"cat",
		"cAt",
		"aaa aaa aaa",
		"Aaa Aaa Aaa",
		"123a456",
		"double-blind",
		"ÿøû"}
	for _, v := range s {
		fmt.Println(convert.Title(v))
	}

	// Output:
	//
	// A
	// Cat
	// CAt
	// Aaa Aaa Aaa
	// Aaa Aaa Aaa
	// 123a456
	// Double-Blind
	// Ÿøû
}

func ExampleToTitle() {
	s := []string{"",
		"a",
		"cat",
		"cAt",
		"aaa aaa aaa",
		"Aaa Aaa Aaa",
		"123a456",
		"double-blind",
		"ÿøû"}
	for _, v := range s {
		fmt.Println(convert.ToTitle(v))
	}

	// Output:
	//
	// A
	// CAT
	// CAT
	// AAA AAA AAA
	// AAA AAA AAA
	// 123A456
	// DOUBLE-BLIND
	// ŸØÛ
}

func ExampleCamel() {
	s := []string{"",
		"a",
		"cat",
		"cAt",
		" aaa aaa aaa ",
		"_Aaa_Aaa_Aaa_",
		"123a456",
		"douBle-blind",
		"ÿøû"}
	for _, v := range s {
		fmt.Println(convert.Camel(v))
	}

	// Output:
	//
	// A
	// Cat
	// Cat
	// AaaAaaAaa
	// AaaAaaAaa
	// 123a456
	// DoubleBlind
	// Ÿøû
}

func ExampleUnCamel() {
	s := []string{"",
		"A",
		"Cat",
		"cAt",
		"AaaAaaAaa",
		" AaaAaaAaa ",
		"123a456",
		"DoubleBlind",
		"Ÿøû"}
	for _, v := range s {
		fmt.Println(convert.UnCamel(v, "_"))
	}

	// Output:
	//
	// a
	// cat
	// c_at
	// aaa_aaa_aaa
	// aaa_aaa_aaa
	// 123a456
	// double_blind
	// ÿøû
}
