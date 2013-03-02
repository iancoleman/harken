package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"regexp"
)

// Inspired by http://12factor.net/build-release-run
func main() {
	build()
	release()
	run()
}

func build() {
	// fetchDependencies()
	buildRoutes()
	compile()
}

func buildRoutes() {
	fmt.Println("Compiling routes")
	walkExtsForRoutes()
	writeRoutesToFile()
}

var routes []string

func walkExtsForRoutes() {
	err := filepath.Walk("exts", getRoutesFromFile)
	if err != nil {
		panic(err)
	}
}

func getRoutesFromFile(filename string, f os.FileInfo, err error) error {
	p := len(filename)
	isGoFile := filename[p-3:p] == ".go"
	isTestFile := p > 8 && filename[p-8:p] == "_test.go"
	r, err := regexp.Compile(`func\s+([A-Z][A-Za-z0-9_-]+)`)
	if isGoFile && !isTestFile {
		f, err := os.Open(filename)
		if err != nil {
			panic(err)
		}
		defer f.Close()
		bf := bufio.NewReader(f)
		for {
			line, isPrefix, err := bf.ReadLine()
			if err == io.EOF {
				break
			}
			if err != nil {
				panic(err)
			}
			if isPrefix {
				panic("Error: Unexpected ling line")
			}
			instance := r.FindStringSubmatch(string(line))
			if len(instance) > 0 && instance[1] != "Init" {
				_, pkg := path.Split(path.Dir(filename))
				call := pkg + "." + instance[1]
				routes = append(routes, call)
			}
		}
	}
	return err
}

func writeRoutesToFile() {
}

func compile() {
	fmt.Println("Building app")
	compilation := exec.Command("go", "build", "main.go")
	compilation.Start()
}

func release() {
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

func run() {
	fmt.Println("Starting dev server")
	// Get the latest release folder
	runServer := exec.Command("./main")
	stdout, err := runServer.StdoutPipe()
	err = runServer.Start()
	go io.Copy(os.Stdout, stdout)
	runServer.Wait()
	fmt.Println(err)
}
