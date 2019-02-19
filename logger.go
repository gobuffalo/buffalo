package buffalo

import (
	"github.com/gobuffalo/logger"
)

// Logger interface is used throughout Buffalo
// apps to log a whole manner of things.
type Logger = logger.FieldLogger
