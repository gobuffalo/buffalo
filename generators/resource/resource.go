package resource

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/gobuffalo/buffalo/generators"
	"github.com/gobuffalo/makr"
)

// New generates a new actions/resource file and a stub test.
func New(opts Options, data makr.Data) (*makr.Generator, error) {
	g := makr.New()
	data["opts"] = opts
	data["actions"] = []string{"List", "Show", "New", "Create", "Edit", "Update", "Destroy"}

	tmplName := "resource-use_model"

	mimeType := opts.MimeType
	if mimeType == "JSON" || mimeType == "XML" {
		tmplName = "resource-json-xml"
	} else if opts.SkipModel {
		tmplName = "resource-name"
	}

	files, err := generators.Find(filepath.Join(generators.TemplatesPath, "resource"))
	if err != nil {
		return nil, err
	}

	for _, f := range files {
		// Adding the resource template to the generator
		if strings.Contains(f.WritePath, tmplName) {
			folder := opts.FilesPath
			if strings.Contains(f.WritePath, "actions") {
				folder = opts.ActionsPath
			}
			p := strings.Replace(f.WritePath, tmplName, folder, -1)
			g.Add(makr.NewFile(p, f.Body))
		}
		if mimeType == "HTML" {
			// Adding the html templates to the generator
			if strings.Contains(f.WritePath, "model-view-") {
				targetPath := filepath.Join(
					filepath.Dir(f.WritePath),
					opts.FilesPath,
					strings.Replace(filepath.Base(f.WritePath), "model-view-", "", -1),
				)
				g.Add(makr.NewFile(targetPath, f.Body))
			}
		}
	}
	g.Add(&makr.Func{
		Should: func(data makr.Data) bool { return true },
		Runner: func(root string, data makr.Data) error {
			return generators.AddInsideAppBlock(fmt.Sprintf("app.Resource(\"/%s\", %sResource{&buffalo.BaseResource{}})", opts.Name.URL(), opts.Name.ModelPlural()))
		},
	})

	if !opts.SkipModel {
		g.Add(modelCommand(opts))
	}

	return g, nil
}

func modelCommand(opts Options) makr.Command {
	args := opts.Args
	args = append(args[:0], args[0+1:]...)
	args = append([]string{"db", "g", "model", opts.Model.UnderSingular()}, args...)

	if opts.SkipMigration {
		args = append(args, "--skip-migration")
	}
	return makr.NewCommand(exec.Command("buffalo", args...))
}
