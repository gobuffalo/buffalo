package build

import (
	"fmt"
	"html/template"
	"strings"

	"github.com/gobuffalo/genny/v2"
	"github.com/gobuffalo/packd"
	"github.com/gobuffalo/plush"
	"github.com/markbates/safe"
)

// TemplateValidator is given a file and returns an
// effort if there is a template validation error
// with the template
type TemplateValidator func(f genny.File) error

// ValidateTemplates returns a genny.RunFn that will walk the
// given box and run each of the files found through each of the
// template validators
func ValidateTemplates(walk packd.Walker, tvs []TemplateValidator) genny.RunFn {
	if len(tvs) == 0 {
		return func(r *genny.Runner) error {
			return nil
		}
	}
	return func(r *genny.Runner) error {
		var errs []string
		err := packd.SkipWalker(walk, packd.CommonSkipPrefixes, func(path string, file packd.File) error {
			info, err := file.FileInfo()
			if err != nil {
				return err
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
					return err
				}
			}

			return nil
		})
		if err != nil {
			return err
		}
		if len(errs) == 0 {
			return nil
		}
		return fmt.Errorf(strings.Join(errs, "\n"))
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
