package setup

// The hope is to have this automatically populated in the future, so that
// routes can be specified in each extension and there's no need for devs to
// touch any code in harken/base.

import (
	// "harken/base/http"

	// "harken/exts/db"
)

func initExts() {
	// Put any extension initialisation scripts here
	// Don't forget to include the package in the imports

	// db.Init()
}

func bufferRoutes() {
	// Put any extension routes here
	// Don't forget to include the package in the imports

	// http.Routes["users.New"] = users.New
}
