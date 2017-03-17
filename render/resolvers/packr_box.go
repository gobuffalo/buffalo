package resolvers

import (
	"path/filepath"

	"github.com/gobuffalo/packr"
)

// PackrBox for resolving files using Packr
type PackrBox struct {
	Box packr.Box
}

// Read a file from a Packr box
func (p PackrBox) Read(name string) ([]byte, error) {
	return p.Box.MustBytes(name)
}

// Resolve the file path of a file inside of a Packr box
func (p PackrBox) Resolve(name string) (string, error) {
	return filepath.Join(p.Box.Path, name), nil
}
