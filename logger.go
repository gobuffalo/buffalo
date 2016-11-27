package buffalo

import (
	"time"

	"github.com/Sirupsen/logrus"
	humanize "github.com/dustin/go-humanize"
	"github.com/markbates/going/randx"
)

var (
	_ Logger = &MultiLogger{}
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

type MultiLogger struct {
	Loggers []logrus.FieldLogger
}

func NewLogger() Logger {
	return &MultiLogger{Loggers: []logrus.FieldLogger{}}
}

func (m *MultiLogger) WithField(key string, value interface{}) Logger {
	lgs := []logrus.FieldLogger{}
	for _, l := range m.Loggers {
		lgs = append(lgs, l.WithField(key, value))
	}
	return &MultiLogger{Loggers: lgs}
}

func (m *MultiLogger) WithFields(fields map[string]interface{}) Logger {
	lgs := []logrus.FieldLogger{}
	for _, l := range m.Loggers {
		lgs = append(lgs, l.WithFields(fields))
	}
	return &MultiLogger{Loggers: lgs}
}

func (m *MultiLogger) WithError(err error) Logger {
	lgs := []logrus.FieldLogger{}
	for _, l := range m.Loggers {
		lgs = append(lgs, l.WithError(err))
	}
	return &MultiLogger{Loggers: lgs}
}

func (m *MultiLogger) Debugf(format string, args ...interface{}) {
	for _, l := range m.Loggers {
		l.Debugf(format, args...)
	}
}

func (m *MultiLogger) Infof(format string, args ...interface{}) {
	for _, l := range m.Loggers {
		l.Infof(format, args...)
	}
}

func (m *MultiLogger) Printf(format string, args ...interface{}) {
	for _, l := range m.Loggers {
		l.Printf(format, args...)
	}
}

func (m *MultiLogger) Warnf(format string, args ...interface{}) {
	for _, l := range m.Loggers {
		l.Warnf(format, args...)
	}
}

func (m *MultiLogger) Warningf(format string, args ...interface{}) {
	for _, l := range m.Loggers {
		l.Warningf(format, args...)
	}
}

func (m *MultiLogger) Errorf(format string, args ...interface{}) {
	for _, l := range m.Loggers {
		l.Errorf(format, args...)
	}
}

func (m *MultiLogger) Fatalf(format string, args ...interface{}) {
	for _, l := range m.Loggers {
		l.Fatalf(format, args...)
	}
}

func (m *MultiLogger) Debug(args ...interface{}) {
	for _, l := range m.Loggers {
		l.Debug(args...)
	}
}

func (m *MultiLogger) Info(args ...interface{}) {
	for _, l := range m.Loggers {
		l.Info(args...)
	}
}

func (m *MultiLogger) Warn(args ...interface{}) {
	for _, l := range m.Loggers {
		l.Warn(args...)
	}
}

func (m *MultiLogger) Error(args ...interface{}) {
	for _, l := range m.Loggers {
		l.Error(args...)
	}
}

func (m *MultiLogger) Fatal(args ...interface{}) {
	for _, l := range m.Loggers {
		l.Fatal(args...)
	}
}

func (m *MultiLogger) Panic(args ...interface{}) {
	for _, l := range m.Loggers {
		l.Panic(args...)
	}
}
