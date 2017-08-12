package highlight_go_test

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"sort"
	"text/template"

	"github.com/shurcooL/go/reflectsource"
	"github.com/shurcooL/highlight_go"
	"github.com/sourcegraph/annotate"
	"github.com/sourcegraph/syntaxhighlight"
)

// debugPrinter implements syntaxhighlight.Printer and prints the parameters it's given.
type debugPrinter struct{ syntaxhighlight.Printer }

func (p debugPrinter) Print(w io.Writer, kind syntaxhighlight.Kind, tokText string) error {
	fmt.Println(reflectsource.GetParentFuncArgsAsString(kind, tokText))

	return p.Printer.Print(w, kind, tokText)
}

func ExamplePrint() {
	src := []byte(`package main

import "fmt"

func main() {
	fmt.Println("Hey there, Go.")
}
`)

	// debugPrinter implements syntaxhighlight.Printer and prints the parameters it's given.
	p := debugPrinter{Printer: syntaxhighlight.HTMLPrinter(syntaxhighlight.DefaultHTMLConfig)}

	var buf bytes.Buffer
	highlight_go.Print(src, &buf, p)

	io.Copy(os.Stdout, &buf)

	// Output:
	// Print(syntaxhighlight.Keyword, "package")
	// Print(syntaxhighlight.Whitespace, " ")
	// Print(syntaxhighlight.Plaintext, "main")
	// Print(syntaxhighlight.Whitespace, "\n\n")
	// Print(syntaxhighlight.Keyword, "import")
	// Print(syntaxhighlight.Whitespace, " ")
	// Print(syntaxhighlight.String, "\"fmt\"")
	// Print(syntaxhighlight.Whitespace, "\n\n")
	// Print(syntaxhighlight.Keyword, "func")
	// Print(syntaxhighlight.Whitespace, " ")
	// Print(syntaxhighlight.Plaintext, "main")
	// Print(syntaxhighlight.Plaintext, "(")
	// Print(syntaxhighlight.Plaintext, ")")
	// Print(syntaxhighlight.Whitespace, " ")
	// Print(syntaxhighlight.Plaintext, "{")
	// Print(syntaxhighlight.Whitespace, "\n\t")
	// Print(syntaxhighlight.Plaintext, "fmt")
	// Print(syntaxhighlight.Plaintext, ".")
	// Print(syntaxhighlight.Plaintext, "Println")
	// Print(syntaxhighlight.Plaintext, "(")
	// Print(syntaxhighlight.String, "\"Hey there, Go.\"")
	// Print(syntaxhighlight.Plaintext, ")")
	// Print(syntaxhighlight.Whitespace, "\n")
	// Print(syntaxhighlight.Plaintext, "}")
	// Print(syntaxhighlight.Whitespace, "\n")
	// <span class="kwd">package</span> <span class="pln">main</span>
	//
	// <span class="kwd">import</span> <span class="str">&#34;fmt&#34;</span>
	//
	// <span class="kwd">func</span> <span class="pln">main</span><span class="pln">(</span><span class="pln">)</span> <span class="pln">{</span>
	// 	<span class="pln">fmt</span><span class="pln">.</span><span class="pln">Println</span><span class="pln">(</span><span class="str">&#34;Hey there, Go.&#34;</span><span class="pln">)</span>
	// <span class="pln">}</span>
}

func ExamplePrint_whitespace() {
	src := []byte("  package    main      \n\t\n")

	highlight_go.Print(src, ioutil.Discard, debugPrinter{Printer: syntaxhighlight.HTMLPrinter(syntaxhighlight.DefaultHTMLConfig)})

	// Output:
	// Print(syntaxhighlight.Whitespace, "  ")
	// Print(syntaxhighlight.Keyword, "package")
	// Print(syntaxhighlight.Whitespace, "    ")
	// Print(syntaxhighlight.Plaintext, "main")
	// Print(syntaxhighlight.Whitespace, "      \n\t\n")
}

// debugAnnotator implements syntaxhighlight.Annotator and prints the parameters it's given.
type debugAnnotator struct{ syntaxhighlight.Annotator }

func (a debugAnnotator) Annotate(start int, kind syntaxhighlight.Kind, tokText string) (*annotate.Annotation, error) {
	fmt.Println(reflectsource.GetParentFuncArgsAsString(start, kind, tokText))

	return a.Annotator.Annotate(start, kind, tokText)
}

func ExampleAnnotate() {
	src := []byte(`package main

import "fmt"

func main() {
	fmt.Println("Hey there, Go.")
}
`)

	// debugAnnotator implements syntaxhighlight.Annotator and prints the parameters it's given.
	p := debugAnnotator{Annotator: syntaxhighlight.HTMLAnnotator(syntaxhighlight.DefaultHTMLConfig)}

	anns, err := highlight_go.Annotate(src, p)
	if err != nil {
		log.Fatalln(err)
	}

	sort.Sort(anns)

	b, err := annotate.Annotate(src, anns, template.HTMLEscape)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println(string(b))

	// Output:
	// Annotate(0, syntaxhighlight.Keyword, "package")
	// Annotate(8, syntaxhighlight.Plaintext, "main")
	// Annotate(14, syntaxhighlight.Keyword, "import")
	// Annotate(21, syntaxhighlight.String, "\"fmt\"")
	// Annotate(28, syntaxhighlight.Keyword, "func")
	// Annotate(33, syntaxhighlight.Plaintext, "main")
	// Annotate(37, syntaxhighlight.Plaintext, "(")
	// Annotate(38, syntaxhighlight.Plaintext, ")")
	// Annotate(40, syntaxhighlight.Plaintext, "{")
	// Annotate(43, syntaxhighlight.Plaintext, "fmt")
	// Annotate(46, syntaxhighlight.Plaintext, ".")
	// Annotate(47, syntaxhighlight.Plaintext, "Println")
	// Annotate(54, syntaxhighlight.Plaintext, "(")
	// Annotate(55, syntaxhighlight.String, "\"Hey there, Go.\"")
	// Annotate(71, syntaxhighlight.Plaintext, ")")
	// Annotate(73, syntaxhighlight.Plaintext, "}")
	// <span class="kwd">package</span> <span class="pln">main</span>
	//
	// <span class="kwd">import</span> <span class="str">&#34;fmt&#34;</span>
	//
	// <span class="kwd">func</span> <span class="pln">main</span><span class="pln">(</span><span class="pln">)</span> <span class="pln">{</span>
	// 	<span class="pln">fmt</span><span class="pln">.</span><span class="pln">Println</span><span class="pln">(</span><span class="str">&#34;Hey there, Go.&#34;</span><span class="pln">)</span>
	// <span class="pln">}</span>
}
