// +build appengine

package buffalo

import (
	"io"
)

func checkIfTerminal(w io.Writer) bool {
	return true
}
