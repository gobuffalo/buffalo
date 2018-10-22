package buffalo

import (
	"github.com/gobuffalo/logger"
	"github.com/markbates/oncer"
)

// Logger interface is used throughout Buffalo
// apps to log a whole manner of things.
type Logger = logger.FieldLogger

// NewLogger is deprecated. Use github.com/gobuffalo/logger.New instead.
func NewLogger(level string) logger.FieldLogger {
	oncer.Deprecate(0, "github.com/gobuffalo/buffalo#NewLogger", "Use github.com/gobuffalo/logger.New instead.")
	return logger.NewLogger(level)
}
