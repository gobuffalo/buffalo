package build

import (
	"strings"
	"sync"

	"github.com/gobuffalo/genny/v2"
)

func transformMain(opts *Options) genny.RunFn {
	if opts.rollback == nil {
		opts.rollback = &sync.Map{}
	}
	return func(r *genny.Runner) error {
		f, err := r.FindFile("main.go")
		if err != nil {
			return err
		}
		opts.rollback.Store(f.Name(), f.String())
		s := strings.Replace(f.String(), "func main()", "func originalMain()", -1)
		f = genny.NewFile(f.Name(), strings.NewReader(s))
		return r.File(f)
	}
}
