package buffalo

import (
	"net/http"
	"strings"
	"testing"

	"github.com/gobuffalo/envy"
	"github.com/stretchr/testify/require"
)

func TestOptions_NewOptions(t *testing.T) {
	tests := []struct {
		name      string
		env       string
		secret    string
		expectErr string
	}{
		{name: "Development doesn't fail with no secret", env: "development", secret: "", expectErr: "securecookie: the value is not valid"},
		{name: "Development doesn't fail with secret set", env: "development", secret: "secrets", expectErr: "securecookie: the value is not valid"},
		{name: "Test doesn't fail with secret set", env: "test", secret: "", expectErr: "securecookie: the value is not valid"},
		{name: "Test doesn't fail with secret set", env: "test", secret: "secrets", expectErr: "securecookie: the value is not valid"},
		{name: "Production fails with no secret", env: "production", secret: "", expectErr: "securecookie: hash key is not set"},
		{name: "Production doesn't fail with secret set", env: "production", secret: "secrets", expectErr: "securecookie: the value is not valid"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			r := require.New(t)
			envy.Temp(func() {
				envy.Set("GO_ENV", test.env)
				envy.Set("SESSION_SECRET", test.secret)

				opts := NewOptions()

				req, _ := http.NewRequest("GET", "/", strings.NewReader(""))
				req.AddCookie(&http.Cookie{Name: "_buffalo_session"})

				_, err := opts.SessionStore.New(req, "_buffalo_session")

				r.Error(err)
				r.Equal(test.expectErr, err.Error())
			})
		})
	}
}
