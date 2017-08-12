package indentwriter_test

import (
	"io"
	"os"

	"github.com/shurcooL/go/indentwriter"
)

func Example() {
	iw := indentwriter.New(os.Stdout, 1)

	io.WriteString(iw, `IndentWriter is simple Go package you can import for the following task.

You take an existing io.Writer, and an integer "indent",
and create this IndentWriter that implements io.Writer too, but it prepends every line with
indent number of tabs.

Note that only non-empty lines get indented.
`)

	// Output:
	// 	IndentWriter is simple Go package you can import for the following task.
	//
	// 	You take an existing io.Writer, and an integer "indent",
	// 	and create this IndentWriter that implements io.Writer too, but it prepends every line with
	// 	indent number of tabs.
	//
	// 	Note that only non-empty lines get indented.
	//
}
