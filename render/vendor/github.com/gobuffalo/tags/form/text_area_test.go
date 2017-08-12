package form_test

import (
	"testing"

	"github.com/gobuffalo/tags"
	"github.com/gobuffalo/tags/form"
	"github.com/markbates/pop/nulls"
	"github.com/stretchr/testify/require"
)

func Test_Form_TextArea(t *testing.T) {
	r := require.New(t)
	f := form.New(tags.Options{})
	ta := f.TextArea(tags.Options{
		"value": "hi",
	})
	r.Equal(`<textarea>hi</textarea>`, ta.String())
}

func Test_Form_TextArea_nullsString(t *testing.T) {
	r := require.New(t)
	f := form.New(tags.Options{})
	ta := f.TextArea(tags.Options{
		"value": nulls.NewString("hi"),
	})
	r.Equal(`<textarea>hi</textarea>`, ta.String())
}

func Test_Form_TextArea_nullsString_empty(t *testing.T) {
	r := require.New(t)
	f := form.New(tags.Options{})
	ta := f.TextArea(tags.Options{
		"value": nulls.String{},
	})
	r.Equal(`<textarea></textarea>`, ta.String())
}
