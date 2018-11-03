package newapp

import (
	"testing"

	"github.com/gobuffalo/meta"
	"github.com/gobuffalo/packr"
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
	Templates = packr.NewBox(".")

	g := Generator{}
	err := g.Validate()
	r.Error(err)
	r.Equal(ErrTemplatesNotFound, errors.Cause(err))
}

func Test_Validate_InGoPath(t *testing.T) {
	r := require.New(t)

	tests := []struct {
		g   Generator
		err error
	}{
		{
			g: Generator{
				App: meta.App{WithModules: true},
			},
			err: nil,
		},
		{
			g: Generator{
				App: meta.App{WithModules: true, WithDep: true},
			},
			err: errGoModulesWithDep,
		},
	}

	for _, test := range tests {
		r.Equal(test.err, test.err)
	}
}
