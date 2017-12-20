package meta

import (
	"encoding/json"
	"os"
	"path"
	"path/filepath"
	"runtime"

	"github.com/gobuffalo/envy"
)

// App represents meta data for a Buffalo application on disk
type App struct {
	Pwd         string `json:"pwd"`
	Root        string `json:"root"`
	GoPath      string `json:"go_path"`
	Name        Name   `json:"name"`
	Bin         string `json:"bin"`
	PackagePkg  string `json:"package_path"`
	ActionsPkg  string `json:"actions_path"`
	ModelsPkg   string `json:"models_path"`
	GriftsPkg   string `json:"grifts_path"`
	WithPop     bool   `json:"with_pop"`
	WithDep     bool   `json:"with_dep"`
	WithWebpack bool   `json:"with_webpack"`
	WithYarn    bool   `json:"with_yarn"`
	WithDocker  bool   `json:"with_docker"`
	WithGrifts  bool   `json:"with_grifts"`
}

// New App based on the details found at the provided root path
func New(root string) App {
	pwd, _ := os.Getwd()
	if root == "." {
		root = pwd
	}
	name := Name(filepath.Base(root))
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

	if _, err := os.Stat(filepath.Join(root, "database.yml")); err == nil {
		app.WithPop = true
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

	return app
}

func (a App) String() string {
	b, _ := json.Marshal(a)
	return string(b)
}
