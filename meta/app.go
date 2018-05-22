package meta

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"runtime"

	"github.com/gobuffalo/envy"
	"github.com/markbates/inflect"
)

// App represents meta data for a Buffalo application on disk
type App struct {
	Pwd         string       `json:"pwd"`
	Root        string       `json:"root"`
	GoPath      string       `json:"go_path"`
	Name        inflect.Name `json:"name"`
	Bin         string       `json:"bin"`
	PackagePkg  string       `json:"package_path"`
	ActionsPkg  string       `json:"actions_path"`
	ModelsPkg   string       `json:"models_path"`
	GriftsPkg   string       `json:"grifts_path"`
	VCS         string       `json:"vcs"`
	WithPop     bool         `json:"with_pop"`
	WithSQLite  bool         `json:"with_sqlite"`
	WithDep     bool         `json:"with_dep"`
	WithWebpack bool         `json:"with_webpack"`
	WithYarn    bool         `json:"with_yarn"`
	WithDocker  bool         `json:"with_docker"`
	WithGrifts  bool         `json:"with_grifts"`
}

// New App based on the details found at the provided root path
func New(root string) App {
	pwd, _ := os.Getwd()
	if root == "." {
		root = pwd
	}

	// Handle symlinks
	var oldPwd = pwd
	pwd = ResolveSymlinks(pwd)
	os.Chdir(pwd)
	if runtime.GOOS != "windows" {
		// On Non-Windows OS, os.Getwd() uses PWD env var as a preferred
		// way to get the working dir.
		os.Setenv("PWD", pwd)
	}
	defer func() {
		// Restore PWD
		os.Chdir(oldPwd)
		if runtime.GOOS != "windows" {
			os.Setenv("PWD", oldPwd)
		}
	}()

	// Gather meta data
	name := inflect.Name(filepath.Base(root))
	pp := envy.CurrentPackage()
	if filepath.Base(pp) != string(name) {
		pp = path.Join(pp, string(name))
	}

	app := App{
		Pwd:        pwd,
		Root:       root,
		GoPath:     envy.GoPath(),
		Name:       name,
		PackagePkg: pp,
		ActionsPkg: pp + "/actions",
		ModelsPkg:  pp + "/models",
		GriftsPkg:  pp + "/grifts",
	}

	app.Bin = filepath.Join("bin", filepath.Base(root))

	if runtime.GOOS == "windows" {
		app.Bin += ".exe"
	}
	db := filepath.Join(root, "database.yml")
	if _, err := os.Stat(db); err == nil {
		app.WithPop = true
		if b, err := ioutil.ReadFile(db); err == nil {
			app.WithSQLite = bytes.Contains(bytes.ToLower(b), []byte("sqlite"))
		}
	}
	if _, err := os.Stat(filepath.Join(root, "Gopkg.toml")); err == nil {
		app.WithDep = true
	}
	if _, err := os.Stat(filepath.Join(root, "webpack.config.js")); err == nil {
		app.WithWebpack = true
	}
	if _, err := os.Stat(filepath.Join(root, "yarn.lock")); err == nil {
		app.WithYarn = true
	}
	if _, err := os.Stat(filepath.Join(root, "Dockerfile")); err == nil {
		app.WithDocker = true
	}
	if _, err := os.Stat(filepath.Join(root, "grifts")); err == nil {
		app.WithGrifts = true
	}
	if _, err := os.Stat(filepath.Join(root, ".git")); err == nil {
		app.VCS = "git"
	} else if _, err := os.Stat(filepath.Join(root, ".bzr")); err == nil {
		app.VCS = "bzr"
	}

	return app
}

// ResolveSymlinks takes a path and gets the pointed path
// if the original one is a symlink.
func ResolveSymlinks(p string) string {
	cd, err := os.Lstat(p)
	if err != nil {
		return p
	}
	if cd.Mode()&os.ModeSymlink != 0 {
		// This is a symlink
		r, err := filepath.EvalSymlinks(p)
		if err != nil {
			return p
		}
		return r
	}
	return p
}

func (a App) String() string {
	b, _ := json.Marshal(a)
	return string(b)
}
