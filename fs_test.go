package buffalo

import (
	"io"
	"io/fs"
	"testing"

	"github.com/gobuffalo/buffalo/internal/testdata/embedded"
	"github.com/stretchr/testify/require"
)

func Test_FS_Disallows_Parent_Folders(t *testing.T) {
	r := require.New(t)

	fsys := NewFS(embedded.FS(), "internal/testdata/disk")
	r.NotNil(fsys)

	f, err := fsys.Open("../panic.txt")
	r.ErrorIs(err, fs.ErrNotExist)
	r.Nil(f)

	f, err = fsys.Open("try/../to/../trick/../panic.txt")
	r.ErrorIs(err, fs.ErrNotExist)
	r.Nil(f)
}

func Test_FS_Hides_embed_go(t *testing.T) {
	r := require.New(t)

	fsys := NewFS(embedded.FS(), "internal/testdata/disk")
	r.NotNil(fsys)

	f, err := fsys.Open("embed.go")
	r.ErrorIs(err, fs.ErrNotExist)
	r.Nil(f)
}

func Test_FS_Prioritizes_Disk(t *testing.T) {
	r := require.New(t)

	fsys := NewFS(embedded.FS(), "internal/testdata/disk")
	r.NotNil(fsys)

	f, err := fsys.Open("file.txt")
	r.NoError(err)

	b, err := io.ReadAll(f)
	r.NoError(err)

	r.Equal("This file is on disk.", string(b))
}

func Test_FS_Uses_Embed_If_No_Disk(t *testing.T) {
	r := require.New(t)

	fsys := NewFS(embedded.FS(), "internal/testdata/empty")
	r.NotNil(fsys)

	f, err := fsys.Open("file.txt")
	r.NoError(err)

	b, err := io.ReadAll(f)
	r.NoError(err)

	r.Equal("This file is embedded.", string(b))
}

func Test_FS_ReadDirFile(t *testing.T) {
	r := require.New(t)

	fsys := NewFS(embedded.FS(), "internal/testdata/disk")
	r.NotNil(fsys)

	f, err := fsys.Open(".")
	r.NoError(err)

	dir, ok := f.(fs.ReadDirFile)
	r.True(ok, "folder does not implement fs.ReadDirFile interface")

	// First read should return 1 file
	entries, err := dir.ReadDir(1)
	r.NoError(err)

	// The actual len will be 0 because the first file read is the embed.go file
	// this is counter-intuitive, but it's how the fs.ReadDirFile interface is specified;
	// if err == nil, just continue to call ReadDir until io.EOF is returned.
	r.LessOrEqual(len(entries), 1, "a call to ReadDir must at most return n entries")

	// First read should return at most 2 files
	entries, err = dir.ReadDir(2)
	r.NoError(err)

	// The actual len will be 2 (file.txt & file2.txt)
	r.LessOrEqual(len(entries), 2, "a call to ReadDir must at most return n entries")

	// trying to read next 2 files (none left)
	entries, err = dir.ReadDir(2)
	r.ErrorIs(err, io.EOF)
	r.Empty(entries)
}
