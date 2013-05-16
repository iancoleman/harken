Start the server
---

$ go build main.go
$ ./main

if it doesn't work, check that you have the GOPATH variable set correctly
to /path/to/project/ which should then contain src/harken
eg /path/to/project/src/harken should be a valid folder.

Set this variable by executing the command

$ export GOPATH=/path/to/project

Adding features
---

Add your features to `exts` folder.

Make sure there's at least an index.html file in base/static with enough basic
code in it to connect to the websocket. This will be done with javascript.
An example index.html is provided.

Developing with it
---

Start the development server with the command `go run dev_server.go`

Navigate to `http://localhost:8000/` in your browser

Use it
---

Test the base packages

`$ go test ./base/*`

Test the extensions you've installed

`$ go test ./exts/*`

Extensions
---

You'll have to write your own. A basic sample extension called test is provided.

This software is still very young.
