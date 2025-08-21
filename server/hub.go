package server

import (
	"log"
	"strings"
	"terminal-chat/models"
	"terminal-chat/utils"
	"time"
)

// Hub maintains the set of active clients and broadcasts messages
type Hub struct {
	clients    map[*Client]bool
	rooms      map[string]map[*Client]bool
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
	userColors map[string]string
}

// NewHub creates a new Hub
func NewHub() *Hub {
	return &Hub{
		clients:    make(map[*Client]bool),
		rooms:      make(map[string]map[*Client]bool),
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		userColors: make(map[string]string),
	}
}

// Run starts the hub
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.registerClient(client)

		case client := <-h.unregister:
			h.unregisterClient(client)

		case message := <-h.broadcast:
			h.broadcastMessage(message)
		}
	}
}

func (h *Hub) registerClient(client *Client) {
	// Clean the room name to avoid encoding issues
	client.Room = strings.TrimSpace(client.Room)

	h.clients[client] = true

	// Add to room
	if h.rooms[client.Room] == nil {
		h.rooms[client.Room] = make(map[*Client]bool)
	}
	h.rooms[client.Room][client] = true

	// Assign color to user
	if h.userColors[client.Username] == "" {
		colorIndex := len(h.userColors) % 6
		colors := []string{"red", "green", "yellow", "blue", "magenta", "cyan"}
		h.userColors[client.Username] = colors[colorIndex]
	}

	log.Printf("✓ User %s joined room '%s'", client.Username, client.Room)

	// Send join message
	joinMsg := models.NewMessage(models.MessageTypeJoin, client.Username,
		"joined the chat", client.Room)
	joinMsg.Color = h.userColors[client.Username]

	// Send messages separately with a small delay
	h.broadcastToRoom(joinMsg.ToJSON(), client.Room)

	// Small delay before sending user list
	time.Sleep(10 * time.Millisecond)
	h.sendUserList(client.Room)
}

func (h *Hub) unregisterClient(client *Client) {
	if _, ok := h.clients[client]; ok {
		delete(h.clients, client)
		close(client.send)

		// Remove from room
		if room, exists := h.rooms[client.Room]; exists {
			delete(room, client)
			if len(room) == 0 {
				delete(h.rooms, client.Room)
			}
		}

		log.Printf("%s User %s left room %s",
			utils.ColorRed("✗"), client.Username, client.Room)

		// Send leave message
		leaveMsg := models.NewMessage(models.MessageTypeLeave, client.Username,
			"left the chat", client.Room)
		leaveMsg.Color = h.userColors[client.Username]
		h.broadcastToRoom(leaveMsg.ToJSON(), client.Room)

		// Send updated user list
		h.sendUserList(client.Room)
	}
}

func (h *Hub) broadcastMessage(message []byte) {
	msg, err := models.MessageFromJSON(message)
	if err != nil {
		log.Printf("Error parsing message: %v", err)
		return
	}

	// Add color to message if not set
	if msg.Color == "" {
		msg.Color = h.userColors[msg.Username]
	}
	msg.Timestamp = time.Now()

	// Debug logging
	log.Printf("Broadcasting message - Type: %s, User: %s, Content: %s, Room: %s",
		msg.Type, msg.Username, msg.Content, msg.Room)

	// Broadcast to all clients in the room
	h.broadcastToRoom(msg.ToJSON(), msg.Room)
}

func (h *Hub) broadcastToRoom(message []byte, room string) {
	if roomClients, exists := h.rooms[room]; exists {
		log.Printf("📡 Broadcasting to %d clients in room '%s'", len(roomClients), room)

		successCount := 0
		for client := range roomClients {
			select {
			case client.send <- message:
				successCount++
				log.Printf("✅ Message sent to client: %s", client.Username)
			default:
				// Client's send channel is full or closed
				log.Printf("❌ Failed to send to client: %s (removing)", client.Username)
				close(client.send)
				delete(h.clients, client)
				delete(roomClients, client)
			}
		}
		log.Printf("📊 Successfully sent to %d/%d clients", successCount, len(roomClients))
	} else {
		log.Printf("⚠️ Room '%s' not found for broadcasting", room)
	}
}

func (h *Hub) sendUserList(room string) {
	var users []models.User
	if roomClients, exists := h.rooms[room]; exists {
		for client := range roomClients {
			user := models.User{
				Username: client.Username,
				Room:     client.Room,
				JoinedAt: time.Now(),
				Color:    h.userColors[client.Username],
			}
			users = append(users, user)
		}
	}

	userListMsg := models.NewMessage(models.MessageTypeUserList, "system", "", room)
	// You can add users to message content as JSON if needed

	h.broadcastToRoom(userListMsg.ToJSON(), room)
}
