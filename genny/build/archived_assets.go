package build

import (
	"archive/zip"
	"bytes"
	"io"
	"path/filepath"
	"strings"

	"github.com/gobuffalo/genny"
	"github.com/gobuffalo/packr/v2"
	"github.com/pkg/errors"
)

func archivedAssets(opts *Options) (*genny.Generator, error) {
	g := genny.New()

	if err := opts.Validate(); err != nil {
		return g, errors.WithStack(err)
	}

	app := opts.App

	outputDir := filepath.Dir(filepath.Join(app.Root, app.Bin))
	target := filepath.Join(outputDir, "assets.zip")
	source := filepath.Join(app.Root, "public", "assets")

	g.RunFn(func(r *genny.Runner) error {
		bb := &bytes.Buffer{}
		archive := zip.NewWriter(bb)
		defer archive.Close()

		box := packr.New("buffalo:build:assets", "")
		box.ResolutionDir = source
		err := box.Walk(func(path string, file packr.File) error {
			info, err := file.FileInfo()
			if err != nil {
				return errors.WithStack(err)
			}

			header, err := zip.FileInfoHeader(info)
			if err != nil {
				return errors.WithStack(err)
			}

			var baseDir string
			if info.IsDir() {
				baseDir = filepath.Base(source)
			}
			if baseDir != "" {
				rel, err := filepath.Rel(source, path)
				if err != nil {
					return errors.WithStack(err)
				}
				header.Name = filepath.Join(baseDir, rel)
			}

			if info.IsDir() {
				header.Name += "/"
			} else {
				header.Method = zip.Deflate
			}

			writer, err := archive.CreateHeader(header)
			if err != nil {
				return errors.WithStack(err)
			}

			if info.IsDir() {
				return nil
			}

			if _, err = io.Copy(writer, file); err != nil {
				return errors.WithStack(err)
			}
			return r.File(genny.NewFile(target, bb))
		})
		if err != nil {
			return errors.WithStack(err)
		}
		return nil
	})

	g.RunFn(func(r *genny.Runner) error {
		f, err := r.FindFile("actions/app.go")
		if err != nil {
			return errors.WithStack(err)
		}
		opts.rollback.Store(f.Name(), f.String())
		body := strings.Replace(f.String(), `app.ServeFiles("/assets"`, `// app.ServeFiles("/assets"`, 1)
		body = strings.Replace(body, `app.ServeFiles("/"`, `// app.ServeFiles("/"`, 1)
		return r.File(genny.NewFileS(f.Name(), body))
	})

	return g, nil
}
