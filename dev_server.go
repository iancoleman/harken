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
	"strings"
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
	getRootFolder()
	walkExtsForRoutes()
	writeRoutesToFile()
}

var imports []string
var routes []string
var inits []string
var	rootFolder string

func getRootFolder() {
	rootPath, _ := os.Getwd()
	_, rootFolder = path.Split(rootPath)
}

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
	funcRegexp, err := regexp.Compile(`func\s+([A-Z][A-Za-z0-9_-]+)`)
	pkgRegexp, err := regexp.Compile(`package\s+([A-Za-z0-9]+)`)
	if isGoFile && !isTestFile {
		hasRoutes := false
		f, err := os.Open(filename)
		if err != nil {
			panic(err)
		}
		defer f.Close()
		bf := bufio.NewReader(f)
		pkg := ""
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
			if len(pkg) == 0 {
				pkgMatch := pkgRegexp.FindStringSubmatch(string(line))
				if len(pkgMatch) > 0 {
					pkg = pkgMatch[1]
				}
			}
			funcName := funcRegexp.FindStringSubmatch(string(line))
			if len(funcName) > 0 {
				hasRoutes = true
				call := pkg + "." + funcName[1]
				if funcName[1] == "Init" {
					inits = append(inits, call)
				} else {
					routes = append(routes, call)
				}
			}
		}
		if hasRoutes {
			fullImport := path.Join(rootFolder, path.Dir(filename))
			imports = append(imports, fullImport)
		}
	}
	return err
}

func writeRoutesToFile() {
	timports := ""
	for _, i := range imports {
		timports = timports + "\t\"" + i + "\"\n"
	}
	template = strings.Replace(template, "{{imports}}", timports, 1)

	tinits := ""
	for _, i := range inits {
		tinits = tinits + "\t" + i + "()\n"
	}
	template = strings.Replace(template, "{{inits}}", tinits, 1)

	troutes := ""
	for _, i := range routes {
		troutes = troutes + "\thttp.Routes[\"" + i + "\"] = " + i + "\n"
	}
	template = strings.Replace(template, "{{routes}}", troutes, 1)

	f, err := os.Create("base/setup/exts.go")
	defer f.Close()
	if err != nil {
		panic(err)
	}

	_, err = f.WriteString(template)
	if err != nil {
		panic(err)
	}
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

var template = `
package setup

// Do not change this file.
// This file has been automatically generated.
// Any changes will not be retained.

import (
	"harken/base/http"
{{imports}})

func initExts() {
{{inits}}}

func bufferRoutes() {
{{routes}}}
`
