package middleware

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_maskSecrets(t *testing.T) {
	r := require.New(t)
	pl := parameterLogger{}

	filteredForm := pl.maskSecrets(url.Values{
		"FirstName":            []string{"Antonio"},
		"MiddleName":           []string{"José"},
		"LastName":             []string{"Pagano"},
		"Password":             []string{"Secret!"},
		"password":             []string{"Other"},
		"pAssWorD":             []string{"Weird one"},
		"PasswordConfirmation": []string{"Secret!"},

		"SomeCVC": []string{"Untouched"},
	})

	r.Equal(filteredForm.Get("Password"), filteredIndicator[0])
	r.Equal(filteredForm.Get("password"), filteredIndicator[0])
	r.Equal(filteredForm.Get("pAssWorD"), filteredIndicator[0])
	r.Equal(filteredForm.Get("PasswordConfirmation"), filteredIndicator[0])
	r.Equal(filteredForm.Get("LastName"), "Pagano")
	r.Equal(filteredForm.Get("SomeCVC"), "Untouched")
}

func Test_maskSecretsCustom(t *testing.T) {
	r := require.New(t)
	ParameterFilterBlackList = []string{
		"FirstName", "LastName", "MiddleName",
	}

	pl := parameterLogger{}

	filteredForm := pl.maskSecrets(url.Values{
		"FirstName":            []string{"Antonio"},
		"MiddleName":           []string{"José"},
		"LastName":             []string{"Pagano"},
		"Password":             []string{"Secret!"},
		"password":             []string{"Other"},
		"pAssWorD":             []string{"Weird one"},
		"PasswordConfirmation": []string{"Secret!"},

		"SomeCVC": []string{"Untouched"},
	})

	r.Equal(filteredForm.Get("Password"), "Secret!")
	r.Equal(filteredForm.Get("password"), "Other")
	r.Equal(filteredForm.Get("LastName"), filteredIndicator[0])
	r.Equal(filteredForm.Get("SomeCVC"), "Untouched")
}
