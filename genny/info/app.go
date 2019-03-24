package info

import (
	"fmt"
	"reflect"
	"text/tabwriter"

	"github.com/gobuffalo/clara/genny/rx"
	"github.com/gobuffalo/genny"
)

func appDetails(opts *Options) genny.RunFn {
	return func(r *genny.Runner) error {
		rx.Header(opts.Out, "Buffalo: Application Details")
		rv := reflect.ValueOf(opts.App)
		rt := rv.Type()

		w := tabwriter.NewWriter(opts.Out, 0, 0, 1, ' ', 0)
		defer w.Flush()
		var err error
		for i := 0; i < rt.NumField(); i++ {
			f := rt.Field(i)
			if !rv.FieldByName(f.Name).CanInterface() {
				continue
			}
			m := fmt.Sprintf("%s\t%v\n", f.Name, rv.FieldByName(f.Name).Interface())
			_, err = w.Write([]byte(m))
			if err != nil {
				return err
			}
		}

		return nil
	}
}
