package buffalo

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"golang.org/x/net/websocket"

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
		logger:   &multiLogger{},
	}

	err := c.Render(123, render.String("{{greet}} {{params.name}}!"))
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
		return c.Render(201, nil)
	})

	w := willie.New(a)
	uv := url.Values{"first_name": []string{"Mark"}}
	res := w.Request("/").Post(uv)
	r.Equal(201, res.Code)

	r.Equal("Mark", user.FirstName)
}

func Test_DefaultContext_Bind_Default_BlankFields(t *testing.T) {
	r := require.New(t)

	user := struct {
		FirstName string `schema:"first_name"`
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

func Test_DefaultContext_Error_Default(t *testing.T) {
	r := require.New(t)

	res := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/", nil)
	r.NoError(err)
	c := DefaultContext{
		response: res,
		request:  req,
		logger:   &multiLogger{},
	}

	c.Error(123, errors.New("Boom!"))
	r.Equal(123, res.Code)
	r.Contains(res.Body.String(), "Boom!")
}

func Test_DefaultContext_Error_JSON(t *testing.T) {
	r := require.New(t)

	res := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/", nil)
	r.NoError(err)
	req.Header.Set("Content-Type", "application/json")

	c := DefaultContext{
		response: res,
		request:  req,
		logger:   &multiLogger{},
	}

	c.Error(123, errors.New("Boom!"))
	r.Equal(123, res.Code)
	r.Contains(res.Body.String(), "Boom!")
	r.Contains(res.Body.String(), `"error":"Boom!`)
	r.Contains(res.Body.String(), `"code":123`)
}

func Test_DefaultContext_Websocket(t *testing.T) {
	r := require.New(t)

	type Message struct {
		Original  string    `json:"original"`
		Formatted string    `json:"formatted"`
		Received  time.Time `json:"received"`
	}

	a := Automatic(Options{})
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

	ws, err := websocket.Dial(wsURL, "", "http://127.0.0.1")
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
