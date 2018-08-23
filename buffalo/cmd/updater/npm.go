package updater

import "github.com/gobuffalo/buffalo/buffalo/cmd/fix"

// PackageJSONCheck will compare the current default Buffalo
// package.json against the applications package.json. If they are
// different you have the option to overwrite the existing package.json
// file with the new one.
var PackageJSONCheck = fix.PackageJSONCheck
