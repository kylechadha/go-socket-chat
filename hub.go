package main

type hub struct {
	clients    map[*client]bool
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
			break

		case c := <-h.unregister:
			_, ok := h.clients[c]
			if ok {
				delete(h.clients, c)
				close(c.send)
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
		case c.send <- []byte(message):
			break

		// Client channel has closed / is no longer available.
		default:
			close(c.send)
			delete(h.clients, c)
		}
	}
}
