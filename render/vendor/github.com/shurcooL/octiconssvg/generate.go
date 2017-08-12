// +build ignore

package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"go/format"
	"io"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/shurcooL/go-goon"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

var oFlag = flag.String("o", "", "write output to `file` (default standard output)")

func main() {
	flag.Parse()

	err := run()
	if err != nil {
		log.Fatalln(err)
	}
}

func run() error {
	f, err := os.Open(filepath.Join("_data", "svg.json"))
	if err != nil {
		return err
	}
	defer f.Close()

	var octicons map[string]string
	err = json.NewDecoder(f).Decode(&octicons)
	if err != nil {
		return err
	}

	var names []string
	for name := range octicons {
		names = append(names, name)
	}
	sort.Strings(names)

	var buf bytes.Buffer
	fmt.Fprint(&buf, `package octiconssvg

import (
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

// Icon returns the named Octicon SVG node.
// It returns nil if name is not a valid Octicon symbol name.
func Icon(name string) *html.Node {
	switch name {
`)
	for _, name := range names {
		fmt.Fprintf(&buf, "	case %q:\n		return %v()\n", name, dashSepToMixedCaps(name))
	}
	fmt.Fprint(&buf, `	default:
		return nil
	}
}
`)
	for _, name := range names {
		processOcticon(&buf, octicons, name)
	}

	b, err := format.Source(buf.Bytes())
	if err != nil {
		return fmt.Errorf("error from format.Source(): %v", err)
	}

	var w io.Writer
	switch *oFlag {
	case "":
		w = os.Stdout
	default:
		f, err := os.Create(*oFlag)
		if err != nil {
			return err
		}
		defer f.Close()
		w = f
	}

	_, err = w.Write(b)
	return err
}

func processOcticon(w io.Writer, octicons map[string]string, name string) {
	svg := parseOcticon(octicons[name])

	// Clear these fields to remove cycles in the data structure, since go-goon
	// cannot print those in a way that's valid Go code. The generated data structure
	// is not a proper *html.Node with all fields set, but it's enough for rendering
	// to be successful.
	svg.LastChild = nil
	svg.FirstChild.Parent = nil

	fmt.Fprintln(w)
	fmt.Fprintf(w, "// %s returns an %q Octicon SVG node.\n", dashSepToMixedCaps(name), name)
	fmt.Fprintf(w, "func %s() *html.Node {\n", dashSepToMixedCaps(name))
	fmt.Fprint(w, "	return ")
	goon.Fdump(w, svg)
	fmt.Fprintln(w, "}")
}

func parseOcticon(svgXML string) *html.Node {
	e, err := html.ParseFragment(strings.NewReader(svgXML), nil)
	if err != nil {
		panic(fmt.Errorf("internal error: html.ParseFragment failed: %v", err))
	}
	svg := e[0].LastChild.FirstChild // TODO: Is there a better way to just get the <svg>...</svg> element directly, skipping <html><head></head><body><svg>...</svg></body></html>?
	svg.Parent.RemoveChild(svg)
	for i, attr := range svg.Attr {
		if attr.Namespace == "" && attr.Key == "width" {
			svg.Attr[i].Val = "16"
			break
		}
	}
	svg.Attr = append(svg.Attr, html.Attribute{Key: atom.Style.String(), Val: `fill: currentColor; vertical-align: top;`})
	return svg
}

// dashSepToMixedCaps converts "string-URL-append" to "StringURLAppend" form.
func dashSepToMixedCaps(in string) string {
	var out string
	ss := strings.Split(in, "-")
	for _, s := range ss {
		initialism := strings.ToUpper(s)
		if _, ok := initialisms[initialism]; ok {
			out += initialism
		} else {
			out += strings.Title(s)
		}
	}
	return out
}

// initialisms is the set of initialisms in Go-style Mixed Caps case.
var initialisms = map[string]struct{}{
	"API":   {},
	"ASCII": {},
	"CPU":   {},
	"CSS":   {},
	"DNS":   {},
	"EOF":   {},
	"GUID":  {},
	"HTML":  {},
	"HTTP":  {},
	"HTTPS": {},
	"ID":    {},
	"IP":    {},
	"JSON":  {},
	"LHS":   {},
	"QPS":   {},
	"RAM":   {},
	"RHS":   {},
	"RPC":   {},
	"SLA":   {},
	"SMTP":  {},
	"SQL":   {},
	"SSH":   {},
	"TCP":   {},
	"TLS":   {},
	"TTL":   {},
	"UDP":   {},
	"UI":    {},
	"UID":   {},
	"UUID":  {},
	"URI":   {},
	"URL":   {},
	"UTF8":  {},
	"VM":    {},
	"XML":   {},
	"XSRF":  {},
	"XSS":   {},

	"RSS": {},
}
