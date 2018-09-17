package fix

import "github.com/gobuffalo/packr/builder"

func PackrClean(r *Runner) error {
	builder.Clean(r.App.Root)
	return nil
}
