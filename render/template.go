package render

import (
	"html/template"
	"io"
	"path/filepath"
	"strings"

	"github.com/shurcooL/github_flavored_markdown"
)

type templateRenderer struct {
	*Engine
	contentType string
	names       []string
}

func (s templateRenderer) ContentType() string {
	return s.contentType
}

func (s templateRenderer) Render(w io.Writer, data Data) error {
	var body template.HTML
	var err error
	for _, name := range s.names {

		body, err = s.exec(name, data)
		if err != nil {
			return err
		}
		data["yield"] = body
	}
	w.Write([]byte(body))
	return nil
}

func (s templateRenderer) exec(name string, data Data) (template.HTML, error) {
	var body string
	source, err := s.Resolver().Read(filepath.Join(s.TemplatesPath, name))
	if err != nil {
		return "", err
	}

	opts := TemplateOptions{
		Data:    data,
		Helpers: s.Helpers,
	}

	opts.Helpers["partial"] = func(name string) (template.HTML, error) {
		d, f := filepath.Split(name)
		name = filepath.Join(d, "_"+f)
		return s.exec(name, data)
	}

	body, err = s.TemplateEngine(string(source), opts)
	if err != nil {
		return "", err
	}

	if strings.ToLower(filepath.Ext(name)) == ".md" {
		b := github_flavored_markdown.Markdown([]byte(body))
		body = string(b)
	}
	return template.HTML(body), nil
}

// Template renders the named files using the specified
// content type and the github.com/aymerick/raymond
// package for templating. If more than 1 file is provided
// the second file will be considered a "layout" file
// and the first file will be the "content" file which will
// be placed into the "layout" using "{{yield}}".
func Template(c string, names ...string) Renderer {
	e := New(Options{})
	return e.Template(c, names...)
}

// Template renders the named files using the specified
// content type and the github.com/aymerick/raymond
// package for templating. If more than 1 file is provided
// the second file will be considered a "layout" file
// and the first file will be the "content" file which will
// be placed into the "layout" using "{{yield}}".
func (e *Engine) Template(c string, names ...string) Renderer {
	return templateRenderer{
		Engine:      e,
		contentType: c,
		names:       names,
	}
}
