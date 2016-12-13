package render

// Options for render.Engine
type Options struct {
	// HTMLLayout is the default layout to be used with all HTML renders.
	HTMLLayout string
	// TemplatesPath is the location of the templates directory on disk.
	TemplatesPath string
	// Helpers to be rendered with the templates
	Helpers map[string]interface{}
	// CacheTemplates reduced overheads, but won't reload changed templates.
	// This should only be set to true in production environments.
	CacheTemplates bool
}
