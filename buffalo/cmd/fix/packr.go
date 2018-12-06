package fix

import "github.com/gobuffalo/packr/v2/packr2/cmd/fix"

// PackrClean will remove any packr files
func PackrClean(r *Runner) error {
	fix.YesToAll = YesToAll
	return fix.Run()
}
