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
		return c.Redirect(http.StatusFound, u)
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
			S: http.StatusPermanentRedirect,
		},
		{
			S: http.StatusInternalServerError,
		},
	}

	for _, tt := range table {
		a := New(Options{})
		a.GET("/foo/{bar}", func(c Context) error {
			return c.Render(http.StatusOK, render.String(c.Param("bar")))
		})
		a.GET("/", func(c Context) error {
			return c.Redirect(http.StatusPermanentRedirect, "fooBarPath()", tt.I)
		})
		a.GET("/nomap", func(c Context) error {
			return c.Redirect(http.StatusPermanentRedirect, "rootPath()")
		})

		w := httptest.New(a)
		res := w.HTML("/").Get()
		r.Equal(tt.S, res.Code)
		r.Equal(tt.E, res.Location())

		res = w.HTML("/nomap").Get()
		r.Equal(http.StatusPermanentRedirect, res.Code)
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

	r.Equal(http.StatusOK, res.Code)
	r.Equal("Mark", name)
}

func Test_DefaultContext_Param_Multiple(t *testing.T) {
	r := require.New(t)

	app := New(Options{})
	var params ParamValues
	var param string
	app.POST("/{id}", func(c Context) error {
		params = c.Params()
		param = c.Param("id")
		return nil
	})

	w := httptest.New(app)
	res := w.HTML("/a?id=c&y=z&id=d").Post(map[string]string{
		"id": "b",
	})
	paramsExpected := url.Values{
		"id": []string{"a", "b", "c", "d"},
		"y":  []string{"z"},
	}

	r.Equal(200, res.Code)
	r.Equal(paramsExpected, params.(url.Values))
	r.Equal("a", param)
}

func Test_DefaultContext_GetSet(t *testing.T) {
	r := require.New(t)
	c := basicContext()
	r.Nil(c.Value("name"))

	c.Set("name", "Mark")
	r.NotNil(c.Value("name"))
	r.Equal("Mark", c.Value("name").(string))
}

func Test_DefaultContext_Set_not_configured(t *testing.T) {
	r := require.New(t)
	c := DefaultContext{}

	c.Set("name", "Yonghwan")
	r.NotNil(c.Value("name"))
	r.Equal("Yonghwan", c.Value("name").(string))
}

func Test_DefaultContext_Value(t *testing.T) {
	r := require.New(t)
	c := basicContext()
	r.Nil(c.Value("name"))

	c.Set("name", "Mark")
	r.NotNil(c.Value("name"))
	r.Equal("Mark", c.Value("name").(string))
}

func Test_DefaultContext_Value_not_configured(t *testing.T) {
	r := require.New(t)
	c := DefaultContext{}
	r.Nil(c.Value("name"))
}

func Test_DefaultContext_Render(t *testing.T) {
	r := require.New(t)

	c := basicContext()
	res := httptest.NewRecorder()
	c.response = res
	c.params = url.Values{"name": []string{"Mark"}}
	c.Set("greet", "Hello")

	err := c.Render(http.StatusTeapot, render.String(`<%= greet %> <%= params["name"] %>!`))
	r.NoError(err)

	r.Equal(http.StatusTeapot, res.Code)
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
		return c.Render(http.StatusCreated, nil)
	})

	w := httptest.New(a)
	uv := url.Values{"first_name": []string{"Mark"}}
	res := w.HTML("/").Post(uv)
	r.Equal(http.StatusCreated, res.Code)

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
			return c.Error(http.StatusUnprocessableEntity, err)
		}
		return c.Render(http.StatusCreated, nil)
	})

	bb := &bytes.Buffer{}
	req, err := http.NewRequest("POST", "/", bb)
	r.NoError(err)
	req.Header.Del("Content-Type")
	res := httptest.NewRecorder()
	a.ServeHTTP(res, req)
	r.Equal(http.StatusUnprocessableEntity, res.Code)
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
			return c.Error(http.StatusUnprocessableEntity, err)
		}
		return c.Render(http.StatusCreated, nil)
	})

	bb := &bytes.Buffer{}
	req, err := http.NewRequest("POST", "/", bb)
	r.NoError(err)
	// Want to make sure that an empty string value does not cause an error on `split`
	req.Header.Set("Content-Type", "")
	res := httptest.NewRecorder()
	a.ServeHTTP(res, req)
	r.Equal(http.StatusUnprocessableEntity, res.Code)
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
		return c.Render(http.StatusCreated, nil)
	})

	w := httptest.New(a)
	uv := url.Values{"first_name": []string{""}}
	res := w.HTML("/").Post(uv)
	r.Equal(http.StatusCreated, res.Code)

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
		return c.Render(http.StatusCreated, nil)
	})

	w := httptest.New(a)
	res := w.JSON("/").Post(map[string]string{
		"first_name": "Mark",
	})
	r.Equal(http.StatusCreated, res.Code)

	r.Equal("Mark", user.FirstName)
}

func Test_DefaultContext_Data(t *testing.T) {
	r := require.New(t)
	c := basicContext()

	r.EqualValues(map[string]interface{}{}, c.Data())
}

func Test_DefaultContext_Data_not_configured(t *testing.T) {
	r := require.New(t)
	c := DefaultContext{}

	r.EqualValues(map[string]interface{}{}, c.Data())
}

func Test_DefaultContext_String(t *testing.T) {
	r := require.New(t)
	c := basicContext()
	c.Set("name", "Buffalo")
	c.Set("language", "go")

	r.EqualValues("language: go\n\nname: Buffalo", c.String())
}

func Test_DefaultContext_String_EmptyData(t *testing.T) {
	r := require.New(t)
	c := basicContext()
	r.EqualValues("", c.String())
}

func Test_DefaultContext_String_EmptyData_not_configured(t *testing.T) {
	r := require.New(t)
	c := DefaultContext{}

	r.EqualValues("", c.String())
}

func Test_DefaultContext_MarshalJSON(t *testing.T) {
	r := require.New(t)
	c := basicContext()
	c.Set("name", "Buffalo")
	c.Set("language", "go")

	jb, err := c.MarshalJSON()
	r.NoError(err)
	r.EqualValues(`{"language":"go","name":"Buffalo"}`, string(jb))
}

func Test_DefaultContext_MarshalJSON_EmptyData(t *testing.T) {
	r := require.New(t)
	c := basicContext()

	jb, err := c.MarshalJSON()
	r.NoError(err)
	r.EqualValues(`{}`, string(jb))
}

func Test_DefaultContext_MarshalJSON_EmptyData_not_configured(t *testing.T) {
	r := require.New(t)
	c := DefaultContext{}

	jb, err := c.MarshalJSON()
	r.NoError(err)
	r.EqualValues(`{}`, string(jb))
}
