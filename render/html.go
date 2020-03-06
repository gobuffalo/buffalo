package render

import (
	"html"

	"github.com/gobuffalo/github_flavored_markdown"
	"github.com/gobuffalo/plush/v4"
)

// HTML renders the named files using the 'text/html'
// content type and the github.com/gobuffalo/plush
// package for templating. If more than 1 file is provided
// the second file will be considered a "layout" file
// and the first file will be the "content" file which will
// be placed into the "layout" using "<%= yield %>".
func HTML(names ...string) Renderer {
	e := New(Options{})
	return e.HTML(names...)
}

// HTML renders the named files using the 'text/html'
// content type and the github.com/gobuffalo/plush
// package for templating. If more than 1 file is provided
// the second file will be considered a "layout" file
// and the first file will be the "content" file which will
// be placed into the "layout" using "<%= yield %>". If no
// second file is provided and an `HTMLLayout` is specified
// in the options, then that layout file will be used
// automatically.
func (e *Engine) HTML(names ...string) Renderer {
	if e.HTMLLayout != "" && len(names) == 1 {
		names = append(names, e.HTMLLayout)
	}
	hr := &templateRenderer{
		Engine:      e,
		contentType: "text/html; charset=utf-8",
		names:       names,
	}
	return hr
}

// MDTemplateEngine runs the input through github flavored markdown before sending it to the Plush engine.
func MDTemplateEngine(input string, data map[string]interface{}, helpers map[string]interface{}) (string, error) {
	if ct, ok := data["contentType"].(string); ok && ct == "text/plain" {
		return plush.BuffaloRenderer(input, data, helpers)
	}
	source := github_flavored_markdown.Markdown([]byte(input))
	source = []byte(html.UnescapeString(string(source)))
	return plush.BuffaloRenderer(string(source), data, helpers)
}
