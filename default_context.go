package buffalo

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/gobuffalo/buffalo/render"
	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
)

// assert that DefaultContext is implementing Context
var _ Context = &DefaultContext{}
var _ context.Context = &DefaultContext{}

// DefaultContext is, as its name implies, a default
// implementation of the Context interface.
type DefaultContext struct {
	context.Context
	response    http.ResponseWriter
	request     *http.Request
	params      url.Values
	logger      Logger
	session     *Session
	contentType string
	data        map[string]interface{}
	flash       *Flash
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
// an int. It will return an error if there is a problem.
func (d *DefaultContext) ParamInt(key string) (int, error) {
	warningMsg := "Context#ParamInt is deprecated, and will be removed in v0.9.0."

	_, file, no, ok := runtime.Caller(1)
	if ok {
		warningMsg = fmt.Sprintf("%s Called from %s:%d", warningMsg, file, no)
	}

	d.Logger().Warn(warningMsg)

	k := d.Params().Get(key)
	i, err := strconv.Atoi(k)
	return i, errors.WithMessage(err, fmt.Sprintf("could not convert %s to an int", k))
}

// Set a value onto the Context. Any value set onto the Context
// will be automatically available in templates.
func (d *DefaultContext) Set(key string, value interface{}) {
	d.data[key] = value
}

// Get is deprecated. Please use Value instead.
func (d *DefaultContext) Get(key string) interface{} {
	warningMsg := "Context#Get is deprecated, and will be removed in v0.9.0. Please use Context#Value instead."

	_, file, no, ok := runtime.Caller(1)
	if ok {
		warningMsg = fmt.Sprintf("%s Called from %s:%d", warningMsg, file, no)
	}

	d.Logger().Warn(warningMsg)
	return d.Value(key)
}

// Value that has previously stored on the context.
func (d *DefaultContext) Value(key interface{}) interface{} {
	if k, ok := key.(string); ok {
		if v, ok := d.data[k]; ok {
			return v
		}
	}
	return d.Context.Value(key)
}

// Session for the associated Request.
func (d *DefaultContext) Session() *Session {
	return d.session
}

// Flash messages for the associated Request.
func (d *DefaultContext) Flash() *Flash {
	return d.flash
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
		data := d.data
		pp := map[string]string{}
		for k, v := range d.params {
			pp[k] = v[0]
		}
		data["params"] = pp
		data["flash"] = d.Flash().data
		bb := &bytes.Buffer{}

		err := rr.Render(bb, data)
		if err != nil {
			return HTTPError{Status: 500, Cause: errors.WithStack(err)}
		}

		if d.Session() != nil {
			d.Flash().Clear()
			d.Flash().persist(d.Session())
		}

		d.Response().Header().Set("Content-Type", rr.ContentType())
		d.Response().WriteHeader(status)
		_, err = io.Copy(d.Response(), bb)
		if err != nil {
			return HTTPError{Status: 500, Cause: errors.WithStack(err)}
		}

		return nil
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
	ct := strings.ToLower(d.Request().Header.Get("Content-Type"))
	if ct != "" {
		cts := strings.Split(ct, ";")
		c := cts[0]
		if b, ok := binders[strings.TrimSpace(c)]; ok {
			return b(d.Request(), value)
		}
		return errors.Errorf("could not find a binder for %s", c)
	}
	return errors.New("blank content type")
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
	return HTTPError{Status: status, Cause: errors.WithStack(err)}
}

// Websocket returns an upgraded github.com/gorilla/websocket.Conn
// that can then be used to work with websockets easily.
func (d *DefaultContext) Websocket() (*websocket.Conn, error) {
	return defaultUpgrader.Upgrade(d.Response(), d.Request(), nil)
}

// Redirect a request with the given status to the given URL.
func (d *DefaultContext) Redirect(status int, url string, args ...interface{}) error {
	d.Flash().persist(d.Session())

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
