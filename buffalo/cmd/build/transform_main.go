package build

import (
	"bytes"
	"io"
	"os"
	"path/filepath"

	"github.com/gobuffalo/plush"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

func (b *Builder) transformMain() error {
	logrus.Debug("transforming main() to originalMain()")

	return b.transform("main.go", func(body []byte, w io.Writer) error {
		body = bytes.Replace(body, []byte("func main()"), []byte("func originalMain()"), 1)
		_, err := w.Write(body)
		if err != nil {
			return errors.WithStack(err)
		}
		return nil
	})
}

func (b *Builder) createBuildMain() error {
	pkg := packagePath(b.Root)

	ctx := plush.NewContext()
	ctx.Set("root", b.Root)
	ctx.Set("hasDB", b.HasDB)
	if b.HasDB {
		ctx.Set("modelsPack", pkg+"/models")
	}
	if _, err := os.Stat(filepath.Join(b.Root, "grifts")); err == nil {
		ctx.Set("griftsPack", pkg+"/grifts")
	}
	ctx.Set("aPack", pkg+"/a")
	ctx.Set("name", filepath.Base(b.Root))

	t, err := templates.MustString("main.go.tmpl")
	if err != nil {
		return errors.WithStack(err)
	}

	s, err := plush.Render(t, ctx)
	if err != nil {
		return errors.WithStack(err)
	}

	bbm := filepath.Join(b.Root, "buffalo_build_main.go")
	logrus.Debugf("creating %s", bbm)
	f, err := os.Create(bbm)
	if err != nil {
		return errors.WithStack(err)
	}
	defer f.Close()
	b.cleanups = append(b.cleanups, func() error {
		return os.RemoveAll(bbm)
	})
	f.WriteString(s)
	return nil
}
