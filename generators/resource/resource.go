package resource

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"

	"github.com/gobuffalo/buffalo/generators"
	"github.com/gobuffalo/makr"
)

// Run generates a new actions/resource file and a stub test.
func (res Generator) Run(root string, data makr.Data) error {
	g := makr.New()
	defer g.Fmt(root)

	data["opts"] = res
	data["actions"] = []string{"List", "Show", "New", "Create", "Edit", "Update", "Destroy"}

	tmplName := "resource-use_model"

	mimeType := res.MimeType
	if mimeType == "JSON" || mimeType == "XML" {
		tmplName = "resource-json-xml"
	} else if res.SkipModel {
		tmplName = "resource-name"
	}

	files, err := generators.Find(filepath.Join(generators.TemplatesPath, "resource"))
	if err != nil {
		return errors.WithStack(err)
	}

	for _, f := range files {
		// Adding the resource template to the generator
		if strings.Contains(f.WritePath, tmplName) {
			folder := res.FilesPath
			if strings.Contains(f.WritePath, "actions") {
				folder = res.ActionsPath
			}
			p := strings.Replace(f.WritePath, tmplName, folder, -1)
			g.Add(makr.NewFile(p, f.Body))
		}
		if mimeType == "HTML" {
			// Adding the html templates to the generator
			if strings.Contains(f.WritePath, "model-view-") {
				targetPath := filepath.Join(
					filepath.Dir(f.WritePath),
					res.FilesPath,
					strings.Replace(filepath.Base(f.WritePath), "model-view-", "", -1),
				)
				g.Add(makr.NewFile(targetPath, f.Body))
			}
		}
	}
	g.Add(&makr.Func{
		Should: func(data makr.Data) bool { return true },
		Runner: func(root string, data makr.Data) error {
			return generators.AddInsideAppBlock(fmt.Sprintf("app.Resource(\"/%s\", %sResource{&buffalo.BaseResource{}})", res.Name.URL(), res.Name.ModelPlural()))
		},
	})

	if !res.SkipModel {
		g.Add(res.modelCommand())
	}

	return g.Run(root, data)
}

func (res Generator) modelCommand() makr.Command {
	args := res.Args
	args = append(args[:0], args[0+1:]...)
	args = append([]string{"db", "g", "model", res.Model.UnderSingular()}, args...)

	if res.SkipMigration {
		args = append(args, "--skip-migration")
	}
	return makr.NewCommand(exec.Command("buffalo", args...))
}
