package core

import "fmt"

// ErrGoModulesWithDep is thrown when trying to use both dep and go modules
var ErrGoModulesWithDep = fmt.Errorf("dep and modules can not be used at the same time")

// ErrNotInGoPath is thrown when not using go modules outside of GOPATH
var ErrNotInGoPath = fmt.Errorf("currently not in a $GOPATH")
