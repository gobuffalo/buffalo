package render

import (
	"html/template"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	// this blank import is here because dep doesn't
	// handle transitive dependencies correctly
	_ "github.com/russross/blackfriday"
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
	ct := strings.ToLower(s.contentType)
	data["contentType"] = ct

	if filepath.Ext(name) == "" {
		switch {
		case strings.Contains(ct, "html"):
			name += ".html"
		case strings.Contains(ct, "javascript"):
			name += ".js"
		case strings.Contains(ct, "markdown"):
			name += ".md"
		}
	}

	// Try to use localized version
	templateName := name
	if languages, ok := data["languages"].([]string); ok {
		ll := len(languages)
		if ll > 0 {
			// Default language is the last in the list
			defaultLanguage := languages[ll-1]
			ext := filepath.Ext(name)
			rawName := strings.TrimSuffix(name, ext)

			for _, l := range languages {
				var candidateName string
				if l == defaultLanguage {
					break
				}
				candidateName = rawName + "." + strings.ToLower(l) + ext
				if s.TemplatesBox.Has(candidateName) {
					// Replace name with the existing suffixed version
					templateName = candidateName
					break
				}
			}
		}
	}

	source, err := s.TemplatesBox.MustBytes(templateName)
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

	body := string(source)
	for _, ext := range s.exts(name) {
		te, ok := s.TemplateEngines[ext]
		if !ok {
			logrus.Errorf("could not find a template engine for %s\n", ext)
			continue
		}
		body, err = te(body, data, helpers)
		if err != nil {
			return "", errors.Wrap(err, name)
		}
	}

	return template.HTML(body), nil
}

func (s templateRenderer) exts(name string) []string {
	exts := []string{}
	for {
		ext := filepath.Ext(name)
		if ext == "" {
			break
		}
		name = strings.TrimSuffix(name, ext)
		exts = append(exts, strings.ToLower(ext[1:]))
	}
	if len(exts) == 0 {
		return []string{"html"}
	}
	sort.Sort(sort.Reverse(sort.StringSlice(exts)))
	return exts
}

func (s templateRenderer) assetPath(file string) (string, error) {

	if len(assetMap) == 0 || os.Getenv("GO_ENV") != "production" {
		manifest, err := s.AssetsBox.MustString("manifest.json")

		if err != nil {
			manifest, err = s.AssetsBox.MustString("assets/manifest.json")
			if err != nil {
				return assetPathFor(file), nil
			}
		}

		err = loadManifest(manifest)
		if err != nil {
			return assetPathFor(file), errors.Wrap(err, "your manifest.json is not correct")
		}
	}

	return assetPathFor(file), nil
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
