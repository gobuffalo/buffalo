package vbuffalo

import (
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"

	"html/template"

	"github.com/gobuffalo/buffalo/generators/newapp"
	"github.com/pkg/errors"
)

func depEnsure() error {
	toml := filepath.Join(pwd, "Gopkg.toml")
	b, err := ioutil.ReadFile(toml)
	if err != nil {
		return errors.WithStack(err)
	}

	addPrune := !bytes.Contains(b, []byte("[prune]"))

	make := func() error {
		f, err := os.Create(toml)
		if err != nil {
			return errors.WithStack(err)
		}
		defer f.Close()
		f.Write(b)

		t, err := template.New("toml template").Parse(newapp.GopkgTomlTmpl)
		if err != nil {
			return errors.WithStack(err)
		}

		err = t.Execute(f, map[string]interface{}{
			"opts":     app,
			"addPrune": addPrune,
		})
		if err != nil {
			return errors.WithStack(err)
		}

		return run("dep", []string{"ensure", "-v"})
	}

	if addPrune {
		if err := make(); err != nil {
			return errors.WithStack(err)
		}
	}
	if !bytes.Contains(b, []byte("[[prune.project]] # buffalo")) {
		if err := make(); err != nil {
			return errors.WithStack(err)
		}
	}
	return nil
}
