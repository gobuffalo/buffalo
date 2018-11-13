package render

import (
	"github.com/gobuffalo/plush"
	"github.com/markbates/oncer"
)

// JavaScript renders the named files using the 'application/javascript'
// content type and the github.com/gobuffalo/plush
// package for templating. If more than 1 file is provided
// the second file will be considered a "layout" file
// and the first file will be the "content" file which will
// be placed into the "layout" using "<%= yield %>".
func JavaScript(names ...string) Renderer {
	e := New(Options{})
	return e.JavaScript(names...)
}

// JavaScript renders the named files using the 'application/javascript'
// content type and the github.com/gobuffalo/plush
// package for templating. If more than 1 file is provided
// the second file will be considered a "layout" file
// and the first file will be the "content" file which will
// be placed into the "layout" using "<%= yield %>". If no
// second file is provided and an `JavaScriptLayout` is specified
// in the options, then that layout file will be used
// automatically.
func (e *Engine) JavaScript(names ...string) Renderer {
	if e.JavaScriptLayout != "" && len(names) == 1 {
		names = append(names, e.JavaScriptLayout)
	}
	hr := &templateRenderer{
		Engine:      e,
		contentType: "application/javascript",
		names:       names,
	}
	return hr
}

// JSTemplateEngine renders files with a `.js` extension through Plush.
// Deprecated: use github.com/gobuffalo/plush.BuffaloRenderer instead.
func JSTemplateEngine(input string, data map[string]interface{}, helpers map[string]interface{}) (string, error) {
	oncer.Deprecate(0, "render.JSTemplateEngine", "Use github.com/gobuffalo/plush.BuffaloRenderer instead.")
	return plush.BuffaloRenderer(input, data, helpers)
}
