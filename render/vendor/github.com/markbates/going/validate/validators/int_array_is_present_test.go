package validators_test

import (
	"testing"

	"github.com/markbates/going/validate"
	. "github.com/markbates/going/validate/validators"
	"github.com/stretchr/testify/assert"
)

func Test_IntArrayIsPresent(t *testing.T) {
	assert := assert.New(t)

	v := IntArrayIsPresent{"Name", []int{1}}
	errors := validate.NewErrors()
	v.IsValid(errors)
	assert.Equal(errors.Count(), 0)

	v = IntArrayIsPresent{"Name", []int{}}
	v.IsValid(errors)
	assert.Equal(errors.Count(), 1)
	assert.Equal(errors.Get("name"), []string{"Name can not be empty."})
}
