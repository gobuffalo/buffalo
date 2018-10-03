package standard

import (
	"github.com/gobuffalo/buffalo/generators"
	"github.com/gobuffalo/makr"
	"github.com/gobuffalo/packr"
	"github.com/pkg/errors"
)

// Run standard assets generator for those wishing to not use webpack
func Run(root string, data makr.Data) error {
	files, err := generators.FindByBox(packr.NewBox("../standard/templates"))
	if err != nil {
		return errors.WithStack(err)
	}
	g := makr.New()
	for _, f := range files {
		g.Add(makr.NewFile(f.WritePath, f.Body))
	}
	return g.Run(root, data)
}
