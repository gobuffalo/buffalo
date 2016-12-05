package render

// HTML renders the named files using the 'text/html'
// content type and the github.com/aymerick/raymond
// package for templating. If more than 1 file is provided
// the second file will be considered a "layout" file
// and the first file will be the "content" file which will
// be placed into the "layout" using "{{yield}}".
func HTML(names ...string) Renderer {
	e := New(Options{})
	return e.HTML(names...)
}

// HTML renders the named files using the 'text/html'
// content type and the github.com/aymerick/raymond
// package for templating. If more than 1 file is provided
// the second file will be considered a "layout" file
// and the first file will be the "content" file which will
// be placed into the "layout" using "{{yield}}". If no
// second file is provided and an `HTMLLayout` is specified
// in the options, then that layout file will be used
// automatically.
func (e *Engine) HTML(names ...string) Renderer {
	if e.HTMLLayout != "" {
		names = append(names, e.HTMLLayout)
	}
	hr := &templateRenderer{
		Engine:      e,
		contentType: "text/html",
		names:       names,
	}
	return hr
}
