package buffalo

import (
	"time"

	"github.com/Sirupsen/logrus"
	humanize "github.com/dustin/go-humanize"
	"github.com/markbates/going/randx"
)

var RequestLogger = RequestLoggerFunc

func RequestLoggerFunc(h Handler) Handler {
	return func(c Context) error {
		now := time.Now()
		c.LogFields(logrus.Fields{
			"request_id": randx.String(10),
			"method":     c.Request().Method,
			"path":       c.Request().URL.String(),
		})
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
