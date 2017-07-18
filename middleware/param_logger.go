package middleware

import (
	"encoding/json"

	"github.com/gobuffalo/buffalo"
)

// ParameterLogger logs form and parameter values to the logger
func ParameterLogger(next buffalo.Handler) buffalo.Handler {
	return func(c buffalo.Context) error {
		defer func() {
			req := c.Request()
			if req.Method != "GET" {
				b, err := json.Marshal(req.Form)
				if err == nil {
					c.LogField("form", string(b))
				}
			}
			b, err := json.Marshal(c.Params())
			if err == nil {
				c.LogField("params", string(b))
			}
		}()
		return next(c)
	}
}
