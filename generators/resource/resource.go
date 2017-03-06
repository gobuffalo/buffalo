package resource

import (
	"fmt"
	"strings"

	"github.com/gobuffalo/buffalo/generators"
	"github.com/gobuffalo/makr"
)

// New generates a new actions/resource file and a stub test.
func New(data makr.Data) (*makr.Generator, error) {
	g := makr.New()
	files, err := generators.Find("resource")
	if err != nil {
		return nil, err
	}
	for _, f := range files {
		g.Add(makr.NewFile(strings.Replace(f.WritePath, "resource-name", data["under"].(string), -1), f.Body))
	}
	g.Add(&makr.Func{
		Should: func(data makr.Data) bool { return true },
		Runner: func(root string, data makr.Data) error {
			return generators.AddInsideAppBlock(fmt.Sprintf("var %sResource buffalo.Resource", data["downFirstCap"]),
				fmt.Sprintf("%sResource = %sResource{&buffalo.BaseResource{}}", data["downFirstCap"], data["camel"]),
				fmt.Sprintf("app.Resource(\"/%s\", %sResource)", data["under"], data["downFirstCap"]),
			)
		},
	})
	g.Add(makr.NewCommand(makr.GoFmt()))
	return g, nil
}
