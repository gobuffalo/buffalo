package buffalo

import (
	"fmt"
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

// Open opens the named file.
//
// When Open returns an error, it should be of type *PathError with the Op
// field set to "open", the Path field set to name, and the Err field
// describing the problem.
//
// Open should reject attempts to open names that do not satisfy
// ValidPath(name), returning a *PathError with Err set to ErrInvalid or
// ErrNotExist.
func (f FS) Open(name string) (fs.File, error) {
	if name == "embed.go" {
		return nil, &fs.PathError{
			Op:   "open",
			Path: name,
			Err:  fs.ErrNotExist,
		}
	}
	file, err := f.getFile(name)
	if name == "." {
		// NOTE: It always returns the root from the "disk" instead
		// "embed". However, it could be fine since the the purpose
		// of buffalo.FS isn't supporting full featured filesystem.
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

	entries, err = dir.ReadDir(n)
	entries = hideEmbedFile(entries)
	return entries, err
}

func hideEmbedFile(entries []fs.DirEntry) []fs.DirEntry {
	result := make([]fs.DirEntry, 0, len(entries))

	for _, entry := range entries {
		if entry.Name() != "embed.go" {
			result = append(result, entry)
		}
	}
	return result
}
