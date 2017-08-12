// Package importgraphutil augments "golang.org/x/tools/refactor/importgraph" with a way to build graphs ignoring tests.
package importgraphutil

import (
	"go/build"
	"sync"

	"golang.org/x/tools/go/buildutil"
	"golang.org/x/tools/refactor/importgraph"
)

// BuildNoTests is like "golang.org/x/tools/refactor/importgraph".Build but doesn't consider test imports.
func BuildNoTests(ctxt *build.Context) (forward, reverse importgraph.Graph, errors map[string]error) {
	type importEdge struct {
		from, to string
	}
	type pathError struct {
		path string
		err  error
	}

	ch := make(chan interface{})

	var wg sync.WaitGroup
	buildutil.ForEachPackage(ctxt, func(path string, err error) {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err != nil {
				ch <- pathError{path, err}
				return
			}
			bp, err := ctxt.Import(path, "", 0)
			if _, ok := err.(*build.NoGoError); ok {
				return // empty directory is not an error
			}
			if err != nil {
				ch <- pathError{path, err}
				return
			}
			for _, imp := range bp.Imports {
				ch <- importEdge{path, imp}
			}
			// Ignore test imports.
		}()
	})
	go func() {
		wg.Wait()
		close(ch)
	}()

	forward = make(importgraph.Graph)
	reverse = make(importgraph.Graph)

	for e := range ch {
		switch e := e.(type) {
		case pathError:
			if errors == nil {
				errors = make(map[string]error)
			}
			errors[e.path] = e.err

		case importEdge:
			if e.to == "C" {
				continue // "C" is fake
			}
			addEdge(forward, e.from, e.to)
			addEdge(reverse, e.to, e.from)
		}
	}

	return forward, reverse, errors
}

func addEdge(g importgraph.Graph, from, to string) {
	edges := g[from]
	if edges == nil {
		edges = make(map[string]bool)
		g[from] = edges
	}
	edges[to] = true
}
