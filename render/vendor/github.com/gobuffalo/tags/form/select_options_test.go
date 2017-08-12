package form_test

import (
	"testing"

	"github.com/gobuffalo/tags/form"
	"github.com/stretchr/testify/require"
)

func Test_SelectionOption_String(t *testing.T) {
	r := require.New(t)
	so := form.SelectOption{
		Value: 1,
		Label: "one",
	}
	r.Equal(`<option value="1">one</option>`, so.String())
}

func Test_SelectionOption_Selected_String(t *testing.T) {
	r := require.New(t)
	so := form.SelectOption{
		Value:         1,
		Label:         "one",
		SelectedValue: "1",
	}
	r.Equal(`<option value="1" selected>one</option>`, so.String())
}
