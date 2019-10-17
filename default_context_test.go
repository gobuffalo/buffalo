package buffalo

import (
	"bytes"
	"context"
	"net/http"
	"net/url"
	"sync"
	"testing"

	"github.com/gobuffalo/buffalo/render"
	"github.com/gobuffalo/httptest"
	"github.com/gobuffalo/logger"
	"github.com/stretchr/testify/require"
)

func basicContext() DefaultContext {
	return DefaultContext{
		Context: context.Background(),
		logger:  logger.New(logger.DebugLevel),
		data:    &sync.Map{},
		flash:   &Flash{data: make(map[string][]string)},
	}
}

func Test_DefaultContext_Redirect(t *testing.T) {
	r := require.New(t)
	a := New(Options{})
	u := "/foo?bar=http%3A%2F%2Flocalhost%3A3000%2Flogin%2Fcallback%2Ffacebook"
	a.GET("/", func(c Context) error {
		return c.Redirect(302, u)
	})

	w := httptest.New(a)
	res := w.HTML("/").Get()
	r.Equal(u, res.Location())
}

func Test_DefaultContext_Redirect_Helper(t *testing.T) {
	r := require.New(t)

	table := []struct {
		E string
		I map[string]interface{}
		S int
	}{
		{
			E: "/foo/baz/",
			I: map[string]interface{}{"bar": "baz"},
			S: 302,
		},
		{
			S: 500,
		},
	}

	for _, tt := range table {
		a := New(Options{})
		a.GET("/foo/{bar}", func(c Context) error {
			return c.Render(200, render.String(c.Param("bar")))
		})
		a.GET("/", func(c Context) error {
			return c.Redirect(302, "fooPath()", tt.I)
		})
		a.GET("/nomap", func(c Context) error {
			return c.Redirect(302, "rootPath()")
		})

		w := httptest.New(a)
		res := w.HTML("/").Get()
		r.Equal(tt.S, res.Code)
		r.Equal(tt.E, res.Location())

		res = w.HTML("/nomap").Get()
		r.Equal(302, res.Code)
		r.Equal("/", res.Location())
	}
}

func Test_DefaultContext_Param(t *testing.T) {
	r := require.New(t)
	c := DefaultContext{
		params: url.Values{
			"name": []string{"Mark"},
		},
	}

	r.Equal("Mark", c.Param("name"))
}

func Test_DefaultContext_Param_form(t *testing.T) {
	r := require.New(t)

	app := New(Options{})
	var name string
	app.POST("/", func(c Context) error {
		name = c.Param("name")
		return nil
	})

	w := httptest.New(app)
	res := w.HTML("/").Post(map[string]string{
		"name": "Mark",
	})

	r.Equal(200, res.Code)
	r.Equal("Mark", name)
}

func Test_DefaultContext_GetSet(t *testing.T) {
	r := require.New(t)
	c := basicContext()
	r.Nil(c.Value("name"))

	c.Set("name", "Mark")
	r.NotNil(c.Value("name"))
	r.Equal("Mark", c.Value("name").(string))
}

func Test_DefaultContext_Value(t *testing.T) {
	r := require.New(t)
	c := basicContext()
	r.Nil(c.Value("name"))

	c.Set("name", "Mark")
	r.NotNil(c.Value("name"))
	r.Equal("Mark", c.Value("name").(string))
	r.Equal("Mark", c.Value("name").(string))
}

func Test_DefaultContext_Render(t *testing.T) {
	r := require.New(t)

	c := basicContext()
	res := httptest.NewRecorder()
	c.response = res
	c.params = url.Values{"name": []string{"Mark"}}
	c.Set("greet", "Hello")

	err := c.Render(123, render.String(`<%= greet %> <%= params["name"] %>!`))
	r.NoError(err)

	r.Equal(123, res.Code)
	r.Equal("Hello Mark!", res.Body.String())
}

func Test_DefaultContext_Bind_Default(t *testing.T) {
	r := require.New(t)

	user := struct {
		FirstName string `form:"first_name"`
	}{}

	a := New(Options{})
	a.POST("/", func(c Context) error {
		err := c.Bind(&user)
		if err != nil {
			return err
		}
		return c.Render(201, nil)
	})

	w := httptest.New(a)
	uv := url.Values{"first_name": []string{"Mark"}}
	res := w.HTML("/").Post(uv)
	r.Equal(201, res.Code)

	r.Equal("Mark", user.FirstName)
}

func Test_DefaultContext_Bind_No_ContentType(t *testing.T) {
	r := require.New(t)

	user := struct {
		FirstName string `form:"first_name"`
	}{
		FirstName: "Mark",
	}

	a := New(Options{})
	a.POST("/", func(c Context) error {
		err := c.Bind(&user)
		if err != nil {
			return c.Error(422, err)
		}
		return c.Render(201, nil)
	})

	bb := &bytes.Buffer{}
	req, err := http.NewRequest("POST", "/", bb)
	r.NoError(err)
	req.Header.Del("Content-Type")
	res := httptest.NewRecorder()
	a.ServeHTTP(res, req)
	r.Equal(422, res.Code)
	r.Contains(res.Body.String(), "blank content type")
}

func Test_DefaultContext_Bind_Empty_ContentType(t *testing.T) {
	r := require.New(t)

	user := struct {
		FirstName string `form:"first_name"`
	}{
		FirstName: "Mark",
	}

	a := New(Options{})
	a.POST("/", func(c Context) error {
		err := c.Bind(&user)
		if err != nil {
			return c.Error(422, err)
		}
		return c.Render(201, nil)
	})

	bb := &bytes.Buffer{}
	req, err := http.NewRequest("POST", "/", bb)
	r.NoError(err)
	// Want to make sure that an empty string value does not cause an error on `split`
	req.Header.Set("Content-Type", "")
	res := httptest.NewRecorder()
	a.ServeHTTP(res, req)
	r.Equal(422, res.Code)
	r.Contains(res.Body.String(), "blank content type")
}

func Test_DefaultContext_Bind_Default_BlankFields(t *testing.T) {
	r := require.New(t)

	user := struct {
		FirstName string `form:"first_name"`
	}{
		FirstName: "Mark",
	}

	a := New(Options{})
	a.POST("/", func(c Context) error {
		err := c.Bind(&user)
		if err != nil {
			return err
		}
		return c.Render(201, nil)
	})

	w := httptest.New(a)
	uv := url.Values{"first_name": []string{""}}
	res := w.HTML("/").Post(uv)
	r.Equal(201, res.Code)

	r.Equal("", user.FirstName)
}

func Test_DefaultContext_Bind_JSON(t *testing.T) {
	r := require.New(t)

	user := struct {
		FirstName string `json:"first_name"`
	}{}

	a := New(Options{})
	a.POST("/", func(c Context) error {
		err := c.Bind(&user)
		if err != nil {
			return err
		}
		return c.Render(201, nil)
	})

	w := httptest.New(a)
	res := w.JSON("/").Post(map[string]string{
		"first_name": "Mark",
	})
	r.Equal(201, res.Code)

	r.Equal("Mark", user.FirstName)
}
