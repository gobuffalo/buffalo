package resolvers

import (
	"path/filepath"

	"github.com/gobuffalo/packr"
)

type PackrBox struct {
	Box packr.Box
}

func (p PackrBox) Read(name string) ([]byte, error) {
	return p.Box.MustBytes(name)
}

func (p PackrBox) Resolve(name string) (string, error) {
	return filepath.Join(p.Box.Path, name), nil
}
