package webpack

import "github.com/gobuffalo/buffalo-cli/genny/assets/webpack"

// BinPath is the path to the local install of webpack
var BinPath = webpack.BinPath

// Templates used for generating webpack
// (exported mostly for the "fix" command)
var Templates = webpack.Templates

var New = webpack.New

type Options = webpack.Options
