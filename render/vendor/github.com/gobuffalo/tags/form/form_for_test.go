package form_test

import (
	"testing"
	"time"

	"github.com/gobuffalo/tags"
	"github.com/gobuffalo/tags/form"
	"github.com/markbates/pop/nulls"
	"github.com/stretchr/testify/require"
)

type Talk struct {
	Date time.Time `format:"01-02-2006"`
}

func Test_NewFormFor(t *testing.T) {
	r := require.New(t)

	f := form.NewFormFor(Talk{}, tags.Options{
		"action": "/users/1",
	})
	r.Equal("form", f.Name)
	r.Equal(`<form action="/users/1" id="talk-form" method="POST"></form>`, f.String())
}

func Test_FormFor_InputValue(t *testing.T) {
	r := require.New(t)
	f := form.NewFormFor(Talk{}, tags.Options{
		"action": "/users/1",
	})

	l := f.InputTag("Name", tags.Options{"value": "Something"})

	r.Equal(`<input id="talk-Name" name="Name" type="text" value="Something" />`, l.String())
}

func Test_FormFor_InputValueFormat(t *testing.T) {
	r := require.New(t)
	f := form.NewFormFor(Talk{}, tags.Options{
		"action": "/users/1",
	})

	l := f.InputTag("Date", tags.Options{})
	r.Equal(`<input id="talk-Date" name="Date" type="text" value="01-01-0001" />`, l.String())

	l = f.InputTag("Date", tags.Options{"format": "01/02"})
	r.Equal(`<input id="talk-Date" name="Date" type="text" value="01/01" />`, l.String())
}

func Test_NewFormFor_With_AuthenticityToken(t *testing.T) {
	r := require.New(t)

	f := form.NewFormFor(Talk{}, tags.Options{
		"action": "/users/1",
	})
	f.SetAuthenticityToken("12345")
	r.Equal("form", f.Name)
	r.Equal(`<form action="/users/1" id="talk-form" method="POST"><input name="authenticity_token" type="hidden" value="12345" /></form>`, f.String())
}

func Test_NewFormFor_With_NotPostMethod(t *testing.T) {
	r := require.New(t)

	f := form.NewFormFor(Talk{}, tags.Options{
		"action": "/users/1",
		"method": "put",
	})
	r.Equal("form", f.Name)
	r.Equal(`<form action="/users/1" id="talk-form" method="POST"><input name="_method" type="hidden" value="PUT" /></form>`, f.String())
}

func Test_FormFor_Label(t *testing.T) {
	r := require.New(t)
	f := form.NewFormFor(Talk{}, tags.Options{})
	l := f.Label("Name", tags.Options{})
	r.Equal(`<label>Name</label>`, l.String())
}

func Test_FormFor_FieldDoesntExist(t *testing.T) {
	r := require.New(t)
	f := form.NewFormFor(Talk{}, tags.Options{})
	l := f.InputTag("IDontExist", tags.Options{})
	r.Equal(`<input id="talk-IDontExist" name="IDontExist" type="text" value="" />`, l.String())
}

func Test_FormFor_NullableField(t *testing.T) {
	r := require.New(t)
	model := struct {
		Name       string
		CreditCard nulls.String
		Floater    nulls.Float64
		Other      nulls.Bool
	}{
		CreditCard: nulls.NewString("Hello"),
	}

	f := form.NewFormFor(model, tags.Options{})

	cases := map[string][]string{
		"CreditCard": {`<input id="-CreditCard" name="CreditCard" type="text" value="Hello" />`},
		"Floater":    {`<input id="-Floater" name="Floater" type="text" value="" />`},
		"Other":      {`<input id="-Other" name="Other" type="text" value="" />`},
	}

	for field, html := range cases {
		l := f.InputTag(field, tags.Options{})
		r.Equal(html[0], l.String())
	}
}
