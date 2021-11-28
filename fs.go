package buffalo

import (
	"fmt"
	"io"
	"io/fs"
	"os"
)

// FS wraps a directory and an embed FS that are expected to have the same contents.
// it prioritizes the directory FS and falls back to the embedded FS if the file cannot
// be found on disk. This is useful during development or when deploying with
// assets not embedded in the binary.
//
// Additionally FS hiddes any file named embed.go from the FS.
type FS struct {
	embed fs.FS
	dir   fs.FS
}

// NewFS returns a new FS that wraps the given directory and embedded FS.
// the embed.FS is expected to embed the same files as the directory FS.
func NewFS(embed fs.ReadDirFS, dir string) FS {
	return FS{
		embed: embed,
		dir:   os.DirFS(dir),
	}
}

// Open implements the FS interface.
func (f FS) Open(name string) (fs.File, error) {
	if name == "embed.go" {
		return nil, fs.ErrNotExist
	}
	file, err := f.getFile(name)
	if name == "." {
		return rootFile{file}, err
	}
	return file, err
}

func (f FS) getFile(name string) (fs.File, error) {
	file, err := f.dir.Open(name)
	if err == nil {
		return file, nil
	}

	return f.embed.Open(name)
}

// rootFile wraps the "." directory for hidding the embed.go file.
type rootFile struct {
	fs.File
}

// ReadDir implements the fs.ReadDirFile interface.
func (f rootFile) ReadDir(n int) (entries []fs.DirEntry, err error) {
	dir, ok := f.File.(fs.ReadDirFile)
	if !ok {
		return nil, fmt.Errorf("%T is not a directory", f.File)
	}

	if n <= 0 {
		entries, err = dir.ReadDir(n)
		entries = hideEmbedFile(entries)
	} else {
		entries, err = dir.ReadDir(n + 1)
		entries = hideEmbedFile(entries)
		if len(entries) > n {
			entries = entries[:n-1]
		}
	}

	if len(entries) == 0 {
		return entries, io.EOF
	}
	return entries, err
}

func hideEmbedFile(entries []fs.DirEntry) []fs.DirEntry {
	for i, entry := range entries {
		if entry.Name() == "embed.go" {
			entries = append(entries[:i], entries[i+1:]...)
		}
	}
	return entries
}
