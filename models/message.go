package models

import (
	"encoding/json"
	"time"
)

// MessageType represents different types of messages
type MessageType string

const (
	MessageTypeChat     MessageType = "chat"
	MessageTypeJoin     MessageType = "join"
	MessageTypeLeave    MessageType = "leave"
	MessageTypeSystem   MessageType = "system"
	MessageTypeUserList MessageType = "userlist"
	MessageTypeGIF      MessageType = "gif" // New GIF type
)

// Add GIF-specific fields to Message struct
type Message struct {
	Type      MessageType `json:"type"`
	Username  string      `json:"username"`
	Content   string      `json:"content"`
	Room      string      `json:"room"`
	Timestamp time.Time   `json:"timestamp"`
	Color     string      `json:"color,omitempty"`
	GIFName   string      `json:"gif_name,omitempty"` // GIF identifier
	IsGIF     bool        `json:"is_gif,omitempty"`   // Flag for GIF messages
}

// User represents a connected user
type User struct {
	Username string    `json:"username"`
	Room     string    `json:"room"`
	JoinedAt time.Time `json:"joined_at"`
	Color    string    `json:"color"`
}

// Room represents a chat room
type Room struct {
	Name  string `json:"name"`
	Users []User `json:"users"`
}

// ToJSON converts message to JSON
func (m *Message) ToJSON() []byte {
	data, _ := json.Marshal(m)
	return data
}

// FromJSON creates message from JSON
func MessageFromJSON(data []byte) (*Message, error) {
	var msg Message
	err := json.Unmarshal(data, &msg)
	return &msg, err
}

// NewMessage creates a new message
func NewMessage(msgType MessageType, username, content, room string) *Message {
	return &Message{
		Type:      msgType,
		Username:  username,
		Content:   content,
		Room:      room,
		Timestamp: time.Now(),
	}
}

// FormatTime returns formatted timestamp
func (m *Message) FormatTime() string {
	return m.Timestamp.Format("15:04:05")
}
