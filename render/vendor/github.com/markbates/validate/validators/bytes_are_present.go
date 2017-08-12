package validators

import (
	"fmt"

	"github.com/markbates/validate"
)

type BytesArePresent struct {
	Name  string
	Field []byte
}

func (v *BytesArePresent) IsValid(errors *validate.Errors) {
	if len(v.Field) == 0 {
		errors.Add(GenerateKey(v.Name), fmt.Sprintf("%s can not be blank.", v.Name))
	}
}
