package validators

import (
	"fmt"

	"github.com/markbates/validate"
)

type IntIsPresent struct {
	Name  string
	Field int
}

func (v *IntIsPresent) IsValid(errors *validate.Errors) {
	if v.Field == 0 {
		errors.Add(GenerateKey(v.Name), fmt.Sprintf("%s can not be blank.", v.Name))
	}
}
