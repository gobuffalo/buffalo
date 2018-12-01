package build

import (
	"os/exec"
	"runtime"
	"strings"
	"testing"

	"github.com/gobuffalo/envy"
	"github.com/gobuffalo/genny/movinglater/gotools/gomods"
	"github.com/gobuffalo/meta"
	"github.com/stretchr/testify/require"
)

func Test_buildCmd(t *testing.T) {
	envy.Set("GO_BIN", "go")
	gomods.Force(true)
	r := require.New(t)

	eq := func(s string, c *exec.Cmd) {
		if runtime.GOOS == "windows" {
			s = strings.Replace(s, "bin/build", `bin\build.exe`, 1)
			s = strings.Replace(s, "bin/foo", `bin\foo.exe`, 1)
		}
		r.Equal(s, strings.Join(c.Args, " "))
	}

	opts := &Options{
		App: meta.New("."),
	}
	c, err := buildCmd(opts)
	r.NoError(err)
	eq("go build -o bin/build", c)

	opts.Environment = "bar"
	c, err = buildCmd(opts)
	r.NoError(err)
	eq("go build -tags bar -o bin/build", c)

	opts.App.Bin = "bin/foo"
	c, err = buildCmd(opts)
	r.NoError(err)
	eq("go build -tags bar -o bin/foo", c)

	opts.WithSQLite = true
	c, err = buildCmd(opts)
	r.NoError(err)
	eq("go build -tags bar sqlite -o bin/foo", c)

	opts.LDFlags = "-X foo.Bar=baz"
	c, err = buildCmd(opts)
	r.NoError(err)
	eq("go build -tags bar sqlite -o bin/foo -ldflags -X foo.Bar=baz", c)

	opts.Static = true
	c, err = buildCmd(opts)
	r.NoError(err)
	eq("go build -tags bar sqlite -o bin/foo -ldflags -linkmode external -extldflags \"-static\" -X foo.Bar=baz", c)

	opts.LDFlags = "-X main.BuildTime=asdf"
	c, err = buildCmd(opts)
	r.NoError(err)
	eq("go build -tags bar sqlite -o bin/foo -ldflags -linkmode external -extldflags \"-static\" -X main.BuildTime=asdf", c)

	opts.LDFlags = "-X main.BuildVersion=asdf"
	c, err = buildCmd(opts)
	r.NoError(err)
	eq("go build -tags bar sqlite -o bin/foo -ldflags -linkmode external -extldflags \"-static\" -X main.BuildVersion=asdf", c)
}
