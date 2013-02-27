package config

import (
	"fmt"
	"os"
)

var Port string

func Init() {
	Port = loadEnvVar("HARKENPORT", ":8000")
	initExtsConfig()
}

func loadEnvVar(variable string, defaultValue string) string {
	v := os.Getenv(variable)
	if len(v) == 0 {
		fmt.Println("set " + variable + " environment variable.")
		// set the message to also say how to do that in
		// the current OS
		fmt.Println("using default of " + defaultValue)
		v = defaultValue
	}
	return v
}
