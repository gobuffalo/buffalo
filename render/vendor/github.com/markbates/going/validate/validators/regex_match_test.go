package validators_test

import (
	"testing"

	"github.com/markbates/going/validate"
	. "github.com/markbates/going/validate/validators"
	"github.com/stretchr/testify/assert"
)

func Test_RegexMatch(t *testing.T) {
	assert := assert.New(t)

	v := RegexMatch{"Phone", "555-555-5555", "^([0-9]{3}-[0-9]{3}-[0-9]{4})$"}
	errors := validate.NewErrors()
	v.IsValid(errors)
	assert.Equal(errors.Count(), 0)

	v = RegexMatch{"Phone", "123-ab1-1424", "^([0-9]{3}-[0-9]{3}-[0-9]{4})$"}
	v.IsValid(errors)
	assert.Equal(errors.Count(), 1)
	assert.Equal(errors.Get("phone"), []string{"Phone does not match the expected format."})
}
