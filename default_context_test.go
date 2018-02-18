package buffalo

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"golang.org/x/net/websocket"

	"github.com/gobuffalo/buffalo/render"
	"github.com/markbates/willie"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
)

func basicContext() DefaultContext {
	return DefaultContext{
		Context: context.Background(),
		logger:  NewLogger("debug"),
		data:    make(map[string]interface{}),
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

	w := willie.New(a)
	res := w.Request("/").Get()
	r.Equal(u, res.Location())
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
			return errors.WithStack(err)
		}
		return c.Render(201, nil)
	})

	w := willie.New(a)
	uv := url.Values{"first_name": []string{"Mark"}}
	res := w.Request("/").Post(uv)
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
			return errors.WithStack(err)
		}
		return c.Render(201, nil)
	})

	w := willie.New(a)
	uv := url.Values{"first_name": []string{""}}
	res := w.Request("/").Post(uv)
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
			return errors.WithStack(err)
		}
		return c.Render(201, nil)
	})

	w := willie.New(a)
	res := w.JSON("/").Post(map[string]string{
		"first_name": "Mark",
	})
	r.Equal(201, res.Code)

	r.Equal("Mark", user.FirstName)
}

func Test_DefaultContext_Websocket(t *testing.T) {
	r := require.New(t)

	type Message struct {
		Original  string    `json:"original"`
		Formatted string    `json:"formatted"`
		Received  time.Time `json:"received"`
	}

	a := New(Options{})
	a.GET("/socket", func(c Context) error {
		conn, err := c.Websocket()
		if err != nil {
			return err
		}
		for {

			_, m, err := conn.ReadMessage()
			if err != nil {
				return err
			}

			data := string(m)

			msg := Message{
				Original:  data,
				Formatted: strings.ToUpper(data),
				Received:  time.Now(),
			}

			if err := conn.WriteJSON(msg); err != nil {
				return err
			}
		}
	})

	ts := httptest.NewServer(a)
	defer ts.Close()

	wsURL := strings.Replace(ts.URL, "http", "ws", 1) + "/socket"

	ws, err := websocket.Dial(wsURL, "", ts.URL)
	r.NoError(err)

	_, err = ws.Write([]byte("hello, world!"))
	r.NoError(err)

	msg := make([]byte, 512)
	read, err := ws.Read(msg)
	r.NoError(err)

	var message Message
	err = json.NewDecoder(bytes.NewReader(msg[:read])).Decode(&message)
	r.NoError(err)

	// Create a table of what we expect.
	tests := []struct {
		Got  string
		Want string
	}{
		{message.Formatted, "HELLO, WORLD!"},
		{message.Original, "hello, world!"},
	}

	// Check the different fields.
	for _, tt := range tests {
		r.Equal(tt.Want, tt.Got)
	}
}
