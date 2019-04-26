package render

import (
	"encoding/json"
	"html/template"
	"path/filepath"
	"sync"

	"github.com/gobuffalo/tags"
)

var assetsMutex = &sync.RWMutex{}
var assetMap map[string]string

func loadManifest(manifest string) error {
	assetsMutex.Lock()
	defer assetsMutex.Unlock()

	err := json.Unmarshal([]byte(manifest), &assetMap)
	return err
}

func assetPathFor(file string) string {
	assetsMutex.RLock()
	filePath := assetMap[file]
	assetsMutex.RUnlock()
	if filePath == "" {
		filePath = file
	}
	return filepath.ToSlash(filepath.Join("/assets", filePath))
}

type helperTag struct {
	name string
	fn   func(string, tags.Options) template.HTML
}

func (s templateRenderer) addAssetsHelpers(helpers Helpers) Helpers {
	helpers["assetPath"] = s.assetPath

	ah := []helperTag{
		{"javascriptTag", jsTag},
		{"stylesheetTag", cssTag},
		{"imgTag", imgTag},
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

func jsTag(src string, options tags.Options) template.HTML {
	if options["type"] == nil {
		options["type"] = "text/javascript"
	}

	options["src"] = src
	jsTag := tags.New("script", options)

	return jsTag.HTML()
}

func cssTag(href string, options tags.Options) template.HTML {
	if options["rel"] == nil {
		options["rel"] = "stylesheet"
	}

	if options["media"] == nil {
		options["media"] = "screen"
	}

	options["href"] = href
	cssTag := tags.New("link", options)

	return cssTag.HTML()
}

func imgTag(src string, options tags.Options) template.HTML {
	options["src"] = src
	imgTag := tags.New("img", options)

	return imgTag.HTML()
}
