package updater

import "github.com/gobuffalo/buffalo/buffalo/cmd/fix"

// Check interface for runnable checker functions
type Check = fix.Check

// Runner will run all compatible checks
type Runner = fix.Runner

// Run all compatible checks
var Run = fix.Run
