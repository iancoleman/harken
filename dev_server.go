package main

import (
	"harken/build-release-run/build"
	"harken/build-release-run/release"
	"harken/build-release-run/run"
)

// Inspired by http://12factor.net/build-release-run
func main() {
	build.Run()
	release.Run()
	run.Run()
}
