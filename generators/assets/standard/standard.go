package standard

import (
	"path/filepath"

	"github.com/gobuffalo/buffalo/generators/assets"
	"github.com/gobuffalo/buffalo/generators/common"
	"github.com/markbates/gentronics"
)

func New(data gentronics.Data) (*gentronics.Generator, error) {
	files, err := common.Find(filepath.Join("assets", "standard"))
	if err != nil {
		return nil, err
	}
	g := gentronics.New()
	for _, f := range files {
		g.Add(gentronics.NewFile(f.WritePath, f.Body))
	}
	g.Add(assets.PublicLogo)
	return g, nil
}
