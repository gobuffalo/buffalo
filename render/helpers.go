package render

import (
	"html/template"

	"github.com/gobuffalo/plush"
	"github.com/gobuffalo/tags"
	"github.com/markbates/pop"
	"github.com/pkg/errors"
)

func init() {
	plush.Helpers.Add("paginator", func(pagination *pop.Paginator, opts map[string]interface{}) (template.HTML, error) {
		t, err := tags.Pagination(pagination, opts)
		if err != nil {
			return "", errors.WithStack(err)
		}
		return t.HTML(), nil
	})
}
