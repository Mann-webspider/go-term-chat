package server

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"terminal-chat/models"
	"time"

	"github.com/gorilla/websocket"
)

const (
	writeWait      = 30 * time.Second
	pongWait       = 0
	pingPeriod     = 25 * time.Second
	maxMessageSize = 512
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// Client represents a WebSocket client
type Client struct {
	hub      *Hub
	conn     *websocket.Conn
	send     chan []byte
	Username string
	Room     string
}

// readPump pumps messages from the websocket connection to the hub
func (c *Client) readPump() {
	defer func() {
		log.Printf("Client %s disconnecting from room %s", c.Username, c.Room)
		c.hub.unregister <- c
		c.conn.Close()
	}()

	c.conn.SetReadLimit(maxMessageSize)
	// Remove all read deadlines - let connection stay alive indefinitely
	// c.conn.SetReadDeadline(time.Now().Add(pongWait)) // REMOVE THIS

	// Optional: Keep pong handler for heartbeat but without deadline
	c.conn.SetPongHandler(func(string) error {
		log.Printf("📡 Received pong from client %s", c.Username)
		return nil
	})

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("Unexpected close error for %s: %v", c.Username, err)
			}
			break
		}

		message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))

		var msg models.Message
		if err := json.Unmarshal(message, &msg); err == nil {
			log.Printf("💬 [%s] %s: %s", msg.Room, msg.Username, msg.Content)
		}

		c.hub.broadcast <- message
	}
}

// writePump pumps messages from the hub to the websocket connection
func (c *Client) writePump() {
	// Optional: Remove ticker for ping if you don't want heartbeat
	// ticker := time.NewTicker(pingPeriod)
	defer func() {
		// ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			// Send each message as separate frame
			if err := c.conn.WriteMessage(websocket.TextMessage, message); err != nil {
				log.Printf("❌ Write error for %s: %v", c.Username, err)
				return
			}

			// Remove automatic ping - only disconnect when client chooses to
			// case <-ticker.C:
			//     c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			//     if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
			//         return
			//     }
		}
	}
}

// ServeWS handles websocket requests from clients
func ServeWS(hub *Hub, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	username := r.URL.Query().Get("username")
	room := r.URL.Query().Get("room")

	if username == "" {
		username = "Anonymous"
	}
	if room == "" {
		room = "general"
	}

	client := &Client{
		hub:      hub,
		conn:     conn,
		send:     make(chan []byte, 256),
		Username: username,
		Room:     room,
	}

	client.hub.register <- client

	go client.writePump()
	go client.readPump()
}
