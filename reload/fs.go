package reload

import (
	"fmt"
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
	file, err := f.getFile(name)
	if name != "." {
		return file, err
	}
	return rootFile{file}, err
}

func (f FS) getFile(name string) (fs.File, error) {
	cfgFile := "./.buffalo.dev.yml"
	if _, err := os.Stat(cfgFile); err != nil {
		return f.embed.Open(name)
	}
	return f.dir.Open(name)
}

type rootFile struct {
	fs.File
}

func (f rootFile) ReadDir(n int) ([]fs.DirEntry, error) {
	dir, ok := f.File.(fs.ReadDirFile)
	if !ok {
		return nil, fmt.Errorf("failed at hiding embed.go file")
	}
	entries, err := dir.ReadDir(n)
	if err != nil {
		return entries, err
	}

	for i, entry := range entries {
		if entry.Name() == "embed.go" {
			entries = append(entries[:i], entries[i+1:]...)
		}
	}
	return entries, nil
}
