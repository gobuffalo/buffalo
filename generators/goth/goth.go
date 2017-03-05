package goth

import (
	"path/filepath"

	"github.com/gobuffalo/buffalo/generators"
	"github.com/markbates/gentronics"
)

// New actions/auth.go file configured to the specified providers.
func New() (*gentronics.Generator, error) {
	g := gentronics.New()
	files, err := generators.Find("goth")
	if err != nil {
		return nil, err
	}
	for _, f := range files {
		g.Add(gentronics.NewFile(f.WritePath, f.Body))
	}
	g.Add(&gentronics.Func{
		Should: func(data gentronics.Data) bool { return true },
		Runner: func(root string, data gentronics.Data) error {
			err := generators.AddInsideAppBlock("auth := app.Group(\"/auth\")",
				"auth.GET(\"/{provider}\", buffalo.WrapHandlerFunc(gothic.BeginAuthHandler))",
				"auth.GET(\"/{provider}/callback\", AuthCallback)")
			if err != nil {
				return err
			}
			return generators.AddImport(filepath.Join("actions", "app.go"), "github.com/markbates/goth/gothic")
		},
	})
	g.Add(gentronics.NewCommand(generators.GoGet("github.com/markbates/goth/...")))
	g.Add(gentronics.NewCommand(generators.GoFmt()))
	return g, nil
}
