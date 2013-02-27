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
	broadcast chan http.OutgoingMsg
	register chan *http.Connection
	unregister chan *http.Connection
}

var socketGroups = make(map[string]socketGroup)

func (g *socketGroup) run() {
	for {
		select {
		case c := <-g.register:
			g.connections[c] = true
		case c := <-g.unregister:
			delete(g.connections, c)
			close(c.Send)
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
	/*
	if the group doesn't exist, make it.
	then register the user to the group.
	g := socketGroups[id]
	g.register <- c
	*/

}

func Unsubscribe(c *http.Connection, groupId string) {
	/*
	unregister the user from the group.
	if the group has zero members, delete it.
	g := socketGroups[id]
	g.unregister <- c
	*/
}
