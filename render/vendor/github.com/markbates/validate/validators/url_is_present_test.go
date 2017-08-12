package validators

import (
	"testing"

	"github.com/markbates/validate"
	"github.com/stretchr/testify/require"
)

func Test_URLIsPresent(t *testing.T) {
	r := require.New(t)

	var tests = []struct {
		url   string
		valid bool
	}{
		{"", false},
		{"http://", false},
		{"https://", false},
		{"http", false},
		{"google.com", false},
		{"http://www.google.com", true},
		{"http://google.com", true},
		{"google.com", false},
		{"https://www.google.cOM", true},
		{"ht123tps://www.google.cOM", false},
		{"https://golang.Org", true},
		{"https://invalid#$%#$@.Org", false},
	}
	for _, test := range tests {
		v := URLIsPresent{Name: "URL", Field: test.url}
		errors := validate.NewErrors()
		v.IsValid(errors)
		r.Equal(test.valid, !errors.HasAny(), test.url, errors.Error())
	}
}
