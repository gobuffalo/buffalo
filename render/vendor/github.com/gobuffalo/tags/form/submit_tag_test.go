package form_test

import (
	"testing"

	"github.com/gobuffalo/tags"
	"github.com/gobuffalo/tags/form"
	"github.com/stretchr/testify/require"
)

func Test_Form_SubmitTag(t *testing.T) {
	r := require.New(t)
	f := form.New(tags.Options{})
	i := f.SubmitTag("Send", tags.Options{"class": "btn btn-primary"})

	r.Equal(`<input class="btn btn-primary" type="submit" value="Send" />`, i.String())
}
