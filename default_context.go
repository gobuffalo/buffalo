package buffalo

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"

	"github.com/gobuffalo/buffalo/binding"
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

// Set a value onto the Context. Any value set onto the Context
// will be automatically available in templates.
func (d *DefaultContext) Set(key string, value interface{}) {
	d.data[key] = value
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

// Cookies for the associated request and response.
func (d *DefaultContext) Cookies() *Cookies {
	return &Cookies{d.request, d.response}
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
	start := time.Now()
	defer func() {
		d.LogField("render", time.Since(start))
	}()
	if rr != nil {
		data := d.data
		pp := map[string]string{}
		for k, v := range d.params {
			pp[k] = v[0]
		}
		data["params"] = pp
		data["flash"] = d.Flash().data
		data["session"] = d.Session()
		data["request"] = d.Request()
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
// is "application/xml" it will use "xml.NewDecoder". See the
// github.com/gobuffalo/buffalo/binding package for more details.
func (d *DefaultContext) Bind(value interface{}) error {
	return binding.Exec(d.Request(), value)
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

	if len(args) > 0 {
		url = fmt.Sprintf(url, args...)
	}
	http.Redirect(d.Response(), d.Request(), url, status)
	return nil
}

// Data contains all the values set through Get/Set.
func (d *DefaultContext) Data() map[string]interface{} {
	return d.data
}

func (d *DefaultContext) String() string {
	bb := make([]string, 0, len(d.data))

	for k, v := range d.data {
		if _, ok := v.(RouteHelperFunc); !ok {
			bb = append(bb, fmt.Sprintf("%s: %s", k, v))
		}
	}
	sort.Strings(bb)
	return strings.Join(bb, "\n\n")
}

// File returns an uploaded file by name, or an error
func (d *DefaultContext) File(name string) (binding.File, error) {
	req := d.Request()
	if err := req.ParseMultipartForm(5 * 1024 * 1024); err != nil {
		return binding.File{}, err
	}
	f, h, err := req.FormFile(name)
	bf := binding.File{
		File:       f,
		FileHeader: h,
	}
	if err != nil {
		return bf, errors.WithStack(err)
	}
	return bf, nil
}

var defaultUpgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}
