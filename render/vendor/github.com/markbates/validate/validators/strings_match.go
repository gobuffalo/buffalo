package validators

import (
	"fmt"
	"github.com/markbates/validate"
	"strings"
)

type StringsMatch struct {
	Name    string
	Field   string
	Field2  string
	Message string
}

// IsValid performs the validation equality of two strings.
func (v *StringsMatch) IsValid(errors *validate.Errors) {
	if strings.TrimSpace(v.Field) != strings.TrimSpace(v.Field2) {
		if v.Message == "" {
			v.Message = fmt.Sprintf("%s does not equal %s.", v.Field, v.Field2)
		}
		errors.Add(GenerateKey(v.Name), v.Message)
	}
}
