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
	"github.com/markbates/buffalo/render"
	"github.com/pkg/errors"
)

type DefaultContext struct {
	response    http.ResponseWriter
	request     *http.Request
	params      url.Values
	logger      Logger
	session     *Session
	contentType string
	data        map[string]interface{}
}

func (d *DefaultContext) Response() http.ResponseWriter {
	return d.response
}

func (d *DefaultContext) Request() *http.Request {
	return d.request
}

func (d *DefaultContext) Params() ParamValues {
	return d.params
}

func (d *DefaultContext) Logger() Logger {
	return d.logger
}

func (d *DefaultContext) Param(key string) string {
	return d.Params().Get(key)
}

func (d *DefaultContext) ParamInt(key string) (int, error) {
	k := d.Params().Get(key)
	i, err := strconv.Atoi(k)
	return i, errors.Wrapf(err, "could not convert %s to an int", k)
}

func (d *DefaultContext) Set(key string, value interface{}) {
	d.data[key] = value
}

func (d *DefaultContext) Get(key string) interface{} {
	return d.data[key]
}

func (d *DefaultContext) Session() *Session {
	return d.session
}

func (d *DefaultContext) Render(status int, rr render.Renderer) error {
	now := time.Now()
	defer func() {
		d.LogField("render", time.Now().Sub(now))
	}()
	d.response.Header().Set("Content-Type", rr.ContentType())
	d.response.WriteHeader(status)
	data := d.data
	pp := map[string]string{}
	for k, v := range d.params {
		pp[k] = v[0]
	}
	data["params"] = pp
	err := rr.Render(d.response, data)
	return err
}

func (d *DefaultContext) Bind(value interface{}) error {
	switch strings.ToLower(d.request.Header.Get("Content-Type")) {
	case "application/json":
		return json.NewDecoder(d.request.Body).Decode(value)
	case "application/xml":
		return xml.NewDecoder(d.request.Body).Decode(value)
	default:
		err := d.request.ParseForm()
		if err != nil {
			return err
		}
		return schema.NewDecoder().Decode(value, d.request.PostForm)
	}
}

func (d *DefaultContext) NoContent(status int) error {
	d.response.WriteHeader(status)
	return nil
}

func (d *DefaultContext) LogField(key string, value interface{}) {
	d.logger = d.logger.WithField(key, value)
}

func (d *DefaultContext) LogFields(values map[string]interface{}) {
	d.logger = d.logger.WithFields(values)
}

func (d *DefaultContext) Error(status int, err error) error {
	err = errors.WithStack(err)
	d.Logger().Errorln(err)
	msg := fmt.Sprintf("%+v", err)
	d.response.WriteHeader(status)
	_, err = d.response.Write([]byte(msg))
	return err
}
