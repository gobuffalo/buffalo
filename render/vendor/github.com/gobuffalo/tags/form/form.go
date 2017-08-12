package form

import (
	"strings"

	"github.com/gobuffalo/tags"
)

type Form struct {
	*tags.Tag
}

func (f *Form) SetAuthenticityToken(s string) {
	f.Prepend(tags.New("input", tags.Options{
		"value": s,
		"type":  "hidden",
		"name":  "authenticity_token",
	}))
}

func (f Form) Label(value string, opts tags.Options) *tags.Tag {
	opts["body"] = value
	return tags.New("label", opts)
}

func New(opts tags.Options) *Form {
	if opts["method"] == nil {
		opts["method"] = "POST"
	}

	if opts["multipart"] != nil {
		opts["enctype"] = "multipart/form-data"
		delete(opts, "multipart")
	}

	form := &Form{
		Tag: tags.New("form", opts),
	}

	m := strings.ToUpper(form.Options["method"].(string))
	if m != "POST" && m != "GET" {
		form.Options["method"] = "POST"
		form.Prepend(tags.New("input", tags.Options{
			"value": m,
			"type":  "hidden",
			"name":  "_method",
		}))
	}

	return form
}
