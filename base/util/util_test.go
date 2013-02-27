package util

import (
	"testing"
)

func TestToken(t *testing.T) {
	// CreateToken creates a token of random alphanumeric characters of the
	// specified length.
	token := CreateToken(14)
	if len(token) != 14 {
		t.Error("Token length expected to be 14 but is %s", len(token))
	}
	// also need to test for randomness
}
