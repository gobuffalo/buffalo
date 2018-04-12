package main

import (
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/gobuffalo/buffalo/buffalo/cmd"
	"github.com/gobuffalo/buffalo/meta"
	"github.com/gobuffalo/plush"
	"github.com/pkg/errors"
)

func main() {
	if exists(".buffalo.dev.yml") {
		vBuffalo()
		return
	}
	cmd.Execute()
}

var app meta.App

// lol! it's "vBuffalo"!
func vBuffalo() {
	app = meta.New(".")
	// not using dep or there isn't a vendor folder
	if !app.WithDep || !exists("vendor") {
		cmd.Execute()
		return
	}
	if err := execute(); err != nil {
		log.Fatal(err)
	}
}

const cmdPkg = "github.com/gobuffalo/buffalo/buffalo/cmd"

var pwd, _ = os.Getwd()
var exePath = filepath.Join(pwd, ".grifter", "main.go")
var vBin = filepath.Join(pwd, "bin", "vbuffalo")

func execute() (err error) {
	dir := filepath.Dir(exePath)
	err = os.MkdirAll(dir, 0755)
	if err != nil {
		return errors.WithStack(err)
	}
	defer func() {
		if err = os.RemoveAll(dir); err != nil {
			log.Println(err)
		}
	}()

	if err := depEnsure(); err != nil {
		return errors.WithStack(err)
	}

	if err = os.Symlink(filepath.Join(pwd, "vendor"), filepath.Join(filepath.Dir(exePath), "vendor")); err != nil {
		return errors.WithStack(err)
	}

	if err = writeMain(); err != nil {
		return errors.WithStack(err)
	}

	err = cd(filepath.Dir(exePath), func() error {
		args := []string{"build", "-v"}
		if app.WithSQLite {
			args = append(args, "--tags", "sqlite")
		}
		args = append(args, "-o", vBin)
		cmd := exec.Command("go", args...)
		return run(cmd)
	})

	cmd := exec.Command(vBin, os.Args[1:]...)
	return run(cmd)
}

func depEnsure() error {
	// toml := filepath.Join(pwd, "Gopkg.toml")
	// b, err := ioutil.ReadFile(toml)
	// if err != nil {
	// 	return errors.WithStack(err)
	// }
	return nil
	// if bytes.Contains(b, []byte("github.com/gobuffalo/buffalo/...")) {
	// 	return nil
	// }
	//
	// f, err := os.Create(toml)
	// if err != nil {
	// 	return errors.WithStack(err)
	// }
	// f.WriteString("required = [\"github.com/gobuffalo/buffalo/...\"]\n")
	// f.Write(b)
	// if err := f.Close(); err != nil {
	// 	return errors.WithStack(err)
	// }
	//
	// return run(exec.Command("dep", "ensure", "-v"))
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

func writeMain() error {
	f, err := os.Create(exePath)
	if err != nil {
		return errors.WithStack(err)
	}
	defer f.Close()
	s, err := plush.Render(mainTemplate, plush.NewContextWith(map[string]interface{}{
		"app":    app,
		"cmdPkg": cmdPkg,
	}))
	if err != nil {
		return errors.WithStack(err)
	}
	_, err = f.WriteString(s)
	return err
}

func exists(f string) bool {
	_, err := os.Stat(f)
	return err == nil
}

const mainTemplate = `package main

import (
	"fmt"
	"<%= cmdPkg %>"
)

func main() {
	fmt.Println(cmd.Version)
	cmd.Execute()
}
`

const prune = `  [[prune.project]]
    name = "github.com/gobuffalo/buffalo"
    go-tests = true
    non-go = false
    unused-packages = false
`
