package form

import "github.com/gobuffalo/tags"

func (f Form) InputTag(opts tags.Options) *tags.Tag {
	if opts["type"] == nil {
		opts["type"] = "text"
	}
	return tags.New("input", opts)
}