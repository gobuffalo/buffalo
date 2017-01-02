package render

import "github.com/gobuffalo/buffalo/render/resolvers"

// Options for render.Engine
type Options struct {
	// HTMLLayout is the default layout to be used with all HTML renders.
	HTMLLayout string
	// TemplatesPath is the location of the templates directory on disk.
	TemplatesPath string
	// FileResolverFunc will attempt to file a file and return it's bytes, if possible
	FileResolverFunc func() resolvers.FileResolver
	fileResolver     resolvers.FileResolver
	// Helpers to be rendered with the templates
	Helpers map[string]interface{}
	// CacheTemplates reduced overheads, but won't reload changed templates.
	// This should only be set to true in production environments.
	CacheTemplates bool
}

// Resolver calls the FileResolverFunc and returns the resolver. The resolver
// is cached, so the function can be called multiple times without penalty.
// This is necessary because certain resolvers, like the RiceBox one, require
// a fully initialized state to work properly and can not be run directly from
// init functions.
func (o *Options) Resolver() resolvers.FileResolver {
	if o.fileResolver == nil {
		o.fileResolver = o.FileResolverFunc()
	}
	return o.fileResolver
}
