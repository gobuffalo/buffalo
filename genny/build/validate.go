package build

import (
	"bytes"
	"fmt"
	"html/template"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/gobuffalo/genny"
	"github.com/gobuffalo/meta"
	"github.com/gobuffalo/packd"
	"github.com/gobuffalo/plush"
	"github.com/karrick/godirwalk"
	"github.com/markbates/safe"
	"github.com/pkg/errors"
)

// TemplateValidator is given a file and returns an
// effort if there is a template validation error
// with the template
type TemplateValidator func(f genny.File) error

// ValidateTemplates returns a genny.RunFn that will walk the
// given box and run each of the files found through each of the
// template validators
func ValidateTemplates(walk packd.Walkable, tvs []TemplateValidator) genny.RunFn {
	if len(tvs) == 0 {
		return func(r *genny.Runner) error {
			return nil
		}
	}
	return func(r *genny.Runner) error {
		var errs []string
		walk.Walk(func(path string, file packd.File) error {
			info, err := file.FileInfo()
			if err != nil {
				return errors.WithStack(err)
			}
			if info.IsDir() {
				return nil
			}

			f := genny.NewFile(path, file)
			for _, tv := range tvs {
				err := safe.Run(func() {
					if err := tv(f); err != nil {
						errs = append(errs, fmt.Sprintf("template error in file %s: %s", path, err.Error()))
					}
				})
				if err != nil {
					return errors.WithStack(err)
				}
			}

			return nil
		})
		if len(errs) == 0 {
			return nil
		}
		return errors.New(strings.Join(errs, "\n"))
	}
}

// PlushValidator validates the file is a valid
// Plush file if the extension is .md, .html, or .plush
func PlushValidator(f genny.File) error {
	if !genny.HasExt(f, ".html", ".md", ".plush") {
		return nil
	}
	_, err := plush.Parse(f.String())
	return err
}

// GoTemplateValidator validates the file is a
// valid Go text/template file if the extension
// is .tmpl
func GoTemplateValidator(f genny.File) error {
	if !genny.HasExt(f, ".tmpl") {
		return nil
	}
	t := template.New(f.Name())
	_, err := t.Parse(f.String())
	return err
}

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
