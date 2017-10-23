package render

import "github.com/gobuffalo/packr"

// Helpers to be included in all templates
type Helpers map[string]interface{}

// Options for render.Engine
type Options struct {
	// HTMLLayout is the default layout to be used with all HTML renders.
	HTMLLayout string

	// JavaScriptLayout is the default layout to be used with all JavaScript renders.
	JavaScriptLayout string

	// TemplatesBox is the location of the templates directory on disk.
	TemplatesBox packr.Box

	// AssetsBox is the location of the public assets the app will serve.
	AssetsBox packr.Box

	// Helpers to be rendered with the templates
	Helpers Helpers

	// TemplateEngine to be used for rendering HTML templates
	TemplateEngines map[string]TemplateEngine
}
