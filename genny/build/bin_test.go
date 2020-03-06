package build

import (
	"os/exec"
	"runtime"
	"strings"
	"testing"

	"github.com/gobuffalo/envy"
	"github.com/gobuffalo/meta"
	"github.com/stretchr/testify/require"
)

func Test_buildCmd(t *testing.T) {
	envy.Set("GO_BIN", "go")
	r := require.New(t)

	eq := func(s string, c *exec.Cmd) {
		if runtime.GOOS == "windows" {
			s = strings.Replace(s, "bin/build", `bin\build.exe`, 1)
			s = strings.Replace(s, "bin/foo", `bin\foo.exe`, 1)
		}
		r.Equal(s, strings.Join(c.Args, " "))
	}

	opts := &Options{
		App:       meta.New("."),
		GoCommand: "build",
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

func Test_buildCmd_Unix_RemovesExe(t *testing.T) {
	if runtime.GOOS == "windows" {
		return
	}
	envy.Set("GO_BIN", "go")
	gomods.Force(true)
	r := require.New(t)

	eq := func(s string, c *exec.Cmd) {
		r.Equal(s, strings.Join(c.Args, " "))
	}
	app := meta.New(".")
	app.Bin = "bin/build.exe"
	opts := &Options{
		App:       app,
		GoCommand: "build",
	}
	c, err := buildCmd(opts)
	r.NoError(err)
	eq("go build -o bin/build", c)
}

func Test_buildCmd_Windows_AddsExe(t *testing.T) {
	if runtime.GOOS != "windows" {
		return
	}
	envy.Set("GO_BIN", "go")
	gomods.Force(true)
	r := require.New(t)

	eq := func(s string, c *exec.Cmd) {
		r.Equal(s, strings.Join(c.Args, " "))
	}

	app := meta.New(".")
	for _, x := range []string{"bin\\build", "bin\\build.exe"} {
		app.Bin = x
		opts := &Options{
			App: app,
		}
		c, err := buildCmd(opts)
		r.NoError(err)
		eq("go build -o bin\\build.exe", c)
	}
}

func Test_installCmd(t *testing.T) {
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
		App:       meta.New("."),
		GoCommand: "install",
	}
	c, err := buildCmd(opts)
	r.NoError(err)
	eq("go install", c)

	opts.Environment = "bar"
	c, err = buildCmd(opts)
	r.NoError(err)
	eq("go install -tags bar", c)

	opts.App.Bin = "bin/foo"
	c, err = buildCmd(opts)
	r.NoError(err)
	eq("go install -tags bar", c)

	opts.WithSQLite = true
	c, err = buildCmd(opts)
	r.NoError(err)
	eq("go install -tags bar sqlite", c)

	opts.LDFlags = "-X foo.Bar=baz"
	c, err = buildCmd(opts)
	r.NoError(err)
	eq("go install -tags bar sqlite -ldflags -X foo.Bar=baz", c)

	opts.Static = true
	c, err = buildCmd(opts)
	r.NoError(err)
	eq("go install -tags bar sqlite -ldflags -linkmode external -extldflags \"-static\" -X foo.Bar=baz", c)

	opts.LDFlags = "-X main.BuildTime=asdf"
	c, err = buildCmd(opts)
	r.NoError(err)
	eq("go install -tags bar sqlite -ldflags -linkmode external -extldflags \"-static\" -X main.BuildTime=asdf", c)

	opts.LDFlags = "-X main.BuildVersion=asdf"
	c, err = buildCmd(opts)
	r.NoError(err)
	eq("go install -tags bar sqlite -ldflags -linkmode external -extldflags \"-static\" -X main.BuildVersion=asdf", c)
}
