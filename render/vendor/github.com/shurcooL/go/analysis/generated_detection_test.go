package analysis_test

import (
	"fmt"
	"os"

	"github.com/shurcooL/go/analysis"
)

func ExampleIsFileGenerated() {
	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	fmt.Println(analysis.IsFileGenerated(cwd, "testdata/generated_0.go.txt"))
	fmt.Println(analysis.IsFileGenerated(cwd, "testdata/handcrafted_0.go.txt"))
	fmt.Println(analysis.IsFileGenerated(cwd, "testdata/handcrafted_1.go.txt"))
	fmt.Println(analysis.IsFileGenerated(cwd, "vendor/github.com/shurcooL/go/analysis/file.go"))
	fmt.Println(analysis.IsFileGenerated(cwd, "subpkg/vendor/math/math.go"))

	// Output:
	// true <nil>
	// false <nil>
	// false <nil>
	// true <nil>
	// true <nil>
}
