package render

import (
	"encoding/json"
	"html/template"
	"path/filepath"
	"sync"

	"github.com/gobuffalo/tags"
)

var mu = &sync.Mutex{}
var assetMap map[string]string

func loadManifest(manifest string) error {
	mu.Lock()
	defer mu.Unlock()

	err := json.Unmarshal([]byte(manifest), &assetMap)
	return err
}

func assetPathFor(file string) string {
	mu.Lock()
	defer mu.Unlock()

	filePath := assetMap[file]
	if filePath == "" {
		filePath = file
	}

	return filepath.Join("/assets", filePath)
}

func (s templateRenderer) addAssetsHelpers(helpers map[string]interface{}) map[string]interface{} {
	helpers["assetPath"] = func(file string) template.HTML {
		return template.HTML(s.assetPath(file))
	}

	helpers["javascriptTag"] = func(file string, options tags.Options) template.HTML {
		return jsTag(s.assetPath(file), options)
	}

	helpers["stylesheetTag"] = func(file string, options tags.Options) template.HTML {
		return cssTag(s.assetPath(file), options)
	}

	return helpers
}

func jsTag(src string, options tags.Options) template.HTML {

	if options == nil {
		options = tags.Options{}
	}

	if options["type"] == nil {
		options["type"] = "text/javascript"
	}

	options["src"] = src
	jsTag := tags.New("script", options)

	return jsTag.HTML()
}

func cssTag(href string, options tags.Options) template.HTML {
	if options == nil {
		options = tags.Options{}
	}

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
