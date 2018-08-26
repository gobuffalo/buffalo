package updater

import "github.com/gobuffalo/buffalo/buffalo/cmd/fix"

// DeprecrationsCheck will either log, or fix, deprecated items in the application
var DeprecrationsCheck = fix.DeprecrationsCheck
