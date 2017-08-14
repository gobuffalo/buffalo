package render

import (
	"encoding/json"
	"fmt"
	"html"
	"html/template"
	"io"
	"log"
	"path/filepath"
	"strings"
	"github.com/gobuffalo/tags"
	// this blank import is here because dep doesn't
	// handle transitive dependencies correctly
	_ "github.com/russross/blackfriday"
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

func (s templateRenderer) partial(name string, dd Data) (template.HTML, error) {
	d, f := filepath.Split(name)
	name = filepath.Join(d, "_"+f)
	return s.exec(name, dd)
}

func (s templateRenderer) exec(name string, data Data) (template.HTML, error) {
	source, err := s.TemplatesBox.MustBytes(name)
	if err != nil {
		return "", err
	}

	helpers := map[string]interface{}{
		"partial": s.partial,
	}

	helpers = s.addAssetsHelpers(helpers)

	for k, v := range s.Helpers {
		helpers[k] = v
	}

	if strings.ToLower(filepath.Ext(name)) == ".md" && strings.ToLower(s.contentType) != "text/plain" {
		source = github_flavored_markdown.Markdown(source)
		source = []byte(html.UnescapeString(string(source)))
	}

	body, err := s.TemplateEngine(string(source), data, helpers)
	if err != nil {
		return "", err
	}

	return template.HTML(body), nil
}

func (s templateRenderer) addAssetsHelpers(helpers map[string]interface{}) map[string]interface{} {
	helpers["assetPath"] = func(file string) template.HTML {
		return template.HTML(s.assetPath(file))
	}

	helpers["javascriptTag"] = func(file string, options tags.Options) template.HTML {
		if options == nil {
			options = tags.Options{}
		}

		if options["type"] == nil {
			options["type"] = "text/javascript"
		}

		options["src"] = s.assetPath(file)
		jsTag := tags.New("script", options)

		return jsTag.HTML()
	}

	helpers["stylesheetTag"] = func(file string, options tags.Options) template.HTML {
		if options == nil {
			options = tags.Options{}
		}

		if options["rel"] == nil {
			options["rel"] = "stylesheet"
		}

		if options["media"] == nil {
			options["media"] = "screen"
		}

		options["href"] = s.assetPath(file)
		cssTag := tags.New("link", options)

		return cssTag.HTML()
	}

	return helpers
}

func (s templateRenderer) assetPath(file string) string {
	manifest, err := s.AssetsBox.MustString("manifest.json")
	if err != nil {
		log.Println("[INFO] didn't find manifest, using raw path to assets")

		return fmt.Sprintf("assets/%v", file)
	}

	var manifestData map[string]string
	err = json.Unmarshal([]byte(manifest), &manifestData)

	if err != nil {
		log.Println("[Warning] seems your manifest is not correct")
		return ""
	}

	if file == "application.css" {
		file = "main.css"
	}

	if file == "application.js" {
		file = "main.js"
	}

	return fmt.Sprintf("/assets/%v", manifestData[file])
}

// Template renders the named files using the specified
// content type and the github.com/gobuffalo/plush
// package for templating. If more than 1 file is provided
// the second file will be considered a "layout" file
// and the first file will be the "content" file which will
// be placed into the "layout" using "{{yield}}".
func Template(c string, names ...string) Renderer {
	e := New(Options{})
	return e.Template(c, names...)
}

// Template renders the named files using the specified
// content type and the github.com/gobuffalo/plush
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
