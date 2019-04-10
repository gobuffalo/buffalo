package render

import (
	"html/template"
	"net/http"

	"github.com/gobuffalo/plush"
	"github.com/gobuffalo/tags"
)

func init() {
	plush.Helpers.Add("paginator", func(pagination interface{}, opts map[string]interface{}, help plush.HelperContext) (template.HTML, error) {
		if opts["path"] == nil {
			if req, ok := help.Value("request").(*http.Request); ok {
				opts["path"] = req.URL.String()
			}
		}
		t, err := tags.Pagination(pagination, opts)
		if err != nil {
			return "", err
		}
		return t.HTML(), nil
	})
}
