package trim_test

import (
	"fmt"

	"github.com/shurcooL/go/trim"
)

func ExampleLastNewline() {
	fmt.Printf("%q\n", trim.LastNewline("String\n"))
	fmt.Printf("%q\n", trim.LastNewline("String"))
	fmt.Printf("%q\n", trim.LastNewline("\n"))
	fmt.Printf("%q\n", trim.LastNewline(""))
	fmt.Printf("%q\n", trim.LastNewline("  String\n\n"))

	// Output:
	// "String"
	// "String"
	// ""
	// ""
	// "  String\n"
}

func ExampleFirstSpace() {
	fmt.Printf("%q\n", trim.FirstSpace(" String"))
	fmt.Printf("%q\n", trim.FirstSpace("String"))
	fmt.Printf("%q\n", trim.FirstSpace(" "))
	fmt.Printf("%q\n", trim.FirstSpace(""))
	fmt.Printf("%q\n", trim.FirstSpace("  String\n\n"))

	// Output:
	// "String"
	// "String"
	// ""
	// ""
	// " String\n\n"
}
