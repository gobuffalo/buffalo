package form_test

import (
	"testing"

	"github.com/gobuffalo/tags"
	"github.com/gobuffalo/tags/form"
	"github.com/stretchr/testify/require"
)

func Test_Form_CheckboxTag(t *testing.T) {
	r := require.New(t)
	f := form.New(tags.Options{})
	ct := f.CheckboxTag(tags.Options{"name": "Chubby"})
	r.Equal(`<label><input name="Chubby" type="checkbox" value="true" /></label>`, ct.String())
}

func Test_Form_CheckboxTag_WithValue(t *testing.T) {
	r := require.New(t)
	f := form.New(tags.Options{})
	ct := f.CheckboxTag(tags.Options{
		"value":     1,
		"checked":   "1",
		"unchecked": "2",
		"name":      "Chubby",
	})
	r.Equal(`<label><input name="Chubby" type="checkbox" value="1" checked /><input name="Chubby" type="hidden" value="2" /></label>`, ct.String())
}

func Test_Form_CheckboxTag_WithLabel(t *testing.T) {
	r := require.New(t)
	f := form.New(tags.Options{})
	ct := f.CheckboxTag(tags.Options{
		"label": " check me",
		"name":  "Chubby",
	})
	r.Equal(`<label><input name="Chubby" type="checkbox" value="true" /> check me</label>`, ct.String())
}
