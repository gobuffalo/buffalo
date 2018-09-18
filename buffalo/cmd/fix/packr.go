package fix

import "github.com/gobuffalo/packr/builder"

// PackrClean will remove any packr files
func PackrClean(r *Runner) error {
	builder.Clean(r.App.Root)
	return nil
}
