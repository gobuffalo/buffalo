package buffalo

import "github.com/Sirupsen/logrus"

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

func NewLogger(level string) Logger {
	l := logrus.New()
	l.Level, _ = logrus.ParseLevel(level)
	return &multiLogger{Loggers: []logrus.FieldLogger{l}}
}
