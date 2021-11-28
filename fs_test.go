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

	fs := NewFS(embedded.FS(), "internal/testdata/disk")
	r.NotNil(fs)

	f, err := fs.Open("file.txt")
	r.NoError(err)

	b, err := io.ReadAll(f)
	r.NoError(err)

	r.Equal("This file is on disk.", string(b))
}

func Test_FS_Uses_Embed_If_No_Disk(t *testing.T) {
	r := require.New(t)

	fs := NewFS(embedded.FS(), "internal/testdata/empty")
	r.NotNil(fs)

	f, err := fs.Open("file.txt")
	r.NoError(err)

	b, err := io.ReadAll(f)
	r.NoError(err)

	r.Equal("This file is embedded.", string(b))
}
