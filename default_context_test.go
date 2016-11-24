package buffalo

import (
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/Sirupsen/logrus"
	"github.com/markbates/buffalo/render"
	"github.com/markbates/willie"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
)

func Test_DefaultContext_Param(t *testing.T) {
	r := require.New(t)
	c := DefaultContext{
		params: url.Values{
			"name": []string{"Mark"},
		},
	}

	r.Equal("Mark", c.Param("name"))
}

func Test_DefaultContext_ParamInt(t *testing.T) {
	r := require.New(t)
	c := DefaultContext{
		params: url.Values{
			"name": []string{"Mark"},
			"id":   []string{"1"},
		},
	}

	id, err := c.ParamInt("id")
	r.NoError(err)
	r.Equal(1, id)

	_, err = c.ParamInt("badkey")
	r.Error(err)

	_, err = c.ParamInt("name")
	r.Error(err)
}

func Test_DefaultContext_GetSet(t *testing.T) {
	r := require.New(t)
	c := DefaultContext{data: map[string]interface{}{}}
	r.Nil(c.Get("name"))

	c.Set("name", "Mark")
	r.NotNil(c.Get("name"))
	r.Equal("Mark", c.Get("name").(string))
}

func Test_DefaultContext_Render(t *testing.T) {
	r := require.New(t)

	res := httptest.NewRecorder()
	c := DefaultContext{
		response: res,
		params:   url.Values{"name": []string{"Mark"}},
		data:     map[string]interface{}{"greet": "Hello"},
		logger:   logrus.New(),
	}

	err := c.Render(123, render.String("{{.greet}} {{.params.name}}!"))
	r.NoError(err)

	r.Equal(123, res.Code)
	r.Equal("Hello Mark!", res.Body.String())
}

func Test_DefaultContext_Bind_Default(t *testing.T) {
	r := require.New(t)

	user := struct {
		FirstName string `schema:"first_name"`
	}{}

	a := New(Options{})
	a.POST("/", func(c Context) error {
		err := c.Bind(&user)
		if err != nil {
			return errors.WithStack(err)
		}
		return c.NoContent(201)
	})

	w := willie.New(a)
	uv := url.Values{"first_name": []string{"Mark"}}
	res := w.Request("/").Post(uv)
	r.Equal(201, res.Code)

	r.Equal("Mark", user.FirstName)
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
			return errors.WithStack(err)
		}
		return c.NoContent(201)
	})

	w := willie.New(a)
	res := w.JSON("/").Post(map[string]string{
		"first_name": "Mark",
	})
	r.Equal(201, res.Code)

	r.Equal("Mark", user.FirstName)
}

// func Test_DefaultContext_Bind_XML(t *testing.T) {
// 	r := require.New(t)
//
// 	user := struct {
// 		FirstName string `json:"first_name"`
// 	}{}
//
// 	a := New(Options{})
// 	a.POST("/", func(c Context) error {
// 		err := c.Bind(&user)
// 		if err != nil {
// 			return errors.WithStack(err)
// 		}
// 		return c.NoContent(201)
// 	})
//
// 	w := willie.New(a)
// 	req := w.Request("/")
// 	req.Headers["Content-Type"] = "application/xml"
// 	b, err := xml.Marshal(map[string]string{
// 		"first_name": "Mark",
// 	})
// 	r.NoError(err)
// 	res := req.Post(bytes.NewReader(b))
// 	r.Equal(201, res.Code)
//
// 	r.Equal("Mark", user.FirstName)
// }
