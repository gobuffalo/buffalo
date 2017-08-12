package form

import (
	"fmt"
	"html/template"
	"strings"

	"github.com/gobuffalo/tags"
)

func (f Form) RadioButton(opts tags.Options) *tags.Tag {
	opts["type"] = "radio"

	var label string
	if opts["label"] != nil {
		label = fmt.Sprint(opts["label"])
		delete(opts, "label")
	}

	value := opts["value"]
	checked := opts["checked"]
	delete(opts, "checked")
	ct := f.InputTag(opts)
	ct.Checked = template.HTMLEscaper(value) == template.HTMLEscaper(checked)
	tag := tags.New("label", tags.Options{
		"body": strings.Join([]string{ct.String(), label}, " "),
	})
	return tag
}
