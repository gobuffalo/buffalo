package validators_test

import (
	"testing"

	"github.com/markbates/going/validate"
	"github.com/markbates/going/validate/validators"
	"github.com/stretchr/testify/require"
)

func Test_FuncValidator(t *testing.T) {
	r := require.New(t)

	fv := &validators.FuncValidator{
		Field:   "Name",
		Message: "%s can't be blank",
		Fn: func() bool {
			return false
		},
	}

	verrs := validate.NewErrors()
	fv.IsValid(verrs)

	r.Equal([]string{"Name can't be blank"}, verrs.Get("name"))
}
