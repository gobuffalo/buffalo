package vbuffalo

import (
	"os"
	"os/exec"
	"path/filepath"

	"github.com/gobuffalo/buffalo/meta"
	"github.com/pkg/errors"
)

const cmdPkg = "github.com/gobuffalo/buffalo/buffalo/cmd"

var pwd, _ = os.Getwd()
var mainPath = filepath.Join(pwd, ".grifter", "main.go")
var binPath = filepath.Join(pwd, "bin", "vbuffalo")
var app meta.App

// Execute using vbuffalo. If this doesn't meet the vbuffalo
// requirements then it should use the passed in function instead.
func Execute(ex func() error) error {
	if !exists(".buffalo.dev.yml") {
		return ex()
	}
	app = meta.New(".")
	// not using dep or there isn't a vendor folder
	if !app.WithDep || !exists("vendor") {
		return ex()
	}
	return execute()
}

func execute() error {
	dir := filepath.Dir(mainPath)
	err := os.MkdirAll(dir, 0755)
	if err != nil {
		return errors.WithStack(err)
	}
	defer os.RemoveAll(dir)

	if err := depEnsure(); err != nil {
		return errors.WithStack(err)
	}

	if err = writeMain(); err != nil {
		return errors.WithStack(err)
	}

	err = cd(filepath.Dir(mainPath), func() error {
		args := []string{"build", "-v"}
		if app.WithSQLite {
			args = append(args, "--tags", "sqlite")
		}
		args = append(args, "-o", binPath)
		cmd := exec.Command("go", args...)
		return run(cmd)
	})

	cmd := exec.Command(binPath, os.Args[1:]...)
	return run(cmd)
}

func cd(dir string, fn func() error) error {
	defer os.Chdir(pwd)
	os.Chdir(dir)
	return fn()
}

func run(cmd *exec.Cmd) error {
	// fmt.Println(strings.Join(cmd.Args, " "))
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func exists(f string) bool {
	_, err := os.Stat(f)
	return err == nil
}
