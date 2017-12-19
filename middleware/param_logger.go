package middleware

import (
	"encoding/json"
	"mime/multipart"
	"net/url"

	"github.com/gobuffalo/buffalo"
	"github.com/pkg/errors"
)

// ParameterLogger logs form and parameter values to the logger
func ParameterLogger(next buffalo.Handler) buffalo.Handler {
	return func(c buffalo.Context) error {
		defer func() {
			req := c.Request()
			if req.Method != "GET" {
				if err := postParamLogger(c); err != nil {
					c.Logger().Error(err)
				}
			}
			b, err := json.Marshal(c.Params())
			if err != nil {
				c.Logger().Error(err)
			}
			c.LogField("params", string(b))
		}()
		return next(c)
	}
}

func postParamLogger(c buffalo.Context) error {
	req := c.Request()
	mp := req.MultipartForm
	if mp != nil {
		return multipartParamLogger(mp, c)
	}

	b, err := json.Marshal(req.Form)
	if err != nil {
		return errors.WithStack(err)
	}
	c.LogField("form", string(b))
	return nil
}

func multipartParamLogger(mp *multipart.Form, c buffalo.Context) error {
	uv := url.Values{}
	for k, v := range mp.Value {
		for _, vv := range v {
			uv.Add(k, vv)
		}
	}
	for k, v := range mp.File {
		for _, vv := range v {
			uv.Add(k, vv.Filename)
		}
	}
	b, err := json.Marshal(uv)
	if err != nil {
		return errors.WithStack(err)
	}
	c.LogField("form", string(b))
	return nil
}
