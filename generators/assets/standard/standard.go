package standard

import (
	"fmt"
	"path/filepath"

	"github.com/gobuffalo/buffalo/generators"
	"github.com/gobuffalo/buffalo/generators/assets"
	"github.com/gobuffalo/makr"
	"github.com/gobuffalo/packr"
	"github.com/pkg/errors"
)

func init() {
	fmt.Println("github.com/gobuffalo/buffalo/generators/assets/standard has been deprecated in v0.13.0, and will be removed in v0.14.0. Use github.com/gobuffalo/buffalo/genny/assets/standard directly.")
}

var logo = &makr.RemoteFile{
	File:       makr.NewFile(filepath.Join("public", "assets", "images", "logo.svg"), ""),
	RemotePath: assets.LogoURL,
}

// Run standard assets generator for those wishing to not use standard
func Run(root string, data makr.Data) error {
	files, err := generators.FindByBox(packr.NewBox("../standard/templates"))
	if err != nil {
		return errors.WithStack(err)
	}
	g := makr.New()
	for _, f := range files {
		g.Add(makr.NewFile(f.WritePath, f.Body))
	}
	g.Add(logo)
	return g.Run(root, data)
}
