package assets

import "fmt"

// LogoURL is the URL to the SVG logo to be used for assets
const LogoURL = "https://raw.githubusercontent.com/gobuffalo/buffalo/master/logo.svg"

func init() {
	fmt.Println("github.com/gobuffalo/buffalo/generators/assets has been deprecated in v0.13.0, and will be removed in v0.14.0. Use github.com/gobuffalo/buffalo/genny/assets directly.")
}
