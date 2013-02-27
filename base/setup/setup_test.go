package setup

import (
	"testing"
)

func TestInitConfig(t *testing.T) {
	// setup.Init calls config.Init
}

func TestInitSetMaxProcs(t *testing.T) {
	// setup.Init sets the processors to the maximum available
}

func TestInitRoutes(t *testing.T) {
	// setup.Init calls bufferRoutes which loads the routes for the websocket
	// into memory in the variable http.Routes
}

func TestInitExts(t *testing.T) {
	// setup.Init calls initExts which calls any initialisation code for the
	// exts
}
