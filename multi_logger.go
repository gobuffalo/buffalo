package buffalo

import "github.com/Sirupsen/logrus"

var (
	_ Logger = &multiLogger{}
)

type multiLogger struct {
	Loggers []logrus.FieldLogger
}

func (m *multiLogger) WithField(key string, value interface{}) Logger {
	lgs := []logrus.FieldLogger{}
	for _, l := range m.Loggers {
		lgs = append(lgs, l.WithField(key, value))
	}
	return &multiLogger{Loggers: lgs}
}

func (m *multiLogger) WithFields(fields map[string]interface{}) Logger {
	lgs := []logrus.FieldLogger{}
	for _, l := range m.Loggers {
		lgs = append(lgs, l.WithFields(fields))
	}
	return &multiLogger{Loggers: lgs}
}

func (m *multiLogger) Debugf(format string, args ...interface{}) {
	for _, l := range m.Loggers {
		l.Debugf(format, args...)
	}
}

func (m *multiLogger) Infof(format string, args ...interface{}) {
	for _, l := range m.Loggers {
		l.Infof(format, args...)
	}
}

func (m *multiLogger) Printf(format string, args ...interface{}) {
	for _, l := range m.Loggers {
		l.Printf(format, args...)
	}
}

func (m *multiLogger) Warnf(format string, args ...interface{}) {
	for _, l := range m.Loggers {
		l.Warnf(format, args...)
	}
}

func (m *multiLogger) Errorf(format string, args ...interface{}) {
	for _, l := range m.Loggers {
		l.Errorf(format, args...)
	}
}

func (m *multiLogger) Fatalf(format string, args ...interface{}) {
	for _, l := range m.Loggers {
		l.Fatalf(format, args...)
	}
}

func (m *multiLogger) Debug(args ...interface{}) {
	for _, l := range m.Loggers {
		l.Debug(args...)
	}
}

func (m *multiLogger) Info(args ...interface{}) {
	for _, l := range m.Loggers {
		l.Info(args...)
	}
}

func (m *multiLogger) Warn(args ...interface{}) {
	for _, l := range m.Loggers {
		l.Warn(args...)
	}
}

func (m *multiLogger) Error(args ...interface{}) {
	for _, l := range m.Loggers {
		l.Error(args...)
	}
}

func (m *multiLogger) Fatal(args ...interface{}) {
	for _, l := range m.Loggers {
		l.Fatal(args...)
	}
}

func (m *multiLogger) Panic(args ...interface{}) {
	for _, l := range m.Loggers {
		l.Panic(args...)
	}
}
