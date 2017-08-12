package validators

import (
	"fmt"

	"github.com/markbates/validate"
)

type IntIsGreaterThan struct {
	Name     string
	Field    int
	Compared int
}

func (v *IntIsGreaterThan) IsValid(errors *validate.Errors) {
	if !(v.Field > v.Compared) {
		errors.Add(GenerateKey(v.Name), fmt.Sprintf("%d is not greater than %d.", v.Field, v.Compared))
	}
}
