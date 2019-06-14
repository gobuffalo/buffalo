package fix

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"golang.org/x/tools/go/ast/astutil"
)

// ImportConverter will changes imports from a -> b
type ImportConverter struct {
	Data map[string]string
}

// Process will walk all the .go files in an application, excluding ./vendor.
// It will then attempt to convert any old import paths to any new import paths
// used by this version Buffalo.
func (c ImportConverter) Process(r *Runner) error {
	fmt.Println("~~~ Rewriting Imports ~~~")

	err := filepath.Walk(".", c.processFile)
	if err != nil {
		return err
	}

	if !r.App.WithDep {
		return nil
	}

	b, err := ioutil.ReadFile("Gopkg.toml")
	if err != nil {
		return err
	}

	for k := range c.Data {
		if bytes.Contains(b, []byte(k)) {
			r.Warnings = append(r.Warnings, fmt.Sprintf("Your Gopkg.toml contains the following import that need to be changed MANUALLY: %s", k))
		}
	}

	return nil
}

func (c ImportConverter) processFile(p string, info os.FileInfo, err error) error {
	er := onlyRelevantFiles(p, info, err, func(p string) error {
		return c.rewriteFile(p)
	})

	return er
}

func (c ImportConverter) rewriteFile(name string) error {

	// create an empty fileset.
	fset := token.NewFileSet()

	// parse the .go file.
	// we are parsing the entire file with comments, so we don't lose anything
	// if we need to write it back out.
	f, err := parser.ParseFile(fset, name, nil, parser.ParseComments)
	if err != nil {
		e := err.Error()
		msg := "expected 'package', found 'EOF'"
		if e[len(e)-len(msg):] == msg {
			return nil
		}
		return err
	}

	changed := false
	for key, value := range c.Data {
		if !astutil.DeleteImport(fset, f, key) {
			continue
		}

		astutil.AddImport(fset, f, value)
		changed = true
	}

	commentsChanged, err := c.handleFileComments(f)
	if err != nil {
		return err
	}

	changed = changed || commentsChanged

	// if no change occurred, then we don't need to write to disk, just return.
	if !changed {
		return nil
	}

	// since the imports changed, resort them.
	ast.SortImports(fset, f)

	// create a temporary file, this easily avoids conflicts.
	temp, err := writeTempResult(name, fset, f)
	if err != nil {
		return err
	}

	// rename the .temp to .go
	return os.Rename(temp, name)
}

func (c ImportConverter) handleFileComments(f *ast.File) (bool, error) {
	change := false

	for _, cg := range f.Comments {
		for _, cl := range cg.List {
			if !strings.HasPrefix(cl.Text, "// import \"") {
				continue
			}

			// trim off extra comment stuff
			ctext := cl.Text
			ctext = strings.TrimPrefix(ctext, "// import")
			ctext = strings.TrimSpace(ctext)

			// unquote the comment import path value
			ctext, err := strconv.Unquote(ctext)
			if err != nil {
				return false, err
			}

			// match the comment import path with the given replacement map
			if ctext, ok := c.match(ctext); ok {
				cl.Text = "// import " + strconv.Quote(ctext)
				change = true
			}

		}
	}

	return change, nil
}

// match takes an import path and replacement map.
func (c ImportConverter) match(importpath string) (string, bool) {
	for key, value := range c.Data {
		if !strings.HasPrefix(importpath, key) {
			continue
		}

		result := strings.Replace(importpath, key, value, 1)
		return result, true
	}

	return importpath, false
}
