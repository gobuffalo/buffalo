package buffalo

import (
	"time"

	"github.com/Sirupsen/logrus"
	humanize "github.com/dustin/go-humanize"
	"github.com/markbates/going/randx"
)

// RequestLogger can be be overridden to a user specified
// function that can be used to log the request.
var RequestLogger = RequestLoggerFunc

// RequestLoggerFunc is the default implementation of the RequestLogger.
// By default it will log a uniq "request_id", the HTTP Method of the request,
// the path that was requested, the duration (time) it took to process the
// request, the size of the response (and the "human" size), and the status
// code of the response.
func RequestLoggerFunc(h Handler) Handler {
	return func(c Context) error {
		var irid interface{}
		if irid = c.Session().Get("requestor_id"); irid == nil {
			irid = randx.String(10)
			c.Session().Set("requestor_id", irid)
			c.Session().Save()
		}
		now := time.Now()
		c.LogFields(logrus.Fields{
			"request_id": irid.(string) + "-" + randx.String(10),
			"method":     c.Request().Method,
			"path":       c.Request().URL.String(),
		})
		ct := c.Request().Header.Get("Content-Type")
		if ct != "" {
			c.LogField("content_type", ct)
		}
		defer func() {
			ws := c.Response().(*buffaloResponse)
			c.LogFields(logrus.Fields{
				"duration":   time.Now().Sub(now),
				"size":       ws.size,
				"human_size": humanize.Bytes(uint64(ws.size)),
				"status":     ws.status,
			})
			c.Logger().Info()
		}()
		return h(c)
	}
}
