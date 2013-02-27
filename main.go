package main

import (
	"harken/base/http"
	"harken/base/setup"
)

func main() {
	setup.Init()
	http.Start()
}
