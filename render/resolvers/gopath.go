package resolvers

import (
	"os"
	"path/filepath"
)

// GoPathResolver will search your entire $GOPATH to find the file
// in question. I wouldn't really recommend using this approach. It's
// very very slow the first you try to find a file, and there is no
// guarantees of finding the right now.
type GoPathResolver struct {
	Path string
	*RecursiveResolver
}

// Read will search your entire $GOPATH to find the file
// in question. I wouldn't really recommend using this approach. It's
// very very slow the first you try to find a file, and there is no
// guarantees of finding the right now.
func (g *GoPathResolver) Read(name string) ([]byte, error) {
	if g.RecursiveResolver == nil {
		g.RecursiveResolver = &RecursiveResolver{
			Path: filepath.Join(os.Getenv("GOPATH"), "src", g.Path),
		}
	}
	return g.RecursiveResolver.Read(name)
}

// Resolve will search your entire $GOPATH to find the file
// in question. I wouldn't really recommend using this approach. It's
// very very slow the first you try to find a file, and there is no
// guarantees of finding the right now.
func (g *GoPathResolver) Resolve(name string) (string, error) {
	if g.RecursiveResolver == nil {
		g.RecursiveResolver = &RecursiveResolver{
			Path: filepath.Join(os.Getenv("GOPATH"), "src", g.Path),
		}
	}
	return g.RecursiveResolver.Resolve(name)
}
