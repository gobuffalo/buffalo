package render

func HTML(names ...string) Renderer {
	e := New(Options{})
	return e.HTML(names...)
}

func (e *Engine) HTML(names ...string) Renderer {
	if e.HTMLLayout != "" {
		names = append(names, e.HTMLLayout)
	}
	hr := &templateFileRenderer{
		Engine:      e,
		contentType: "text/html",
		names:       names,
		helpers:     e.TemplateFuncs,
	}
	return hr
}
