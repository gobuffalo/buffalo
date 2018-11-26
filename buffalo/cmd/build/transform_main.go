package build

import (
	"bytes"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

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

func (b *Builder) buildVersion(version string) string {
	if b.Options.VCS == "git" {
		_, err := exec.LookPath("git")
		if err != nil {
			return version
		}
		cmd := exec.Command("git", "rev-parse", "--short", "HEAD")
		out := &bytes.Buffer{}
		cmd.Stdout = out
		cmd.Stderr = os.Stderr
		cmd.Stdin = os.Stdin
		err = cmd.Run()
		if err == nil && out.String() != "" {
			version = strings.TrimSpace(out.String())
		}
	} else if b.Options.VCS == "bzr" {
		_, err := exec.LookPath("bzr")
		if err != nil {
			return version
		}
		cmd := exec.Command("bzr", "revno")
		out := &bytes.Buffer{}
		cmd.Stdout = out
		cmd.Stderr = os.Stderr
		cmd.Stdin = os.Stdin
		err = cmd.Run()
		if err == nil && out.String() != "" {
			version = strings.TrimSpace(out.String())
		}
	}
	return version
}

func (b *Builder) createBuildMain() error {
	ctx := plush.NewContext()
	ctx.Set("opts", b.Options)

	bt := time.Now().Format(time.RFC3339)
	ctx.Set("buildTime", bt)
	ctx.Set("buildVersion", b.buildVersion(bt))

	t, err := templates.FindString("main.go.tmpl")
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
