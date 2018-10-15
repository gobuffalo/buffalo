package build

import "path/filepath"

// AbsoluteBinaryPath returns the absolute path to the binary.
func (b *Builder) AbsoluteBinaryPath() string {
	if filepath.IsAbs(b.Bin) {
		return b.Bin
	}
	return filepath.Join(b.Root, b.Bin)
}
