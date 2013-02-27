package config

import (
	"os"
	"testing"
)

func TestInitWithoutExistingPortEnvvar(t *testing.T) {
	// Calling init sets the port variable to default
	os.Setenv("HARKENPORT", "")
	Init()
	if Port != ":8000" {
		t.Error("Expected port to be :8000 but it was '%s'", Port)
	}
}

func TestInitWithExistingPortEnvvar(t *testing.T) {
	// Calling init sets the port variable
	os.Setenv("HARKENPORT", "anyvalue")
	Init()
	if Port != "anyvalue" {
		t.Error("Expected port to be 'anyvalue' but it was '%s'", Port)
	}
}
