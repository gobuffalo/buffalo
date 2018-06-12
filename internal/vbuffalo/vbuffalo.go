package vbuffalo

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/gobuffalo/buffalo/meta"
	"github.com/pkg/errors"
)

var pwd, _ = os.Getwd()
var mainPath = filepath.Join(pwd, ".grifter", "main.go")
var binPath = filepath.Join(pwd, "bin", "vbuffalo")
var app meta.App

func init() {
	fmt.Println("vbuffalo")
	if runtime.GOOS == "windows" {
		binPath += ".exe"
	}
}

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
		return run("go", args)
	})
	if err != nil {
		return errors.WithStack(err)
	}

	return run(binPath, os.Args[1:])
}

func cd(dir string, fn func() error) error {
	defer os.Chdir(pwd)
	os.Chdir(dir)
	return fn()
}

func run(name string, args []string) error {
	// fmt.Println("vbuffalo", name, strings.Join(args, " "))
	cmd := exec.Command(name, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Start(); err != nil {
		return errors.WithStack(err)
	}
	return cmd.Wait()
}

func exists(f string) bool {
	_, err := os.Stat(f)
	return err == nil
}
