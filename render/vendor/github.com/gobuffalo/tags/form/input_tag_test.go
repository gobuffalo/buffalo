package form_test

import (
	"testing"

	"github.com/gobuffalo/tags"
	"github.com/gobuffalo/tags/form"
	"github.com/stretchr/testify/require"
)

func Test_Form_InputTag(t *testing.T) {
	r := require.New(t)
	f := form.New(tags.Options{})
	i := f.InputTag(tags.Options{})
	r.Equal(`<input type="text" />`, i.String())
}
