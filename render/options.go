package render

// Options for render.Engine
type Options struct {
	// HTMLLayout is the default layout to be used with all HTML renders.
	HTMLLayout string
	// TemplatesPath is the location of the templates directory on disk.
	TemplatesPath string
	// TemplateHelpers to be used with all rendered templates.
	// See github.com/aymerick/raymond for more details on helpers.
	TemplateHelpers map[string]interface{}
}
