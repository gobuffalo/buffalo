package buffalo

import (
	"crypto/rand"
	"time"

	humanize "github.com/dustin/go-humanize"
	"github.com/gobuffalo/buffalo/internal/httpx"
)

// RequestLogger can be be overridden to a user specified
// function that can be used to log the request.
var RequestLogger = RequestLoggerFunc

func randString(i int) (string, error) {
	if i == 0 {
		i = 64
	}
	b := make([]byte, i)
	_, err := rand.Read(b)
	return string(b), err
}

// RequestLoggerFunc is the default implementation of the RequestLogger.
// By default it will log a uniq "request_id", the HTTP Method of the request,
// the path that was requested, the duration (time) it took to process the
// request, the size of the response (and the "human" size), and the status
// code of the response.
func RequestLoggerFunc(h Handler) Handler {
	return func(c Context) error {
		rs, err := randString(10)
		if err != nil {
			return err
		}
		var irid interface{}
		if irid = c.Session().Get("requestor_id"); irid == nil {
			rs, err := randString(10)
			if err != nil {
				return err
			}
			irid = rs
			c.Session().Set("requestor_id", irid)
			c.Session().Save()
		}

		rid := irid.(string) + "-" + rs
		c.Set("request_id", rid)
		c.LogField("request_id", rid)

		start := time.Now()
		defer func() {
			ws, ok := c.Response().(*Response)
			if !ok {
				ws = &Response{ResponseWriter: c.Response()}
				ws.Status = 200
			}
			req := c.Request()
			ct := httpx.ContentType(req)
			if ct != "" {
				c.LogField("content_type", ct)
			}
			c.LogFields(map[string]interface{}{
				"method":     req.Method,
				"path":       req.URL.String(),
				"duration":   time.Since(start),
				"size":       ws.Size,
				"human_size": humanize.Bytes(uint64(ws.Size)),
				"status":     ws.Status,
			})
			c.Logger().Info(req.URL.String())
		}()
		return h(c)
	}
}
