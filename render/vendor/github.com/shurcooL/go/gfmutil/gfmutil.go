// Package gfmutil offers functionality to render GitHub Flavored Markdown to io.Writer.
package gfmutil

import (
	"bytes"
	"io"
	"net/http"

	"github.com/shurcooL/github_flavored_markdown"
)

// TODO: Change API to return errors rather than panicking.

// WriteGitHubFlavoredMarkdownViaLocal converts GitHub Flavored Markdown to full HTML page and writes it to w.
// It assumes that GFM CSS is available at /assets/gfm/gfm.css.
func WriteGitHubFlavoredMarkdownViaLocal(w io.Writer, markdown []byte) {
	io.WriteString(w, `<html><head><meta charset="utf-8"><link href="/assets/gfm/gfm.css" media="all" rel="stylesheet" type="text/css" /><link href="//cdnjs.cloudflare.com/ajax/libs/octicons/2.1.2/octicons.css" media="all" rel="stylesheet" type="text/css" /></head><body><article class="markdown-body entry-content" style="padding: 30px;">`)
	w.Write(github_flavored_markdown.Markdown(markdown))
	io.WriteString(w, `</article></body></html>`)
}

// WriteGitHubFlavoredMarkdownViaGitHub converts GitHub Flavored Markdown to full HTML page and writes it to w
// by using GitHub API.
// It assumes that GFM CSS is available at /assets/gfm/gfm.css.
func WriteGitHubFlavoredMarkdownViaGitHub(w io.Writer, markdown []byte) {
	io.WriteString(w, `<html><head><meta charset="utf-8"><link href="/assets/gfm/gfm.css" media="all" rel="stylesheet" type="text/css" /><link href="//cdnjs.cloudflare.com/ajax/libs/octicons/2.1.2/octicons.css" media="all" rel="stylesheet" type="text/css" /></head><body><article class="markdown-body entry-content" style="padding: 30px;">`)

	// Convert GitHub Flavored Markdown to HTML (includes syntax highlighting for diff, Go, etc.)
	resp, err := http.Post("https://api.github.com/markdown/raw", "text/x-markdown", bytes.NewReader(markdown))
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	_, err = io.Copy(w, resp.Body)
	if err != nil {
		panic(err)
	}

	io.WriteString(w, `</article></body></html>`)
}
