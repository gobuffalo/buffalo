package render

import (
	"encoding/json"
	"html/template"
	"io"
	"path/filepath"

	ht "github.com/gobuffalo/helpers/tags"
	"github.com/gobuffalo/tags/v3"
)

type helperTag struct {
	name string
	fn   func(string, tags.Options) template.HTML
}

func (s *templateRenderer) addAssetsHelpers(helpers Helpers) Helpers {
	helpers["assetPath"] = s.assetPath

	ah := []helperTag{
		{"javascriptTag", ht.JS},
		{"stylesheetTag", ht.CSS},
		{"imgTag", ht.Img},
	}

	for _, h := range ah {
		func(h helperTag) {
			helpers[h.name] = func(file string, options tags.Options) (template.HTML, error) {
				if options == nil {
					options = tags.Options{}
				}
				f, err := s.assetPath(file)
				if err != nil {
					return "", err
				}
				return h.fn(f, options), nil
			}
		}(h)
	}

	return helpers
}

var assetMap = stringMap{}

func assetPathFor(file string) string {
	filePath, ok := assetMap.Load(file)
	if filePath == "" || !ok {
		filePath = file
	}
	return filepath.ToSlash(filepath.Join("/assets", filePath))
}

func loadManifest(manifest io.Reader) error {
	m := map[string]string{}

	err := json.NewDecoder(manifest).Decode(&m)
	if err != nil {
		return err
	}
	for k, v := range m {
		assetMap.Store(k, v)
	}
	return nil
}
