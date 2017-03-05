package resource

import (
	"fmt"
	"strings"

	"github.com/gobuffalo/buffalo/generators"
	"github.com/markbates/gentronics"
)

// New generates a new actions/resource file and a stub test.
func New(data gentronics.Data) (*gentronics.Generator, error) {
	g := gentronics.New()
	files, err := generators.Find("resource")
	if err != nil {
		return nil, err
	}
	for _, f := range files {
		g.Add(gentronics.NewFile(strings.Replace(f.WritePath, "resource-name", data["under"].(string), -1), f.Body))
	}
	g.Add(&gentronics.Func{
		Should: func(data gentronics.Data) bool { return true },
		Runner: func(root string, data gentronics.Data) error {
			return generators.AddInsideAppBlock(fmt.Sprintf("var %sResource buffalo.Resource", data["downFirstCap"]),
				fmt.Sprintf("%sResource = %sResource{&buffalo.BaseResource{}}", data["downFirstCap"], data["camel"]),
				fmt.Sprintf("app.Resource(\"/%s\", %sResource)", data["under"], data["downFirstCap"]),
			)
		},
	})
	g.Add(gentronics.NewCommand(generators.GoFmt()))
	return g, nil
}
