package resource

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/gobuffalo/buffalo/generators"
	"github.com/gobuffalo/makr"
	"github.com/markbates/inflect"
)

// New generates a new actions/resource file and a stub test.
func New(data makr.Data) (*makr.Generator, error) {
	g := makr.New()
	files, err := generators.Find("resource")
	if err != nil {
		return nil, err
	}
	// Get the flags
	useModel := data["useModel"].(bool)
	skipModel := data["skipModel"].(bool)

	tmplName := "resource-use_model"

	if skipModel == true {
		tmplName = "resource-name"
	}
	for _, f := range files {
		if strings.Contains(f.WritePath, tmplName) {
			g.Add(makr.NewFile(strings.Replace(f.WritePath, tmplName, data["under"].(string), -1), f.Body))
		}
		if strings.Contains(f.WritePath, "model-view-") {
			targetPath := filepath.Join(
				filepath.Dir(f.WritePath),
				data["modelUnder"].(string),
				strings.Replace(filepath.Base(f.WritePath), "model-view-", "", -1),
			)
			g.Add(makr.NewFile(targetPath, f.Body))
		}
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
	if skipModel == false && useModel == false {
		g.Add(modelCommand(data))
	}

	g.Add(makr.NewCommand(makr.GoFmt()))

	return g, nil
}

func modelCommand(data makr.Data) makr.Command {
	modelName := inflect.Underscore(data["singular"].(string))

	args := data["args"].([]string)
	args = append(args[:0], args[0+1:]...)
	args = append([]string{"db", "g", "model", modelName}, args...)

	if skipMigration := data["skipMigration"].(bool); skipMigration == true {
		args = append(args, "--skip-migration")
	}

	return makr.NewCommand(exec.Command("buffalo", args...))
}
