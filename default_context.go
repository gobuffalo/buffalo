package buffalo

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/schema"
	"github.com/gorilla/websocket"
	"github.com/markbates/buffalo/render"
	"github.com/pkg/errors"
)

// DefaultContext is, as its name implies, a default
// implementation of the Context interface.
type DefaultContext struct {
	response    http.ResponseWriter
	request     *http.Request
	params      url.Values
	logger      Logger
	session     *Session
	contentType string
	data        map[string]interface{}
}

// Response returns the original Response for the request.
func (d *DefaultContext) Response() http.ResponseWriter {
	return d.response
}

// Request returns the original Request.
func (d *DefaultContext) Request() *http.Request {
	return d.request
}

// Params returns all of the parameters for the request,
// including both named params and query string parameters.
// These parameters are automatically available in templates
// as "{{.params}}".
func (d *DefaultContext) Params() ParamValues {
	return d.params
}

// Logger returns the Logger for this context.
func (d *DefaultContext) Logger() Logger {
	return d.logger
}

// Param returns a param, either named or query string,
// based on the key.
func (d *DefaultContext) Param(key string) string {
	return d.Params().Get(key)
}

// ParamInt tries to convert the requested parameter to
// an int. It will  return an error if there is a problem.
func (d *DefaultContext) ParamInt(key string) (int, error) {
	k := d.Params().Get(key)
	i, err := strconv.Atoi(k)
	return i, errors.WithMessage(err, fmt.Sprintf("could not convert %s to an int", k))
}

// Set a value onto the Context. Any value set onto the Context
// will be automatically available in templates.
func (d *DefaultContext) Set(key string, value interface{}) {
	d.data[key] = value
}

// Get a value that was previous set onto the Context.
func (d *DefaultContext) Get(key string) interface{} {
	return d.data[key]
}

// Session for the associated Request.
func (d *DefaultContext) Session() *Session {
	return d.session
}

// Render a status code and render.Renderer to the associated Response.
// The request parameters will be made available to the render.Renderer
// "{{.params}}". Any values set onto the Context will also automatically
// be made available to the render.Renderer. To render "no content" pass
// in a nil render.Renderer.
func (d *DefaultContext) Render(status int, rr render.Renderer) error {
	now := time.Now()
	defer func() {
		d.LogField("render", time.Now().Sub(now))
	}()
	if rr != nil {
		d.Response().Header().Set("Content-Type", rr.ContentType())
		d.Response().WriteHeader(status)
		data := d.data
		pp := map[string]string{}
		for k, v := range d.params {
			pp[k] = v[0]
		}
		data["params"] = pp
		return rr.Render(d.Response(), data)
	}
	d.Response().WriteHeader(status)
	return nil
}

// Bind the interface to the request.Body. The type of binding
// is dependent on the "Content-Type" for the request. If the type
// is "application/json" it will use "json.NewDecoder". If the type
// is "application/xml" it will use "xml.NewDecoder". The default
// binder is "http://www.gorillatoolkit.org/pkg/schema".
func (d *DefaultContext) Bind(value interface{}) error {
	switch strings.ToLower(d.Request().Header.Get("Content-Type")) {
	case "application/json", "text/json", "json":
		return json.NewDecoder(d.Request().Body).Decode(value)
	case "application/xml", "text/xml", "xml":
		return xml.NewDecoder(d.Request().Body).Decode(value)
	default:
		err := d.Request().ParseForm()
		if err != nil {
			return err
		}
		dec := schema.NewDecoder()
		dec.IgnoreUnknownKeys(true)
		dec.ZeroEmpty(true)
		return dec.Decode(value, d.Request().PostForm)
	}
}

// LogField adds the key/value pair onto the Logger to be printed out
// as part of the request logging. This allows you to easily add things
// like metrics (think DB times) to your request.
func (d *DefaultContext) LogField(key string, value interface{}) {
	d.logger = d.logger.WithField(key, value)
}

// LogFields adds the key/value pairs onto the Logger to be printed out
// as part of the request logging. This allows you to easily add things
// like metrics (think DB times) to your request.
func (d *DefaultContext) LogFields(values map[string]interface{}) {
	d.logger = d.logger.WithFields(values)
}

func (d *DefaultContext) Error(status int, err error) error {
	err = errors.WithStack(err)
	d.Logger().Error(err)
	msg := fmt.Sprintf("%+v", err)
	d.Response().WriteHeader(status)

	ct := d.Request().Header.Get("Content-Type")
	switch strings.ToLower(ct) {
	case "application/json", "text/json", "json":
		err = json.NewEncoder(d.Response()).Encode(map[string]interface{}{
			"error": msg,
			"code":  status,
		})
	case "application/xml", "text/xml", "xml":
	default:
		_, err = d.Response().Write([]byte(msg))
	}
	return err
}

// Websocket returns an upgraded github.com/gorilla/websocket.Conn
// that can then be used to work with websockets easily.
func (d *DefaultContext) Websocket() (*websocket.Conn, error) {
	return defaultUpgrader.Upgrade(d.Response(), d.Request(), nil)
}

// Redirect a request with the given status to the given URL.
func (d *DefaultContext) Redirect(status int, url string, args ...interface{}) error {
	http.Redirect(d.Response(), d.Request(), fmt.Sprintf(url, args...), status)
	return nil
}

// Data contains all the values set through Get/Set.
func (d *DefaultContext) Data() map[string]interface{} {
	return d.data
}

var defaultUpgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}
