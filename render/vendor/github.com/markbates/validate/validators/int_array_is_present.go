package validators

import (
	"fmt"

	"github.com/markbates/validate"
)

type IntArrayIsPresent struct {
	Name  string
	Field []int
}

func (v *IntArrayIsPresent) IsValid(errors *validate.Errors) {
	if len(v.Field) == 0 {
		errors.Add(GenerateKey(v.Name), fmt.Sprintf("%s can not be empty.", v.Name))
	}
}
