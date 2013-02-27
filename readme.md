Start the server
---

`$ go run main.go`

or alternatively,

`$ go build main.go && ./main`

Adding features
---

Add your features to `exts` folder.

Add the require routes and initialisation scripts for your features to
`base/setup/exts.go`

Make sure there's at least an index.html file in base/static with enough basic
code in it to connect to the websocket. This will be done with javascript.
An example index.html will be provided in the near future.

Using it
---

Start the server.

Navigate to `http://localhost:8000/` in your browser

Testing it
---

Test the base packages

`$ go test base/*`

Test the extensions you've installed

`$ go test exts/*`

Extensions
---

You'll have to write your own. Examples will be provided in the future.

This software is only extremely young so this stuff isn't written yet.
