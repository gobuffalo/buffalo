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
	files, err := generators.Find(filepath.Join(generators.TemplatesPath, "resource"))
	if err != nil {
		return nil, err
	}
	// Get the flags
	useModel := data["useModel"].(string)
	skipModel := data["skipModel"].(bool)
	mimeType := data["mimeType"].(string)

	tmplName := "resource-use_model"

	if mimeType == "json" || mimeType == "xml" {
		tmplName = "resource-json-xml"
	} else if skipModel {
		tmplName = "resource-name"
	}

	for _, f := range files {
		// Adding the resource template to the generator
		if strings.Contains(f.WritePath, tmplName) {
			folder := data["filesPath"].(string)
			if strings.Contains(f.WritePath, "actions") {
				folder = data["actionsPath"].(string)
			}

			g.Add(makr.NewFile(strings.Replace(f.WritePath, tmplName, folder, -1), f.Body))
		}
		if mimeType == "html" {
			// Adding the html templates to the generator
			if strings.Contains(f.WritePath, "model-view-") {
				targetPath := filepath.Join(
					filepath.Dir(f.WritePath),
					data["filesPath"].(string),
					strings.Replace(filepath.Base(f.WritePath), "model-view-", "", -1),
				)
				g.Add(makr.NewFile(targetPath, f.Body))
			}
		}
	}
	g.Add(&makr.Func{
		Should: func(data makr.Data) bool { return true },
		Runner: func(root string, data makr.Data) error {
			return generators.AddInsideAppBlock(fmt.Sprintf("app.Resource(\"/%s\", %sResource{&buffalo.BaseResource{}})", data["resourceURL"], data["resourceName"]))
		},
	})

	if !skipModel && useModel == "" {
		g.Add(modelCommand(data))
	}

	return g, nil
}

func modelCommand(data makr.Data) makr.Command {
	modelName := inflect.Underscore(data["modelSingularUnder"].(string))

	args := data["args"].([]string)
	args = append(args[:0], args[0+1:]...)
	args = append([]string{"db", "g", "model", modelName}, args...)

	if skipMigration := data["skipMigration"].(bool); skipMigration {
		args = append(args, "--skip-migration")
	}

	return makr.NewCommand(exec.Command("buffalo", args...))
}
