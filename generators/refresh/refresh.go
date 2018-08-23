package refresh

import (
	"fmt"

	"github.com/gobuffalo/buffalo/generators"
	"github.com/gobuffalo/makr"
	"github.com/gobuffalo/packr"
	"github.com/pkg/errors"
)

func init() {
	fmt.Println("github.com/gobuffalo/buffalo/generators/refresh has been deprecated in v0.13.0, and will be removed in v0.14.0. Use github.com/gobuffalo/buffalo/genny/refresh directly.")
}

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
