package buffalo

import (
	"github.com/gorilla/mux"
)

/* TODO: consider to split out Home (or Router, whatever) from App #road-to-v1
   Group and Domain based multi-homing are actually not an App if the concept
   of the App represents the application. The App should be only one for whole
   application.

   For an extreme example, App.Group().Stop() or even App.Group().Serve() are
   still valid function calls while they should not be allowed and the result
   could be strage.
*/

// Home is a container for Domains and Groups that independantly serves a
// group of pages with its own Middleware and ErrorHandlers. It is usually
// a multi-homed server domain or group of paths under a certain prefix.
//
// While the App is for managing whole application life cycle along with its
// default Home, including initializing and stopping its all components such
// as listeners and long-running jobs, Home is only for a specific group of
// services to serve its service logic efficiently.
type Home struct {
	app     *App // will replace App.root
	appSelf *App // temporary while the App is in action.
	// replace Options' Name, Host, and Prefix
	name   string
	host   string
	prefix string

	// moved from App
	// Middleware returns the current MiddlewareStack for the App/Group.
	Middleware    *MiddlewareStack `json:"-"`
	ErrorHandlers ErrorHandlers    `json:"-"`
	router        *mux.Router
	filepaths     []string
}
