package build

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/gobuffalo/plush"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

func (b *Builder) prepAPackage() error {
	a := filepath.Join(b.Root, "a")
	logrus.Debugf("preparing %s", a)
	err := os.MkdirAll(a, 0766)
	if err != nil {
		return errors.WithStack(err)
	}

	infl := filepath.Join(b.Root, "inflections.json")
	if _, err := os.Stat(infl); err == nil {
		logrus.Debugf("preparing %s", infl)
		// there's an inflection file we need to copy it over
		fa, err := os.Open(infl)
		if err != nil {
			return errors.WithStack(err)
		}
		defer fa.Close()
		fb, err := os.Create(filepath.Join(b.Root, "a", "inflections.json"))
		if err != nil {
			return errors.WithStack(err)
		}
		defer fb.Close()
		_, err = io.Copy(fb, fa)
		if err != nil {
			return errors.WithStack(err)
		}
	}

	b.cleanups = append(b.cleanups, func() error {
		return os.RemoveAll(a)
	})
	return nil
}

func (b *Builder) buildAInit() error {
	a := filepath.Join(b.Root, "a", "a.go")
	logrus.Debugf("generating %s", a)
	f, err := os.Create(a)
	if err != nil {
		return errors.WithStack(err)
	}
	defer f.Close()

	ctx := plush.NewContext()
	ctx.Set("opts", b.Options)

	t, err := templates.MustString("a.go.tmpl")
	if err != nil {
		return errors.WithStack(err)
	}

	s, err := plush.Render(t, ctx)
	if err != nil {
		return errors.WithStack(err)
	}

	_, err = f.WriteString(s)
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func (b *Builder) buildADatabase() error {
	ad := filepath.Join(b.Root, "a", "database.go")
	logrus.Debugf("generating %s", ad)

	dgo, err := os.Create(ad)
	if err != nil {
		return errors.WithStack(err)
	}
	defer dgo.Close()

	bb := &bytes.Buffer{}
	if b.WithPop {
		// copy the database.yml file to the migrations folder so it's available through packr
		os.MkdirAll(filepath.Join(b.Root, "migrations"), 0755)
		d, err := os.Open("database.yml")
		if err != nil {
			return errors.WithStack(err)
		}
		defer d.Close()
		_, err = io.Copy(bb, d)
		if err != nil {
			return errors.WithStack(err)
		}
		if bytes.Contains(bb.Bytes(), []byte("sqlite")) {
			b.Tags = append(b.Tags, "sqlite")
			if !b.Static {
				logrus.Debug("you are building a SQLite application, please consider using the `--static` flag to compile a static binary")
			}
		}
	}
	dgo.WriteString("package a\n")
	dgo.WriteString(fmt.Sprintf("var DB_CONFIG = `%s`", bb.String()))
	return nil
}
