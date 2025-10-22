package render

import (
	"io/fs"

	"github.com/gobuffalo/plush/v5/helpers/hctx"
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

	// TemplateMetadataKeys allows users to specify custom keys for template metadata
	// If nil, uses default Buffalo metadata approach
	TemplateMetadataKeys map[string]string

	TemplateBaseDir string
}

// Default metadata keys
var defaultTemplateMetadataKeys = map[string]string{
	"template_file": "_buffalo_template_file",
	"base_name":     "_buffalo_base_name",
	"extension":     "_buffalo_extension",
	"last_modified": "_buffalo_last_modification",
}
