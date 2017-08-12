package form

import "github.com/gobuffalo/tags"

func (f Form) TextArea(opts tags.Options) *tags.Tag {
	if opts["value"] != nil {
		opts["body"] = opts["value"]
		delete(opts, "value")
	}
	return tags.New("textarea", opts)
}
