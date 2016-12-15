package render

import "github.com/markbates/buffalo/render/resolvers"

// Options for render.Engine
type Options struct {
	// HTMLLayout is the default layout to be used with all HTML renders.
	HTMLLayout string
	// TemplatesPath is the location of the templates directory on disk.
	TemplatesPath string
	// FileResolver will attempt to file a file and return it's bytes, if possible
	FileResolver resolvers.FileResolver
	// Helpers to be rendered with the templates
	Helpers map[string]interface{}
	// CacheTemplates reduced overheads, but won't reload changed templates.
	// This should only be set to true in production environments.
	CacheTemplates bool
}
