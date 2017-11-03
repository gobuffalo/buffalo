package docker

import (
	"path/filepath"

	"github.com/gobuffalo/buffalo/generators"
	"github.com/gobuffalo/makr"
	"github.com/pkg/errors"
)

// Run Docker generator
func (d Generator) Run(root string, data makr.Data) error {
	g := makr.New()
	data["opts"] = d
	g.Add(&makr.Func{
		Should: func(data makr.Data) bool {
			return d.Style != "none"
		},
		Runner: func(root string, data makr.Data) error {
			style := d.Style
			if style != "multi" && style != "standard" {
				return errors.Errorf("unknown Docker style: %s", style)
			}
			files, err := generators.Find(filepath.Join(generators.TemplatesPath, "docker", style))
			if err != nil {
				return errors.WithStack(err)
			}
			fg := makr.New()
			for _, f := range files {
				fg.Add(makr.NewFile(f.WritePath, f.Body))
			}
			return fg.Run(root, data)
		},
	})
	return g.Run(root, data)
}
