package render

import "github.com/gobuffalo/buffalo/internal/consts"

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
		contentType: consts.MIME_JavaScript,
		names:       names,
	}
	return hr
}
