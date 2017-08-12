package validators_test

import (
	"testing"

	"github.com/markbates/validate"
	"github.com/markbates/validate/validators"
	"github.com/stretchr/testify/require"
)

func Test_FuncValidator(t *testing.T) {
	r := require.New(t)

	fv := &validators.FuncValidator{
		Name:    "Name",
		Field:   "Field",
		Message: "%s is an invalid name",
		Fn: func() bool {
			return false
		},
	}

	verrs := validate.NewErrors()
	fv.IsValid(verrs)

	r.Equal([]string{"Field is an invalid name"}, verrs.Get("name"))
}

func Test_FuncValidatorNoName(t *testing.T) {
	r := require.New(t)

	fv := &validators.FuncValidator{
		Field:   "Name",
		Message: "%s is invalid",
		Fn: func() bool {
			return false
		},
	}

	verrs := validate.NewErrors()
	fv.IsValid(verrs)

	r.Equal([]string{"Name is invalid"}, verrs.Get("name"))
}
