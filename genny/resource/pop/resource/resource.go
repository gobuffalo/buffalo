package resource

import (
	"fmt"
	"html/template"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/gobuffalo/flect"
	"github.com/gobuffalo/genny"
	"github.com/gobuffalo/genny/movinglater/gotools"
	"github.com/gobuffalo/mapi"
	"github.com/gobuffalo/packr"
	"github.com/pkg/errors"
)

// New resource generator for pop
func New(opts *Options) (*genny.Group, error) {
	gg := &genny.Group{}

	if err := opts.Validate(); err != nil {
		return gg, errors.WithStack(err)
	}

	box := packr.NewBox("../resource/templates")

	tmplName := "resource-use_model"
	if opts.SkipModel {
		tmplName = "resource-name"
	}

	g := genny.New()
	g.Transformer(genny.Replace(tmplName, opts.Name.Tableize().String()))
	g.Transformer(genny.Replace("model-view-", opts.Name.Tableize().String()+string(filepath.Separator)))

	err := box.Walk(func(path string, f packr.File) error {
		// Adding the resource template to the generator
		if strings.Contains(path, tmplName) {
			fmt.Println("### tmplName ->", tmplName)
			fmt.Println("### path ->", path)
			g.File(genny.NewFile(path, f))
		}
		if !opts.SkipTemplates {
			if strings.Contains(path, "model-view-") {
				g.File(genny.NewFile(path, f))
			}
		}
		return nil
	})
	if err != nil {
		return gg, errors.WithStack(err)
	}
	data := mapi.Mapi{
		"opts": opts,
	}
	g.Transformer(gotools.TemplateTransformer(data, template.FuncMap{
		"camelize": flect.Camelize,
	}))

	if !opts.SkipModel && !opts.UseModel {
		g.Command(modelCommand(opts))
	}
	gg.Add(g)

	return gg, nil
}

func modelCommand(opts *Options) *exec.Cmd {
	args := opts.Args
	args = append(args[:0], args[0+1:]...)
	args = append([]string{"pop", "g", "model", opts.Name.Underscore().Singularize().String()}, args...)

	if opts.SkipMigration {
		args = append(args, "--skip-migration")
	}

	return exec.Command("buffalo", args...)
}
