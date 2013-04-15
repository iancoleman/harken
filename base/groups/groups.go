package groups

import (
	"harken/base/http"
)

/*
Provides a facility to subscribe to a 'notification group' where all
members will be notified when events related to that group happen.
*/

type socketGroup struct {
	connections map[*http.Connection]bool
	broadcast   chan http.OutgoingMsg
	register    chan *http.Connection
	unregister  chan *http.Connection
}

var SocketGroups = make(map[string]socketGroup)

func (g *socketGroup) run() {
	for {
		select {
		case c := <-g.register:
			g.connections[c] = true
		case c := <-g.unregister:
			delete(g.connections, c)
		case m := <-g.broadcast:
			for c := range g.connections {
				select {
				case c.Send <- m:
				default:
					delete(g.connections, c)
					close(c.Send)
					go c.Ws.Close()
				}
			}
		}
	}
}

func Subscribe(c *http.Connection, groupId string) {
	if _, exists := SocketGroups[groupId]; !exists {
		createGroup(groupId)
	}
	SocketGroups[groupId].register <- c
}

func Unsubscribe(c *http.Connection, groupId string) {
	SocketGroups[groupId].unregister <- c
	if len(SocketGroups[groupId].connections) == 0 {
		delete(SocketGroups, groupId)
	}
}

func Broadcast(groupId string, data http.OutgoingMsg) {
	SocketGroups[groupId].broadcast <- data
}

func createGroup(groupId string) {
	g := socketGroup{
		broadcast:   make(chan http.OutgoingMsg),
		register:    make(chan *http.Connection),
		unregister:  make(chan *http.Connection),
		connections: make(map[*http.Connection]bool),
	}
	SocketGroups[groupId] = g
	go g.run()
}
