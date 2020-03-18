package fix

import (
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/tools/go/ast/astutil"
)

//MiddlewareTransformer moves from our old middleware package to new one
type MiddlewareTransformer struct {
	PackagesReplacement map[string]string
	Aliases             map[string]string
}

func (mw MiddlewareTransformer) transformPackages(r *Runner) error {
	return filepath.Walk(".", mw.processFile)
}

func (mw MiddlewareTransformer) processFile(p string, fi os.FileInfo, err error) error {
	er := onlyRelevantFiles(p, fi, err, func(p string) error {
		if err := mw.rewriteMiddlewareUses(p); err != nil {
			return err
		}

		fset, f, err := buildASTFor(p)
		if err != nil {
			if e := err.Error(); strings.Contains(e, "expected 'package', found 'EOF'") {
				return nil
			}

			return err
		}

		//Replacing mw packages
		for old, new := range mw.PackagesReplacement {
			deleted := astutil.DeleteImport(fset, f, old)
			if deleted {
				astutil.AddNamedImport(fset, f, mw.Aliases[new], new)
			}
		}

		if err := mw.addMissingRootMiddlewareImports(fset, f, p); err != nil {
			return err
		}

		ast.SortImports(fset, f)

		temp, err := writeTempResult(p, fset, f)
		if err != nil {
			return err
		}

		// rename the .temp to .go
		return os.Rename(temp, p)
	})

	return er
}

func (mw MiddlewareTransformer) addMissingRootMiddlewareImports(fset *token.FileSet, f *ast.File, p string) error {
	read, err := ioutil.ReadFile(p)
	if err != nil {
		return err
	}

	content := string(read)

	astutil.DeleteImport(fset, f, "github.com/gobuffalo/buffalo/middleware")
	if strings.Contains(content, "paramlogger.ParameterLogger") {
		astutil.AddNamedImport(fset, f, "paramlogger", "github.com/gobuffalo/mw-paramlogger")
	}

	if strings.Contains(content, "popmw.Transaction") {
		astutil.AddImport(fset, f, "github.com/gobuffalo/buffalo-pop/v2/pop/popmw")
	}

	if strings.Contains(content, "contenttype.Add") || strings.Contains(content, "contenttype.Set") {
		astutil.AddNamedImport(fset, f, "contenttype", "github.com/gobuffalo/mw-contenttype")
	}

	return ioutil.WriteFile(p, []byte(content), 0)
}

func (mw MiddlewareTransformer) rewriteMiddlewareUses(p string) error {
	read, err := ioutil.ReadFile(p)
	if err != nil {
		return err
	}

	newContents := string(read)
	newContents = strings.Replace(newContents, "middleware.SetContentType", "contenttype.Set", -1)
	newContents = strings.Replace(newContents, "middleware.AddContentType", "contenttype.Add", -1)
	newContents = strings.Replace(newContents, "middleware.ParameterLogger", "paramlogger.ParameterLogger", -1)
	newContents = strings.Replace(newContents, "middleware.PopTransaction", "popmw.Transaction", -1)
	newContents = strings.Replace(newContents, "ssl.ForceSSL", "forcessl.Middleware", -1)

	err = ioutil.WriteFile(p, []byte(newContents), 0)
	return err
}

func writeTempResult(name string, fset *token.FileSet, f *ast.File) (string, error) {
	temp := name + ".temp"
	w, err := os.Create(temp)
	if err != nil {
		return "", err
	}

	// write changes to .temp file, and include proper formatting.
	err = (&printer.Config{Mode: printer.TabIndent | printer.UseSpaces, Tabwidth: 8}).Fprint(w, fset, f)
	if err != nil {
		return "", err
	}

	// close the writer
	err = w.Close()
	if err != nil {
		return "", err
	}

	return temp, nil
}

func buildASTFor(p string) (*token.FileSet, *ast.File, error) {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, p, nil, parser.ParseComments)
	return fset, f, err
}

//onlyRelevantFiles processes only .go files excluding folders like node_modules and vendor.
func onlyRelevantFiles(p string, fi os.FileInfo, err error, fn func(p string) error) error {
	if err != nil {
		return err
	}

	if fi.IsDir() {
		base := filepath.Base(p)
		if strings.HasPrefix(base, "_") {
			return filepath.SkipDir
		}
		for _, n := range []string{"vendor", "node_modules", ".git"} {
			if base == n {
				return filepath.SkipDir
			}
		}
		return nil
	}

	ext := filepath.Ext(p)
	if ext != ".go" {
		return nil
	}

	return fn(p)
}
