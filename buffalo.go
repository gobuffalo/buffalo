/*
Package buffalo is a Go web development eco-system, designed to make your life easier.

Buffalo helps you to generate a web project that already has everything from front-end (JavaScript, SCSS, etc.) to back-end (database, routing, etc.) already hooked up and ready to run. From there it provides easy APIs to build your web application quickly in Go.

Buffalo **isn't just a framework**, it's a holistic web development environment and project structure that **lets developers get straight to the business** of, well, building their business.
*/
package buffalo

// we need to import the runtime package
// as its needed by `buffalo build` and without
// this import the package doesn't get vendored
// by go mod vendor or by dep. this import fixes
// this problem.
import _ "github.com/gobuffalo/buffalo/runtime"
