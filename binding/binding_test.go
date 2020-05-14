package binding

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/require"
)

type blogPost struct {
	Tags     []string
	Dislikes int
	Likes    int32
}

func Test_Register(t *testing.T) {
	r := require.New(t)

	Register("foo/bar", func(*http.Request, interface{}) error {
		return nil
	})

	r.NotNil(BaseRequestBinder.binders["foo/bar"])

	req, err := http.NewRequest("POST", "/", nil)
	r.NoError(err)

	req.Header.Set("Content-Type", "foo/bar")
	req.Form = url.Values{
		"Tags":     []string{"AAA"},
		"Likes":    []string{"12"},
		"Dislikes": []string{"1000"},
	}

	req.ParseForm()

	var post blogPost
	r.NoError(Exec(req, &post))

	r.Equal([]string(nil), post.Tags)
	r.Equal(int32(0), post.Likes)
	r.Equal(0, post.Dislikes)

}

func Test_RegisterCustomDecoder(t *testing.T) {
	r := require.New(t)

	RegisterCustomDecoder(func(vals []string) (interface{}, error) {
		return []string{"X"}, nil
	}, []interface{}{[]string{}}, nil)

	RegisterCustomDecoder(func(vals []string) (interface{}, error) {
		return 0, nil
	}, []interface{}{int(0)}, nil)

	post := blogPost{}
	req, err := http.NewRequest("POST", "/", nil)
	r.NoError(err)

	req.Header.Set("Content-Type", "application/html")
	req.Form = url.Values{
		"Tags":     []string{"AAA"},
		"Likes":    []string{"12"},
		"Dislikes": []string{"1000"},
	}
	req.ParseForm()

	r.NoError(Exec(req, &post))
	r.Equal([]string{"X"}, post.Tags)
	r.Equal(int32(12), post.Likes)
	r.Equal(0, post.Dislikes)
}
