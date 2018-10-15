package updater

import "github.com/gobuffalo/buffalo/buffalo/cmd/fix"

// DepEnsure runs `dep ensure -v` to make sure that any newly changed
// imports are added to dep.
var DepEnsure = fix.DepEnsure
