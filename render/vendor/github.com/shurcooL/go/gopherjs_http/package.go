package gopherjs_http

import (
	"go/build"
	"net/http"
	"os"
)

// Package returns an http.FileSystem that contains a single file at root,
// containing result of building package with importPath using GopherJS.
func Package(importPath string) http.FileSystem {
	return packageFS{importPath: importPath}
}

type packageFS struct {
	importPath string
}

func (fs packageFS) Open(name string) (http.File, error) {
	if name != "/" {
		return nil, &os.PathError{Op: "open", Path: name, Err: os.ErrNotExist}
	}
	p, err := build.Import(fs.importPath, "", build.FindOnly)
	if err != nil {
		return nil, &os.PathError{Op: `"go/build".Import`, Path: fs.importPath, Err: err}
	}
	return (&gopherJSFS{source: http.Dir(p.SrcRoot)}).compileGoPackage(p.ImportPath)
}
