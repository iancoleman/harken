package audience

import (
	"harken/base/http"
)

/*
Provides a facility to subscribe to a 'notification audience' where all
members will be notified when events related to that audience happen.
*/

type socketAudience struct {
	connections map[*http.Connection]bool
	broadcast   chan http.OutgoingMsg
	register    chan *http.Connection
	unregister  chan *http.Connection
}

var SocketAudiences = make(map[string]socketAudience)

func (a *socketAudience) run() {
	for {
		select {
		case c := <-a.register:
			a.connections[c] = true
		case c := <-a.unregister:
			delete(a.connections, c)
		case m := <-a.broadcast:
			for c := range a.connections {
				select {
				case c.Send <- m:
				default:
					delete(a.connections, c)
					close(c.Send)
					go c.Ws.Close()
				}
			}
		}
	}
}

func Subscribe(c *http.Connection, audienceId string) {
	if _, exists := SocketAudiences[audienceId]; !exists {
		createAudience(audienceId)
	}
	SocketAudiences[audienceId].register <- c
}

func Unsubscribe(c *http.Connection, audienceId string) {
	SocketAudiences[audienceId].unregister <- c
	if len(SocketAudiences[audienceId].connections) == 0 {
		delete(SocketAudiences, audienceId)
	}
}

func Broadcast(audienceId string, data http.OutgoingMsg) {
	SocketAudiences[audienceId].broadcast <- data
}

func createAudience(audienceId string) {
	a := socketAudience{
		broadcast:   make(chan http.OutgoingMsg),
		register:    make(chan *http.Connection),
		unregister:  make(chan *http.Connection),
		connections: make(map[*http.Connection]bool),
	}
	SocketAudiences[audienceId] = a
	go a.run()
}
