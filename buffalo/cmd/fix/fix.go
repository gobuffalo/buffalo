package fix

import (
	"bufio"
	"fmt"
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

var replace = map[string]string{
	"github.com/markbates/pop":      "github.com/gobuffalo/pop",
	"github.com/markbates/validate": "github.com/gobuffalo/validate",
	"github.com/satori/go.uuid":     "github.com/gobuffalo/uuid",
}

var ic = ImportConverter{
	Data: replace,
}

var mr = MiddlewareTransformer{
	PackagesReplacement: map[string]string{
		"github.com/gobuffalo/buffalo/middleware/basicauth": "github.com/gobuffalo/mw-basicauth",
		"github.com/gobuffalo/buffalo/middleware/csrf":      "github.com/gobuffalo/mw-csrf",
		"github.com/gobuffalo/buffalo/middleware/i18n":      "github.com/gobuffalo/mw-i18n",
		"github.com/gobuffalo/buffalo/middleware/ssl":       "github.com/gobuffalo/mw-ssl",
		"github.com/gobuffalo/buffalo/middleware/tokenauth": "github.com/gobuffalo/mw-tokenauth",
	},

	Aliases: map[string]string{
		"github.com/gobuffalo/mw-basicauth":          "basicauth",
		"github.com/gobuffalo/mw-csrf":               "csrf",
		"github.com/gobuffalo/mw-i18n":               "i18n",
		"github.com/gobuffalo/mw-ssl":                "ssl",
		"github.com/gobuffalo/mw-tokenauth":          "tokenauth",
		"github.com/gobuffalo/mw-paramlogger":        "paramlogger",
		"github.com/gobuffalo/mw-contenttype":        "contenttype",
		"github.com/gobuffalo/buffalo-pop/pop/popmw": "popmw",
	},
}

var checks = []Check{
	ic.Process,
	mr.transformPackages,
	WebpackCheck,
	PackageJSONCheck,
	DepEnsure,
	DeprecrationsCheck,
}

func ask(q string) bool {
	fmt.Printf("? %s [y/n]\n", q)

	reader := bufio.NewReader(os.Stdin)
	text, _ := reader.ReadString('\n')

	text = strings.ToLower(strings.TrimSpace(text))
	return text == "y" || text == "yes"
}

type MiddlewareTransformer struct {
	PackagesReplacement map[string]string
	Aliases             map[string]string
}

func (mw MiddlewareTransformer) transformPackages(r *Runner) error {
	err := filepath.Walk(".", func(p string, fi os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if fi.IsDir() {
			return nil
		}

		for _, n := range []string{"vendor", "node_modules", ".git"} {
			if strings.HasPrefix(p, n+string(filepath.Separator)) {
				return nil
			}
		}

		ext := filepath.Ext(p)
		if ext != ".go" {
			return nil
		}

		read, err := ioutil.ReadFile(p)
		if err != nil {
			return err
		}

		newContents := string(read)
		newContents = strings.Replace(newContents, "middleware.SetContentType", "contenttype.Set", -1)
		newContents = strings.Replace(newContents, "middleware.AddContentType", "contenttype.Add", -1)
		newContents = strings.Replace(newContents, "middleware.ParameterLogger", "paramlogger.ParameterLogger", -1)
		newContents = strings.Replace(newContents, "middleware.PopTransaction", "popmw.Transaction", -1)

		err = ioutil.WriteFile(p, []byte(newContents), 0)
		if err != nil {
			return err
		}

		// create an empty fileset.
		fset := token.NewFileSet()

		f, err := parser.ParseFile(fset, p, nil, parser.ParseComments)
		if err != nil {
			e := err.Error()
			msg := "expected 'package', found 'EOF'"
			if e[len(e)-len(msg):] == msg {
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

		if strings.Contains(newContents, "paramlogger.ParameterLogger") {
			astutil.DeleteImport(fset, f, "github.com/gobuffalo/buffalo/middleware")
			astutil.AddNamedImport(fset, f, "paramlogger", "github.com/gobuffalo/mw-paramlogger")
		}

		if strings.Contains(newContents, "popmw.Transaction") {
			astutil.DeleteImport(fset, f, "github.com/gobuffalo/buffalo/middleware")
			astutil.AddNamedImport(fset, f, "popmw", "github.com/gobuffalo/buffalo-pop/pop/popmw")
		}

		if strings.Contains(newContents, "contenttype.Add") || strings.Contains(newContents, "contenttype.Set") {
			astutil.DeleteImport(fset, f, "github.com/gobuffalo/buffalo/middleware")
			astutil.AddNamedImport(fset, f, "contenttype", "github.com/gobuffalo/mw-contenttype")
		}

		// since the imports changed, resort them.
		ast.SortImports(fset, f)

		// create a temporary file, this easily avoids conflicts.
		temp, err := mw.writeTempResult(p, fset, f)
		if err != nil {
			return err
		}

		// rename the .temp to .go
		return os.Rename(temp, p)
	})
	return err
}

func (mw MiddlewareTransformer) writeTempResult(name string, fset *token.FileSet, f *ast.File) (string, error) {
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
