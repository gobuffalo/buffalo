package resolvers

import (
	"fmt"
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
	fmt.Printf("### name -> %+v\n", name)
	if g.RecursiveResolver == nil {
		g.RecursiveResolver = &RecursiveResolver{
			Path: filepath.Join(os.Getenv("GOPATH"), "src", g.Path),
		}
	}
	fmt.Printf("### g.RecursiveResolver -> %+v\n", g.RecursiveResolver)
	return g.RecursiveResolver.Read(name)
}

// Resolve will search your entire $GOPATH to find the file
// in question. I wouldn't really recommend using this approach. It's
// very very slow the first you try to find a file, and there is no
// guarantees of finding the right now.
func (g *GoPathResolver) Resolve(name string) (string, error) {
	fmt.Printf("### name -> %+v\n", name)
	if g.RecursiveResolver == nil {
		g.RecursiveResolver = &RecursiveResolver{
			Path: filepath.Join(os.Getenv("GOPATH"), "src", g.Path),
		}
	}
	return g.RecursiveResolver.Resolve(name)
}
