package release

import (
	"fmt"
	"os"
)

func Run() {
	createConfig()
	createRelease()
}

func createConfig() {
	fmt.Println("Configuring for dev")
	// Resource handles to the database, memcached and other backing services.
	// Credentials to external services such as Amazon S3 or Twitter.
	// Per-deploy values such as the canonical hostname for the deploy.
	os.Setenv("HARKENPORT", ":8000")
}

func createRelease() {
	fmt.Println("Creating release")
	// Create a new release folder
	// Copy ./main to the release folder
	// Write the config variables to the release folder (for reference only)
}
