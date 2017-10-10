package meta

import (
	"encoding/json"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/gobuffalo/envy"
)

type App struct {
	Pwd         string `json:"pwd"`
	Root        string `json:"root"`
	GoPath      string `json:"go_path"`
	Name        Name   `json:"name"`
	Bin         string `json:"bin"`
	PackagePath string `json:"package_path"`
	ActionsPath string `json:"actions_path"`
	ModelsPath  string `json:"models_path"`
	WithPop     bool   `json:"with_pop"`
	WithDep     bool   `json:"with_dep"`
	WithWebpack bool   `json:"with_webpack"`
	WithYarn    bool   `json:"with_yarn"`
	WithDocker  bool   `json:"with_docker"`
}

func New(root string) App {
	pwd, _ := os.Getwd()
	pp := packagePath(root)

	app := App{
		Pwd:         pwd,
		Root:        root,
		GoPath:      envy.GoPath(),
		Name:        Name(filepath.Base(root)),
		PackagePath: pp,
		ActionsPath: pp + "/actions",
		ModelsPath:  pp + "/models",
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

	return app
}

func (a App) String() string {
	b, _ := json.Marshal(a)
	return string(b)
}

func packagePath(root string) string {
	src := filepath.ToSlash(filepath.Join(envy.GoPath(), "src"))
	root = filepath.ToSlash(root)
	return strings.Replace(root, src+"/", "", 2)
}
