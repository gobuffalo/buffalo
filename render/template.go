package render

import (
	"fmt"
	"html/template"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"

	"github.com/sirupsen/logrus"
	"golang.org/x/text/language"
)

type templateRenderer struct {
	*Engine
	contentType string
	names       []string
	aliases     sync.Map
}

func (s *templateRenderer) ContentType() string {
	return s.contentType
}

func (s *templateRenderer) resolve(name string) ([]byte, error) {
	if s.TemplatesFS == nil {
		return nil, fmt.Errorf("no templates fs defined")
	}

	f, err := s.TemplatesFS.Open(name)
	if err == nil {
		return io.ReadAll(f)
	}

	v, ok := s.aliases.Load(name)
	if !ok {
		return nil, fmt.Errorf("could not find template %s", name)
	}

	f, err = s.TemplatesFS.Open(v.(string))
	if err != nil {
		return nil, err
	}
	return io.ReadAll(f)
}

func (s *templateRenderer) Render(w io.Writer, data Data) error {

	if err := s.updateAliases(); err != nil {
		return err
	}

	var body template.HTML
	for _, name := range s.names {
		var err error
		body, err = s.exec(name, data)
		if err != nil {
			return fmt.Errorf("%s: %w", name, err)
		}
		data["yield"] = body
	}
	_, err := w.Write([]byte(body))
	return err
}

func (s *templateRenderer) updateAliases() error {
	if s.TemplatesFS == nil {
		return nil
	}

	return fs.WalkDir(s.TemplatesFS, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}

		// index.plush.ko-kr.html as index.ko-kr.html
		shortcut := strings.Replace(path, ".plush.", ".", 1)
		s.aliases.Store(shortcut, path)

		// register short version (lang only) of shortcut
		words := strings.Split(filepath.Base(shortcut), ".")
		if len(words) > 2 {
			for _, w := range words[1 : len(words)-1] {
				if tag, err := language.Parse(w); err == nil {
					base, confidence := tag.Base()
					if confidence == language.Exact || confidence == language.High {
						// index.plush.ko-kr.html as index.ko.html
						shortcut := strings.Replace(shortcut, w, base.String(), 1)
						s.aliases.Store(shortcut, path)

						// index.plush.ko-kr.html as index.plush.ko.html
						shortcut = strings.Replace(path, w, base.String(), 1)
						s.aliases.Store(shortcut, path)
					}
				}

			}
		}

		return nil
	})
}

func fixExtension(name string, ct string) string {
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
	return name
}

// partialFeeder returns template string for the name from `TemplateBox`.
// It should be registered as helper named `partialFeeder` so plush can
// find it with the name.
func (s *templateRenderer) partialFeeder(name string) (string, error) {
	ct := strings.ToLower(s.contentType)

	d, f := filepath.Split(name)
	name = filepath.Join(d, "_"+f)
	name = fixExtension(name, ct)

	b, err := s.resolve(name)
	return string(b), err
}

func (s *templateRenderer) exec(name string, data Data) (template.HTML, error) {
	ct := strings.ToLower(s.contentType)
	data["contentType"] = ct

	name = fixExtension(name, ct)

	source, err := s.localizedResolve(name, data)
	if err != nil {
		return "", err
	}

	helpers := map[string]interface{}{}

	for k, v := range s.Helpers {
		helpers[k] = v
	}

	// Allows to specify custom partialFeeder
	if helpers["partialFeeder"] == nil {
		helpers["partialFeeder"] = s.partialFeeder
	}

	helpers = s.addAssetsHelpers(helpers)

	body := string(source)
	for _, ext := range s.exts(name) {
		te, ok := s.TemplateEngines[ext]
		if !ok {
			logrus.Errorf("could not find a template engine for %s", ext)
			continue
		}
		body, err = te(body, data, helpers)
		if err != nil {
			return "", err
		}
	}

	return template.HTML(body), nil
}

// next step, deprecate if this is no longer required.
func (s *templateRenderer) localizedName(name string, data Data) string {
	templateName := name

	languages, ok := data["languages"].([]string)
	if !ok || len(languages) == 0 {
		return templateName
	}

	ll := len(languages)
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
		if _, err := s.resolve(candidateName); err == nil {
			// Replace name with the existing suffixed version
			templateName = candidateName
			break
		}
	}

	return templateName
}

func (s *templateRenderer) localizedResolve(name string, data Data) ([]byte, error) {
	languages, ok := data["languages"].([]string)
	if !ok || len(languages) == 0 {
		return s.resolve(name)
	}

	defaultLang := languages[len(languages)-1] // default language
	ext := filepath.Ext(name)
	rawName := strings.TrimSuffix(name, ext)

	for _, lang := range languages {
		if lang == defaultLang {
			break
		}

		fullLower := strings.ToLower(lang)        // forms of ko-kr or ko
		short := strings.Split(fullLower, "-")[0] // form of ko

		fullLocale := rawName + "." + fullLower + ext
		if source, err := s.resolve(fullLocale); err == nil {
			return source, nil
		}

		if fullLower != short {
			langOnly := rawName + "." + short + ext
			if source, err := s.resolve(langOnly); err == nil {
				return source, nil
			}
		}
	}

	return s.resolve(name)
}

func (s *templateRenderer) exts(name string) []string {
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

func (s *templateRenderer) assetPath(file string) (string, error) {

	if len(assetMap.Keys()) == 0 || os.Getenv("GO_ENV") != "production" {
		manifest, err := s.AssetsFS.Open("manifest.json")
		if err != nil {
			manifest, err = s.AssetsFS.Open("assets/manifest.json")
			if err != nil {
				return assetPathFor(file), nil
			}
		}
		defer manifest.Close()

		err = loadManifest(manifest)
		if err != nil {
			return assetPathFor(file), fmt.Errorf("your manifest.json is not correct: %s", err)
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
	return &templateRenderer{
		Engine:      e,
		contentType: c,
		names:       names,
	}
}
