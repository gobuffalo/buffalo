package goth

import (
	"path/filepath"

	"github.com/gobuffalo/buffalo/generators"
	"github.com/gobuffalo/makr"
)

// New actions/auth.go file configured to the specified providers.
func New() (*makr.Generator, error) {
	g := makr.New()
	files, err := generators.Find("goth")
	if err != nil {
		return nil, err
	}
	for _, f := range files {
		g.Add(makr.NewFile(f.WritePath, f.Body))
	}
	g.Add(&makr.Func{
		Should: func(data makr.Data) bool { return true },
		Runner: func(root string, data makr.Data) error {
			err := generators.AddInsideAppBlock("auth := app.Group(\"/auth\")",
				"auth.GET(\"/{provider}\", buffalo.WrapHandlerFunc(gothic.BeginAuthHandler))",
				"auth.GET(\"/{provider}/callback\", AuthCallback)")
			if err != nil {
				return err
			}
			return generators.AddImport(filepath.Join("actions", "app.go"), "github.com/markbates/goth/gothic")
		},
	})
	g.Add(makr.NewCommand(makr.GoGet("github.com/markbates/goth/...")))
	return g, nil
}
