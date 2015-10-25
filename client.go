package main

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/websocket"
)

const (
	maxMessageSize = 1024 * 1024
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	writeWait      = 10 * time.Second
)

type client struct {
	conn *websocket.Conn
	send chan []byte
	id   int
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  maxMessageSize,
	WriteBufferSize: maxMessageSize,
}

func socketHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Socket connection opened.")

	if r.Method != "GET" {
		http.Error(w, "Method not allowed.", 405)
		return
	}

	// Upgrade the HTTP connection to WS.
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	// Create a new client.
	c := &client{
		conn: conn,
		send: make(chan []byte, maxMessageSize),
		id:   h.connected + 1,
	}

	// Register the client with the hub.
	h.register <- c

	// Kick off the read / write pumps.
	go c.writePump()
	c.readPump()
}

func (c *client) readPump() {
	// If the for loop exits (there's an error), unregister the client.
	defer func() {
		h.unregister <- c
		c.conn.Close()
	}()

	// Set the Read Limit and Deadline.
	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))

	// Set the Pong handler to reset the Read Deadline.
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		// Read the message from the socket -- this is a blocking operation.
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			break
		}

		// Broadcast the message to all clients.
		h.broadcast <- "User " + strconv.Itoa(c.id) + ": " + string(message)
	}
}

func (c *client) writePump() {
	// Create a ticker based on the pingPeriod.
	ticker := time.NewTicker(pingPeriod)

	// If the for select exits (there's an error), stop the ticker.
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			if !ok {
				// If a zero value was sent on the channel (channel closed),
				// close the socket connection.
				log.Println("Client disconnected.")
				c.write(websocket.CloseMessage, []byte{})
				return
			}

			log.Println("Sending data to the client.")
			if err := c.write(websocket.TextMessage, message); err != nil {
				return
			}
		case <-ticker.C:
			log.Println("Pinging the client.")
			if err := c.write(websocket.PingMessage, []byte{}); err != nil {
				return
			}
		}
	}
}

func (c *client) write(messageType int, message []byte) error {
	c.conn.SetWriteDeadline(time.Now().Add(writeWait))
	return c.conn.WriteMessage(messageType, message)
}
