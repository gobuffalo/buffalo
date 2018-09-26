// A Go web development eco-system, designed to make your life easier.
//
// Buffalo helps you to generate a web project that already has everything from front-end (JavaScript, SCSS, etc.) to back-end (database, routing, etc.) already hooked up and ready to run. From there it provides easy APIs to build your web application quickly in Go.
//
// Buffalo **isn't just a framework**, it's a holistic web development environment and project structure that **lets developers get straight to the business** of, well, building their business.
//
// > I :heart: web dev in go again - Brian Ketelsen
package buffalo

import "github.com/gobuffalo/buffalo/buffalo/cmd"

func init() {
	// fixes https://github.com/gobuffalo/buffalo/issues/1323
	var _ = cmd.RootCmd
}
