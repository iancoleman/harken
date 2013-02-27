package http

import (
	"code.google.com/p/go.net/websocket"
	"net/http"
)

/*
Handles an incoming message and routes it to the appropriate place.
*/

var Routes = make(map[string]func(string, string) *OutgoingMsg)

func call(session string, key string, data string, ch chan *OutgoingMsg) {

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

func initUrls() {
	http.Handle("/ws", websocket.Handler(wsHandler))
	http.Handle("/", http.FileServer(http.Dir("static")))
}
