package run

import (
	"fmt"
	"io"
	"os"
	"os/exec"
)

func Run() {
	fmt.Println("Starting dev server")
	// Get the latest release folder
	// runServer := exec.Command("./main")
	runServer := exec.Command("go", "run", "main.go")
	stdout, err := runServer.StdoutPipe()
	stderr, err := runServer.StderrPipe()
	err = runServer.Start()
	go io.Copy(os.Stdout, stdout)
	go io.Copy(os.Stderr, stderr)
	runServer.Wait()
	if err != nil {
		fmt.Println(err)
	}
}
