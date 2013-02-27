package http

import (
	"code.google.com/p/go.net/websocket"
	"net/http"
)

/*
Handles an incoming message and routes it to the appropriate place.
Also manages permissions for incoming messages. Also see setup/setup.go
*/

var Routes = make(map[string]func(string, string) *OutgoingMsg)

// See
// http://www.mikespook.com/2012/07/function-call-by-name-in-golang/
// http://stackoverflow.com/questions/6769020/go-map-of-functions
func Call(session string, key string, data string, ch chan *OutgoingMsg) {

    // Call the function
    _, containsKey := Routes[key]
    if containsKey {
		ch <- Routes[key](session, data)
	} else {
		msg := "Invalid ext.Method: " + key
		// TODO 'did you mean' - look for similar methods
		response := OutgoingMsg{Ext:"http", Method:"Call", Error:msg}
		ch <- &response
	}
}

// See http://stackoverflow.com/questions/6564558/wildcards-in-the-pattern-for-http-handlefunc

func initUrls() {
	http.Handle("/ws", websocket.Handler(wsHandler))
	http.Handle("/", http.FileServer(http.Dir("static")))
}
