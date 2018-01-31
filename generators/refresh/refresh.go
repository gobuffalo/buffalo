package refresh

import (
	"github.com/gobuffalo/buffalo/generators"
	"github.com/gobuffalo/makr"
	"github.com/gobuffalo/packr"
	"github.com/pkg/errors"
)

// Run generator for a .buffalo.dev.yml file
func Run(root string, data makr.Data) error {
	g := makr.New()

	files, err := generators.FindByBox(packr.NewBox("../refresh/templates"))
	if err != nil {
		return errors.WithStack(err)
	}

	for _, f := range files {
		g.Add(makr.NewFile(f.WritePath, f.Body))
	}

	return g.Run(root, data)
}
