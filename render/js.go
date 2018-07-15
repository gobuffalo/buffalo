package render

import (
	"html/template"
	"strings"

	"github.com/gobuffalo/plush"
	"github.com/pkg/errors"
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
// It also implements a new `partial` helper that will run non-JS partials
// through `JSEscapeString` before injecting.
func JSTemplateEngine(input string, data map[string]interface{}, helpers map[string]interface{}) (string, error) {
	var pf partFunc
	var ok bool
	if pf, ok = helpers["partial"].(func(string, Data) (template.HTML, error)); !ok {
		return "", errors.New("could not find a partial function")
	}

	helpers["partial"] = func(name string, dd Data) (template.HTML, error) {
		if strings.Contains(name, ".js") {
			return pf(name, dd)
		}
		h, err := pf(name, dd)
		if err != nil {
			return "", errors.WithStack(err)
		}
		he := template.JSEscapeString(string(h))
		return template.HTML(he), nil
	}

	return plush.BuffaloRenderer(input, data, helpers)
}

type partFunc func(string, Data) (template.HTML, error)
