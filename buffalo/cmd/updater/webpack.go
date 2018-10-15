package updater

import "github.com/gobuffalo/buffalo/buffalo/cmd/fix"

// WebpackCheck will compare the current default Buffalo
// webpack.config.js against the applications webpack.config.js. If they are
// different you have the option to overwrite the existing webpack.config.js
// file with the new one.
var WebpackCheck = fix.WebpackCheck
