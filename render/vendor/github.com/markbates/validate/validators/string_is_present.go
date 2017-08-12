package validators

import (
	"fmt"
	"strings"

	"github.com/markbates/validate"
)

type StringIsPresent struct {
	Name  string
	Field string
}

func (v *StringIsPresent) IsValid(errors *validate.Errors) {
	if strings.TrimSpace(v.Field) == "" {
		errors.Add(GenerateKey(v.Name), fmt.Sprintf("%s can not be blank.", v.Name))
	}
}
