package standard

import (
	"path/filepath"

	"github.com/gobuffalo/buffalo/generators/assets"
	"github.com/gobuffalo/buffalo/generators/common"
	"github.com/markbates/gentronics"
)

var logo = &gentronics.RemoteFile{
	File:       gentronics.NewFile("public/assets/images/logo.svg", ""),
	RemotePath: assets.LogoURL,
}

// New standard assets generator for those wishing to not use webpack
func New(data gentronics.Data) (*gentronics.Generator, error) {
	files, err := common.Find(filepath.Join("assets", "standard"))
	if err != nil {
		return nil, err
	}
	g := gentronics.New()
	for _, f := range files {
		g.Add(gentronics.NewFile(f.WritePath, f.Body))
	}
	g.Add(logo)
	return g, nil
}
