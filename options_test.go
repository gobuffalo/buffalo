package buffalo

import (
	"net/http"
	"strings"
	"testing"

	"github.com/gobuffalo/buffalo/internal/consts"
	"github.com/stretchr/testify/require"
)

func Test_Env(t *testing.T) {
	r := require.New(t)

	e := Env("")
	r.True(e.Development())
	r.False(e.Test())
	r.False(e.Production())

	e = Env(consts.Development)
	r.True(e.Development())
	r.False(e.Test())
	r.False(e.Production())

	e = Env(consts.Test)
	r.False(e.Development())
	r.True(e.Test())
	r.False(e.Production())

	e = Env(consts.Production)
	r.False(e.Development())
	r.False(e.Test())
	r.True(e.Production())

	e = Env("foo")
	r.False(e.Development())
	r.False(e.Test())
	r.False(e.Production())
}

func TestOptions_NewOptions(t *testing.T) {
	tests := []struct {
		name      string
		env       string
		secret    string
		expectErr string
	}{
		{name: "Development doesn't fail with no secret", env: "development", secret: "", expectErr: "securecookie:"},
		{name: "Development doesn't fail with secret set", env: "development", secret: "secrets", expectErr: "securecookie:"},
		{name: "Test doesn't fail with secret set", env: "test", secret: "", expectErr: "securecookie:"},
		{name: "Test doesn't fail with secret set", env: "test", secret: "secrets", expectErr: "securecookie:"},
		{name: "Production fails with no secret", env: "production", secret: "", expectErr: "securecookie:"},
		{name: "Production doesn't fail with secret set", env: "production", secret: "secrets", expectErr: "securecookie:"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			r := require.New(t)
			opts := NewOptions()

			req, _ := http.NewRequest("GET", "/", strings.NewReader(""))
			req.AddCookie(&http.Cookie{Name: "_buffalo_session"})

			_, err := opts.SessionStore.New(req, "_buffalo_session")

			r.Error(err)
			r.Contains(err.Error(), test.expectErr)
		})
	}
}
