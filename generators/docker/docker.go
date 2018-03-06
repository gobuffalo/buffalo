package docker

import (
	"github.com/gobuffalo/buffalo/generators"
	"github.com/gobuffalo/makr"
	"github.com/gobuffalo/packr"
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
			var box packr.Box
			switch d.Style {
			case "standard":
				box = packr.NewBox("./standard/templates")
			case "multi":
				box = packr.NewBox("./multi/templates")
			default:
				return errors.Errorf("unknown Docker style: %s", d.Style)
			}
			files, err := generators.FindByBox(box)
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
