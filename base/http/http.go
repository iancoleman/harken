package http

/*
Starts the http server and passes requests to Routes
*/

import (
	"code.google.com/p/go.net/websocket"
	"fmt"
	"harken/base/config"
	"harken/base/util"
	"net/http"
)

var Connections = make(map[string]*Connection)

type IncomingMsg struct {
	Ext    string
	Method string
	Data   string
}

type OutgoingMsg struct {
	Ext    string
	Method string
	Data   string
	Error  string
}

type Connection struct {
	Owner   string
	Send    chan OutgoingMsg
	Session string
	Ws      *websocket.Conn
}

func (c *Connection) reader() {
	for {
		var message IncomingMsg
		err := websocket.JSON.Receive(c.Ws, &message)
		fmt.Println("Received via ws", message)
		if err != nil {
			// TODO something
			break
		}
		callId := message.Ext + "." + message.Method
		ch := make(chan *OutgoingMsg)
		go call(c.Session, callId, message.Data, ch)
		response := <-ch
		if response != nil {
			c.Send <- *response
		}
	}
	c.Ws.Close()
}

func (c *Connection) writer() {
	for message := range c.Send {
		err := websocket.JSON.Send(c.Ws, message)
		if err != nil {
			// TODO something
			break
		}
	}
	delete(Connections, c.Session)
	c.Ws.Close()
}

func wsHandler(ws *websocket.Conn) {
	session := util.CreateToken(32)
	c := &Connection{
		Send:    make(chan OutgoingMsg, 256),
		Session: session,
		Ws:      ws}
	Connections[session] = c
	go c.writer()
	c.reader()
}

func Start() {
	fmt.Println("Serving at http://localhost" + config.Port)
	initUrls() // see routes.go
	if err := http.ListenAndServe(config.Port, nil); err != nil {
		fmt.Println("ListenAndServe:", err)
	}
}
