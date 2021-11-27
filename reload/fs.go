package reload

import (
	"io/fs"
	"os"
)

type FS struct {
	embed fs.FS
	dir   fs.FS
}

func NewFS(embed fs.ReadDirFS, dir string) FS {
	return FS{
		embed: embed,
		dir:   os.DirFS(dir),
	}
}

func (f FS) Open(name string) (fs.File, error) {
	if name == "embed.go" {
		return nil, fs.ErrNotExist
	}
	cfgFile := "./.buffalo.dev.yml"
	if _, err := os.Stat(cfgFile); err != nil {
		return f.embed.Open(name)
	}
	return f.dir.Open(name)
}
