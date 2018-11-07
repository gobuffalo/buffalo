package core

import "github.com/pkg/errors"

var ErrGoModulesWithDep = errors.New("dep and modules can not be used at the same time")
var ErrNotInGoPath = errors.New("currently not in a $GOPATH")
