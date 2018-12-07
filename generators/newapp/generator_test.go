package newapp

import (
	"testing"

	"github.com/gobuffalo/packr/v2"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
)

func Test_Validate_TemplatesFound(t *testing.T) {
	r := require.New(t)

	g := Generator{}
	err := g.Validate()
	r.Error(err)
	r.NotEqual(ErrTemplatesNotFound, errors.Cause(err))
}

func Test_Validate_TemplatesMissing(t *testing.T) {
	r := require.New(t)

	obox := Templates
	defer func() {
		Templates = obox
	}()
	Templates = packr.New(".", ".")

	g := Generator{}
	err := g.Validate()
	r.Error(err)
	r.Equal(ErrTemplatesNotFound, errors.Cause(err))
}
