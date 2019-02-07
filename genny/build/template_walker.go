package build

import (
	"bytes"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/gobuffalo/meta"
	"github.com/gobuffalo/packd"
	"github.com/karrick/godirwalk"
	"github.com/pkg/errors"
)

type dirWalker struct {
	dir string
}

func (d dirWalker) WalkPrefix(pre string, fn packd.WalkFunc) error {
	return d.Walk(func(path string, file packd.File) error {
		if strings.HasPrefix(path, pre) {
			return fn(path, file)
		}
		return nil
	})
}

func (d dirWalker) Walk(fn packd.WalkFunc) error {
	callback := func(path string, de *godirwalk.Dirent) error {
		if de != nil && de.IsDir() {
			base := filepath.Base(path)
			for _, pre := range []string{"vendor", ".", "_"} {
				if strings.HasPrefix(base, pre) {
					return filepath.SkipDir
				}
			}
			return nil
		}
		b, err := ioutil.ReadFile(path)
		if err != nil {
			return errors.WithStack(err)
		}
		f, err := packd.NewFile(path, bytes.NewReader(b))
		if err != nil {
			return errors.WithStack(err)
		}
		return fn(path, f)
	}

	godirwalk.Walk(d.dir, &godirwalk.Options{
		FollowSymbolicLinks: true,
		Callback:            callback,
	})
	return nil
}

func templateWalker(app meta.App) packd.Walkable {
	return dirWalker{dir: app.Root}
}
