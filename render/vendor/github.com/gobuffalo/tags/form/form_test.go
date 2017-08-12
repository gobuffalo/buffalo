package form_test

import (
	"testing"

	"github.com/gobuffalo/tags"
	"github.com/gobuffalo/tags/form"
	"github.com/stretchr/testify/require"
)

func Test_NewForm(t *testing.T) {
	r := require.New(t)

	f := form.New(tags.Options{
		"action": "/users/1",
	})
	r.Equal("form", f.Name)
	r.Equal(`<form action="/users/1" method="POST"></form>`, f.String())
}

func Test_NewForm_With_AuthenticityToken(t *testing.T) {
	r := require.New(t)

	f := form.New(tags.Options{
		"action": "/users/1",
	})
	f.SetAuthenticityToken("12345")
	r.Equal("form", f.Name)
	r.Equal(`<form action="/users/1" method="POST"><input name="authenticity_token" type="hidden" value="12345" /></form>`, f.String())
}

func Test_NewForm_With_NotPostMethod(t *testing.T) {
	r := require.New(t)

	f := form.New(tags.Options{
		"action": "/users/1",
		"method": "put",
	})
	r.Equal("form", f.Name)
	r.Equal(`<form action="/users/1" method="POST"><input name="_method" type="hidden" value="PUT" /></form>`, f.String())
}

func Test_Form_Label(t *testing.T) {
	r := require.New(t)
	f := form.New(tags.Options{})
	l := f.Label("Name", tags.Options{})
	r.Equal(`<label>Name</label>`, l.String())
}
