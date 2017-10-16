package grift

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/gobuffalo/makr"
	"github.com/pkg/errors"
)

//Run allows to create a new grift task generator
func (gg Generator) Run(root string, data makr.Data) error {
	g := makr.New()
	defer g.Fmt(root)

	header := tmplHeader
	path := filepath.Join("grifts", gg.Name.File()+".go")

	if _, err := os.Stat(path); err == nil {
		template, err := ioutil.ReadFile(path)
		if err != nil {
			return errors.WithStack(err)
		}
		header = string(template)
	}

	g.Add(makr.NewFile(path, header+tmplBody))

	data["opts"] = gg
	return g.Run(root, data)
}
