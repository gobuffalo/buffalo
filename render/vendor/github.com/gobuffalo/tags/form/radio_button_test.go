package form_test

import (
	"testing"

	"github.com/gobuffalo/tags"
	"github.com/gobuffalo/tags/form"
	"github.com/stretchr/testify/require"
)

func Test_Form_RadioButton(t *testing.T) {
	r := require.New(t)
	f := form.New(tags.Options{})
	ct := f.RadioButton(tags.Options{})
	r.Equal(`<label><input type="radio" checked /> </label>`, ct.String())
}

func Test_Form_RadioButton_WithValue(t *testing.T) {
	r := require.New(t)
	f := form.New(tags.Options{})
	ct := f.RadioButton(tags.Options{
		"value": 1,
	})
	r.Equal(`<label><input type="radio" value="1" /> </label>`, ct.String())
}

func Test_Form_RadioButton_WithValueSelected(t *testing.T) {
	r := require.New(t)
	f := form.New(tags.Options{})
	ct := f.RadioButton(tags.Options{
		"value":   1,
		"checked": "1",
	})
	r.Equal(`<label><input type="radio" value="1" checked /> </label>`, ct.String())
}

func Test_Form_RadioButton_WithLabel(t *testing.T) {
	r := require.New(t)
	f := form.New(tags.Options{})
	ct := f.RadioButton(tags.Options{
		"value": 1,
		"label": "check me",
	})
	r.Equal(`<label><input type="radio" value="1" /> check me</label>`, ct.String())
}
