package build

import (
	"path/filepath"
	"strings"

	"github.com/gobuffalo/envy"
	"github.com/gobuffalo/packr"
)

var templates = packr.NewBox("./templates")

func packagePath(rootPath string) string {
	gosrcpath := filepath.ToSlash(filepath.Join(goPath(rootPath), "src"))
	rootPath = filepath.ToSlash(rootPath)
	return strings.Replace(rootPath, gosrcpath+"/", "", 2)
}

func goPath(root string) string {
	gpMultiple := envy.GoPaths()
	path := ""

	for i := 0; i < len(gpMultiple); i++ {
		if strings.HasPrefix(root, filepath.Join(gpMultiple[i], "src")) {
			path = gpMultiple[i]
			break
		}
	}
	return path
}
