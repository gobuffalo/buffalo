package core

import "fmt"

// ErrNotInGoPath is thrown when not using go modules outside of GOPATH
var ErrNotInGoPath = fmt.Errorf("currently not in a $GOPATH")
