package middleware

import (
	"encoding/json"
	"fmt"
	"mime/multipart"
	"net/url"
	"strings"

	"github.com/gobuffalo/mw-paramlogger"

	"github.com/gobuffalo/buffalo"
	"github.com/pkg/errors"
)

// ParameterExclusionList is the list of parameter names that will be filtered
// from the application logs (see maskSecrets).
// Important: this list will be used in case insensitive.
//
// Deprecated: use github.com/gobuffalo/mw-paramlogger#ParameterExclusionList instead.
var ParameterExclusionList = paramlogger.ParameterExclusionList

var filteredIndicator = []string{"[FILTERED]"}

// ParameterLogger logs form and parameter values to the logger
type parameterLogger struct {
	excluded []string
}

// ParameterLogger logs form and parameter values to the loggers
//
// Deprecated: use github.com/gobuffalo/mw-paramlogger#ParameterLogger instead.
var ParameterLogger = paramlogger.ParameterLogger

// Middleware is a buffalo middleware function to connect this parameter filterer with buffalo
//
// Deprecated: use github.com/gobuffalo/mw-paramlogger#Middleware instead.
func (pl parameterLogger) Middleware(next buffalo.Handler) buffalo.Handler {
	return func(c buffalo.Context) error {
		fmt.Printf("paramlogger Middleware is deprecated and will be removed in the next version. Please use github.com/gobuffalo/mw-paramlogger#Middleware instead.")
		defer func() {
			req := c.Request()
			if req.Method != "GET" {
				if err := pl.logForm(c); err != nil {
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

func (pl parameterLogger) logForm(c buffalo.Context) error {
	req := c.Request()
	mp := req.MultipartForm
	if mp != nil {
		return pl.multipartParamLogger(mp, c)
	}

	if err := pl.addFormFieldTo(c, req.Form); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (pl parameterLogger) multipartParamLogger(mp *multipart.Form, c buffalo.Context) error {
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

	if err := pl.addFormFieldTo(c, uv); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func (pl parameterLogger) addFormFieldTo(c buffalo.Context, form url.Values) error {
	maskedForm := pl.maskSecrets(form)
	b, err := json.Marshal(maskedForm)

	if err != nil {
		return err
	}

	c.LogField("form", string(b))
	return nil
}

// maskSecrets matches ParameterExclusionList against parameters passed in the
// request, and returns a copy of the request parameters replacing excluded params
// with [FILTERED].
func (pl parameterLogger) maskSecrets(form url.Values) url.Values {
	if len(pl.excluded) == 0 {
		pl.excluded = ParameterExclusionList
	}

	copy := url.Values{}
	for key, values := range form {
	blcheck:
		for _, excluded := range pl.excluded {
			copy[key] = values
			if strings.ToUpper(key) == strings.ToUpper(excluded) {
				copy[key] = filteredIndicator
				break blcheck
			}

		}
	}
	return copy
}
