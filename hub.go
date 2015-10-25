package main

import "strconv"

type hub struct {
	clients    map[*client]bool
	connected  int
	broadcast  chan string
	register   chan *client
	unregister chan *client
}

var h = hub{
	clients:    make(map[*client]bool),
	broadcast:  make(chan string),
	register:   make(chan *client),
	unregister: make(chan *client),
}

func (h *hub) activate() {
	for {
		select {
		case c := <-h.register:
			h.clients[c] = true

			// Increment the number of connected clients on the hub.
			h.connected++

			// Broadcast a user join message.
			h.broadcastMessage("User " + strconv.Itoa(c.id) + " has joined the chat.")
			break

		case c := <-h.unregister:
			_, ok := h.clients[c]
			if ok {
				// Remove the client and close its send channel.
				delete(h.clients, c)
				close(c.send)

				// Broadcast a user leave message.
				h.broadcastMessage("User " + strconv.Itoa(c.id) + " has left the chat.")
			}
			break

		case m := <-h.broadcast:
			h.broadcastMessage(m)
			break
		}
	}
}

func (h *hub) broadcastMessage(message string) {
	for c := range h.clients {
		select {
		// Send on the client channel, data will be picked up by the client's writePump.
		case c.send <- []byte(message):
			break

		// If we get here, client channel has closed / is no longer available.
		default:
			close(c.send)
			delete(h.clients, c)
		}
	}
}
