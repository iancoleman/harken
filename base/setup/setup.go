package setup

import (
	"harken/base/config"
	"runtime"
)

func Init() {
	config.Init()
	setMaxProcs()
	bufferRoutes()
	initExts()
}

func setMaxProcs() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}
