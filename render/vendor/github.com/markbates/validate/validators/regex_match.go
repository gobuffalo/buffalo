package validators

import (
	"fmt"
	"regexp"

	"github.com/markbates/validate"
)

// RegexMatch specifies the properties needed by the validation.
type RegexMatch struct {
	Name  string
	Field string
	Expr  string
}

// IsValid performs the validation based on the regexp match.
func (v *RegexMatch) IsValid(errors *validate.Errors) {
	r := regexp.MustCompile(v.Expr)
	if !r.Match([]byte(v.Field)) {
		errors.Add(GenerateKey(v.Name), fmt.Sprintf("%s does not match the expected format.", v.Name))
	}
}
