package render

import (
	"io/fs"

	"github.com/gobuffalo/helpers/hctx"
)

// Helpers to be included in all templates
type Helpers hctx.Map

// Options for render.Engine
type Options struct {
	// HTMLLayout is the default layout to be used with all HTML renders.
	HTMLLayout string

	// JavaScriptLayout is the default layout to be used with all JavaScript renders.
	JavaScriptLayout string

	// TemplateFS is the fs.FS that holds the templates
	TemplatesFS fs.FS

	// AssetsFS is the fs.FS that holds the of the public assets the app will serve.
	AssetsFS fs.FS

	// Helpers to be rendered with the templates
	Helpers Helpers

	// TemplateEngine to be used for rendering HTML templates
	TemplateEngines map[string]TemplateEngine

	// DefaultContentType instructs the engine what it should fall back to if
	// the "content-type" is unknown
	DefaultContentType string

	// Prefix inherits the global single prefix from buffalo.Options
	Prefix string
}
