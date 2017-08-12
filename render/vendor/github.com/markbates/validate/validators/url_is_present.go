package validators

import (
	"fmt"
	"net/url"

	"github.com/markbates/validate"
)

type URLIsPresent struct {
	Name    string
	Field   string
	Message string
}

// IsValid performs the validation to check if URL is formatted correctly
// uses net/url ParseRequestURI to check validity
func (v *URLIsPresent) IsValid(errors *validate.Errors) {
	if v.Field == "http://" || v.Field == "https://" {
		v.Message = fmt.Sprintf("%s url is empty", v.Name)
		errors.Add(GenerateKey(v.Name), v.Message)
	}
	parsedUrl, err := url.ParseRequestURI(v.Field)
	if err != nil {
		if v.Message == "" {
			v.Message = fmt.Sprintf("%s does not match url format. Err: %s", v.Name,
				err)
		}
		errors.Add(GenerateKey(v.Name), v.Message)
	} else {
		if parsedUrl.Scheme != "" && parsedUrl.Scheme != "http" && parsedUrl.Scheme != "https" {
			v.Message = fmt.Sprintf("%s invalid url scheme", v.Name)
			errors.Add(GenerateKey(v.Name), v.Message)
		}
	}
}
