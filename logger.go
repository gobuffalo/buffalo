package buffalo

import (
	"time"

	"github.com/Sirupsen/logrus"
	humanize "github.com/dustin/go-humanize"
	"github.com/markbates/going/randx"
)

type Logger interface {
	WithField(string, interface{}) Logger
	WithFields(map[string]interface{}) Logger
	WithError(error) Logger
	Debugf(string, ...interface{})
	Infof(string, ...interface{})
	Printf(string, ...interface{})
	Warnf(string, ...interface{})
	Warningf(string, ...interface{})
	Errorf(string, ...interface{})
	Fatalf(string, ...interface{})
	Debug(...interface{})
	Info(...interface{})
	Warn(...interface{})
	Error(...interface{})
	Fatal(...interface{})
	Panic(...interface{})
}

var RequestLogger = func(h Handler) Handler {
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

func NewLogger() Logger {
	return &multiLogger{Loggers: []logrus.FieldLogger{}}
}
