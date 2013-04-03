package build

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
)

func Run() {
	// fetchDependencies()
	buildRoutes()
	copyStaticFiles()
	compile()
}

var imports []string
var routes []string
var inits []string
var	rootFolder string

func buildRoutes() {
	fmt.Println("Compiling routes")
	getRootFolder()
	walkExtsForRoutes()
	writeRoutesToFile()
}

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
			funcNameBits := funcRegexp.FindStringSubmatch(string(line))
			if len(funcNameBits) > 0 {
				funcName := funcNameBits[1]
				lastChar := string(funcName[len(funcName)-1])
				hasRoutes = true
				call := pkg + "." + funcName
				if funcName == "Init_" {
					inits = append(inits, call)
				} else if lastChar != "_" { // _ at end of func implies private
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

func copyStaticFiles() {
	filepath.Walk("exts", getStaticFilesFromPath)
}

func getStaticFilesFromPath(src string, f os.FileInfo, err error) error {
	rSep := "[\\\\/]" // Path separator. Hints that regexp maybe non-ideal
	rExtName := "([A-Za-z]+)"
	rFileName := "(.*\\..*)"
	extPath := "exts" + rSep + rExtName + rSep + "static" + rSep + rFileName
	r, _ := regexp.Compile(extPath)
	match := r.FindStringSubmatch(src)
	if len(match) > 0 {
		// Process the file - this is currently pretty basic, could be improved
		// a lot
		ext := match[1]
		fileLocation := match[2]
		filePath := path.Dir(fileLocation)
		fileName := path.Base(fileLocation)
		// Pipeline this file - preprocessing and minification
		// Make sure the directory exists
		dstPath := path.Join("static", ext, filePath)
		dst := path.Join(dstPath, fileName)
		os.MkdirAll(dstPath, 0755)
		// Copy the resulting file content to ./static/EXT/BLAH
		srcContent, err := ioutil.ReadFile(src)
		if err != nil {
			panic(err)
		}
		err = ioutil.WriteFile(dst, srcContent, 0744)
		if err != nil {
			panic(err)
		}
	}
	return nil
}

func compile() {
	// Only need to do this for a deployment, not for dev server
	/*
	fmt.Println("Building app")
	compilation := exec.Command("go", "build", "main.go")
	compilation.Start()
	*/
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
