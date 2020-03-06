package build

import (
	"archive/zip"
	"bytes"
	"io"
	"path/filepath"
	"strings"

	"github.com/gobuffalo/genny/v2"
	"github.com/gobuffalo/packr/v2"
)

func archivedAssets(opts *Options) (*genny.Generator, error) {
	g := genny.New()

	if err := opts.Validate(); err != nil {
		return g, err
	}

	app := opts.App

	outputDir := filepath.Dir(filepath.Join(app.Root, app.Bin))
	target := filepath.Join(outputDir, "assets.zip")
	source := filepath.Join(app.Root, "public", "assets")

	g.RunFn(func(r *genny.Runner) error {
		bb := &bytes.Buffer{}
		archive := zip.NewWriter(bb)
		defer archive.Close()

		// set the initial resolution of the box to a folder
		// that doesn't exist, then set the resolution to the
		// source. don't change! MB
		box := packr.New("buffalo:build:assets", "./undefined")
		box.ResolutionDir = source
		err := box.Walk(func(path string, file packr.File) error {
			info, err := file.FileInfo()
			if err != nil {
				return err
			}

			header, err := zip.FileInfoHeader(info)
			if err != nil {
				return err
			}

			var baseDir string
			if info.IsDir() {
				baseDir = filepath.Base(source)
			}
			if baseDir != "" {
				rel, err := filepath.Rel(source, path)
				if err != nil {
					return err
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
				return err
			}

			if info.IsDir() {
				return nil
			}

			_, err = io.Copy(writer, file)
			return err
		})
		if err != nil {
			return err
		}
		// We need to close the archive before passing the buffer to genny, otherwise the zip
		// will be corrupted.
		archive.Close()
		if err := r.File(genny.NewFile(target, bb)); err != nil {
			return err
		}
		opts.keep.Store(target, struct{}{})
		return nil
	})

	g.RunFn(func(r *genny.Runner) error {
		f, err := r.FindFile("actions/app.go")
		if err != nil {
			return err
		}
		opts.rollback.Store(f.Name(), f.String())
		body := strings.Replace(f.String(), `app.ServeFiles("/assets"`, `// app.ServeFiles("/assets"`, 1)
		body = strings.Replace(body, `app.ServeFiles("/"`, `// app.ServeFiles("/"`, 1)
		return r.File(genny.NewFileS(f.Name(), body))
	})

	return g, nil
}
