package build

import (
	"testing"

	"github.com/gobuffalo/genny/gentest"
	"github.com/gobuffalo/packr"
	"github.com/stretchr/testify/require"
)

func Test_TemplateValidator_Good(t *testing.T) {
	r := require.New(t)

	box := packr.NewBox("../build/_fixtures/template_validator/good")
	tvs := []TemplateValidator{PlushValidator}

	run := gentest.NewRunner()
	run.WithRun(ValidateTemplates(box, tvs))

	r.NoError(run.Run())
}

func Test_TemplateValidator_Bad(t *testing.T) {
	r := require.New(t)

	box := packr.NewBox("../build/_fixtures/template_validator/bad")
	tvs := []TemplateValidator{PlushValidator}

	run := gentest.NewRunner()
	run.WithRun(ValidateTemplates(box, tvs))

	err := run.Run()
	r.Error(err)
	r.Equal("template error in file a.html: line 1: no prefix parse function for > found\ntemplate error in file b.md: line 1: no prefix parse function for > found", err.Error())
}
