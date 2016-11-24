package buffalo

import (
	"time"

	"github.com/Sirupsen/logrus"
	humanize "github.com/dustin/go-humanize"
)

type Logger interface {
	logrus.FieldLogger
}

var RequestLogger = func(h Handler) Handler {
	return func(c Context) error {
		now := time.Now()
		c.LogFields(logrus.Fields{
			"method": c.Request().Method,
			"path":   c.Request().URL,
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
