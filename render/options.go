package render

import (
	"github.com/gobuffalo/helpers/hctx"
	"github.com/gobuffalo/packd"
)

// Helpers to be included in all templates
type Helpers hctx.Map

// Options for render.Engine
type Options struct {
	// HTMLLayout is the default layout to be used with all HTML renders.
	HTMLLayout string

	// JavaScriptLayout is the default layout to be used with all JavaScript renders.
	JavaScriptLayout string

	// TemplatesBox is the location of the templates directory on disk.
	TemplatesBox packd.Box

	// AssetsBox is the location of the public assets the app will serve.
	AssetsBox packd.Box

	// Helpers to be rendered with the templates
	Helpers Helpers

	// TemplateEngine to be used for rendering HTML templates
	TemplateEngines map[string]TemplateEngine

	// DefaultContentType instructs the engine what it should fall back to if
	// the "content-type" is unknown
	DefaultContentType string
}
