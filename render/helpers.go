package render

import (
	"html/template"
	"net/http"

	"github.com/gobuffalo/plush"
	"github.com/gobuffalo/tags"
	"github.com/markbates/pop"
	"github.com/pkg/errors"
)

func init() {
	plush.Helpers.Add("paginator", func(pagination *pop.Paginator, opts map[string]interface{}, help plush.HelperContext) (template.HTML, error) {
		if opts["path"] == nil {
			if req, ok := help.Value("request").(*http.Request); ok {
				opts["path"] = req.URL.String()
			}
		}
		t, err := tags.Pagination(pagination, opts)
		if err != nil {
			return "", errors.WithStack(err)
		}
		return t.HTML(), nil
	})
}
