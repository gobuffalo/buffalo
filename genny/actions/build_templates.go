package actions

import (
	"fmt"

	"github.com/gobuffalo/genny"
)

func buildTemplates(pres *presenter) genny.RunFn {
	return func(r *genny.Runner) error {
		f, err := box.FindString("view.plush.html.tmpl")
		if err != nil {
			return err
		}
		for _, a := range pres.Actions {
			pres.Data["action"] = a
			fn := fmt.Sprintf("templates/%s/%s.plush.html.tmpl", pres.Name.Folder(), a.File())
			xf := genny.NewFileS(fn, f)
			xf, err = transform(pres, xf)
			if err != nil {
				return err
			}
			if err := r.File(xf); err != nil {
				return err
			}
		}
		return nil
	}
}
