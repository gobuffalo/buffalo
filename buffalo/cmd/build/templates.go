package build

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/gobuffalo/packr"
	"github.com/gobuffalo/plush"
	"github.com/pkg/errors"
)

var templates = packr.NewBox("./templates")

func (b *Builder) validateTemplates() error {
	if b.SkipTemplateValidation {
		return nil
	}
	errs := []string{}
	err := filepath.Walk(filepath.Join(b.App.Root, "templates"), func(path string, info os.FileInfo, err error) error {
		if info == nil || info.IsDir() {
			return nil
		}

		ext := filepath.Ext(path)
		if ext != ".html" && ext != ".md" {
			return nil
		}

		b, err := ioutil.ReadFile(path)
		if err != nil {
			return errors.WithStack(err)
		}

		if _, err = plush.Parse(string(b)); err != nil {
			errs = append(errs, fmt.Sprintf("template error in file %s: %s", path, err.Error()))
		}

		return nil
	})

	if err != nil {
		return errors.WithStack(err)
	}

	if len(errs) > 0 {
		return errors.New(strings.Join(errs, "\n"))
	}

	return nil
}
