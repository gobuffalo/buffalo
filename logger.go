package buffalo

import "github.com/Sirupsen/logrus"

type Logger interface {
	logrus.FieldLogger
}
