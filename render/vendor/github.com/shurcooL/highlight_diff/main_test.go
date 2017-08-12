package highlight_diff_test

import (
	"github.com/shurcooL/go-goon"
	"github.com/shurcooL/highlight_diff"
)

func ExampleAnnotate() {
	anns, err := highlight_diff.Annotate([]byte(`@@ -1,6 +1,6 @@
 language: go
 go:
-  - 1.4
+  - 1.5
 install:
   - go get golang.org/x/tools/cmd/vet
 script:
`))
	if err != nil {
		panic(err)
	}

	goon.DumpExpr(len(anns))
	for _, ann := range anns {
		goon.DumpExpr(ann.Start, ann.End)
		goon.DumpExpr(string(ann.Left), string(ann.Right))
		goon.DumpExpr(ann.WantInner)
	}

	// Output:
	// len(anns) = (int)(1)
	// ann.Start = (int)(0)
	// ann.End = (int)(16)
	// string(ann.Left) = (string)("<span class=\"gu input-block\">")
	// string(ann.Right) = (string)("</span>")
	// ann.WantInner = (int)(0)
}
