// +build dev

package gfmstyle

import (
	"go/build"
	"log"
	"net/http"
)

func importPathToDir(importPath string) string {
	p, err := build.Import(importPath, "", build.FindOnly)
	if err != nil {
		log.Fatalln(err)
	}
	return p.Dir
}

// Assets contains the gfm.css style file for rendering GitHub Flavored Markdown.
var Assets = http.Dir(importPathToDir("github.com/shurcooL/github_flavored_markdown/gfmstyle/_data"))
