package client

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"terminal-chat/gifs"
	"terminal-chat/models"
	"terminal-chat/utils"
	"time"

	"github.com/gorilla/websocket"
)

// Client represents a chat client
type Client struct {
	conn     *websocket.Conn
	username string
	room     string
	ui       *UI
	done     chan struct{}
}

// StartClient starts the chat client
func StartClient(host, port string) {
	// ... existing menu code ...

	// Get user input
	username, room, err := GetUserInput()
	if err != nil {
		log.Fatal("Error getting user input:", err)
	}

	// Get server address
	serverAddr, err := ShowConnectionMenu()
	if err != nil {
		log.Fatal("Error selecting server:", err)
	}

	// Create client
	client := &Client{
		username: username,
		room:     room,
		done:     make(chan struct{}),
	}

	// Connect to server
	if err := client.connect(serverAddr); err != nil {
		log.Fatal("Failed to connect to server:", err)
	}

	// Initialize UI with fixed input bar
	client.ui = NewUI(username, room)
	client.ui.InitScreen()

	// Handle graceful shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		client.disconnect()
		os.Exit(0)
	}()

	// Start message handling
	go client.readMessages()

	// Start input handling with fixed input bar
	client.handleInputWithBar()
}

// handleInputWithBar handles user input in the fixed input bar
func (c *Client) handleInputWithBar() {
	// Enable raw mode for character-by-character input
	scanner := bufio.NewScanner(os.Stdin)

	for {
		// Position cursor in input bar
		c.ui.positionCursorForChat()

		fmt.Print(utils.ColorWhite(""))
		if scanner.Scan() {
			input := scanner.Text()

			if input == "" {
				continue
			}

			// Clear the input bar
			c.ui.ClearInput()

			// Handle commands
			if strings.HasPrefix(input, "/") {
				if c.handleCommand(input) {
					break // Exit chat
				}
			} else {
				// Send regular message
				c.sendMessage(input)
			}
		}
	}
}

// sendMessage sends a message to the server
func (c *Client) sendMessage(content string) {
	if c.conn == nil {
		return
	}

	msg := models.NewMessage(models.MessageTypeChat, c.username, content, c.room)

	c.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))

	if err := c.conn.WriteJSON(msg); err != nil {
		log.Printf("Error sending message: %v", err)
		return
	}
}

// readMessages handles incoming messages
func (c *Client) readMessages() {
	defer c.conn.Close()

	for {
		select {
		case <-c.done:
			return
		default:
			_, message, err := c.conn.ReadMessage()
			if err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
					log.Printf("WebSocket error: %v", err)
				}
				c.ui.ShowDisconnected()
				return
			}

			c.processMessage(message)
		}
	}
}

// processMessage processes incoming messages
func (c *Client) processMessage(message []byte) {
	var msg models.Message
	if err := json.Unmarshal(message, &msg); err != nil {
		return
	}

	// Display message in chat area (not mixed with input)
	c.ui.DisplayMessage(msg)

	// Handle user list updates
	if msg.Type == models.MessageTypeUserList {
		c.ui.UpdateUserList([]string{})
	}
}

// handleCommand processes chat commands
func (c *Client) handleCommand(command string) bool {
	parts := strings.Fields(command)
	if len(parts) == 0 {
		return false
	}

	cmd := strings.ToLower(parts[0])

	switch cmd {
	case "/quit", "/exit":
		c.disconnect()
		return true

	case "/help":
		c.ui.ShowHelp()

	case "/users":
		c.ui.ShowUserList()

	case "/clear":
		c.ui.ClearChat()

	case "/gif":
		c.handleGIFCommand(parts)

	case "/gifs":
		c.showAvailableGIFs()

	default:
		systemMsg := models.Message{
			Type:      models.MessageTypeSystem,
			Username:  "system",
			Content:   fmt.Sprintf("Unknown command: %s. Type /help for available commands.", cmd),
			Timestamp: time.Now(),
		}
		c.ui.DisplayMessage(systemMsg)
	}

	return false
}
func (c *Client) handleGIFCommand(parts []string) {
	if len(parts) < 2 {
		systemMsg := models.Message{
			Type:      models.MessageTypeSystem,
			Username:  "system",
			Content:   "Usage: /gif <name>. Type /gifs to see available GIFs.",
			Timestamp: time.Now(),
		}
		c.ui.DisplayMessage(systemMsg)
		return
	}

	gifName := parts[1]
	if _, exists := gifs.GetGIF(gifName); !exists {
		systemMsg := models.Message{
			Type:      models.MessageTypeSystem,
			Username:  "system",
			Content:   fmt.Sprintf("GIF '%s' not found. Type /gifs to see available GIFs.", gifName),
			Timestamp: time.Now(),
		}
		c.ui.DisplayMessage(systemMsg)
		return
	}

	// Send GIF message
	msg := models.Message{
		Type:      models.MessageTypeGIF,
		Username:  c.username,
		Content:   fmt.Sprintf("sent a GIF: %s", gifName),
		Room:      c.room,
		Timestamp: time.Now(),
		GIFName:   gifName,
		IsGIF:     true,
	}

	if err := c.conn.WriteJSON(msg); err != nil {
		log.Printf("Error sending GIF: %v", err)
		return
	}
}

// showAvailableGIFs shows list of available GIFs
func (c *Client) showAvailableGIFs() {
	availableGIFs := gifs.GetAvailableGIFs()
	gifList := strings.Join(availableGIFs, ", ")

	systemMsg := models.Message{
		Type:      models.MessageTypeSystem,
		Username:  "system",
		Content:   fmt.Sprintf("Available GIFs: %s", gifList),
		Timestamp: time.Now(),
	}
	c.ui.DisplayMessage(systemMsg)
}

// connect and disconnect methods remain the same...
func (c *Client) connect(serverAddr string) error {
	u := url.URL{
		Scheme: "ws",
		Host:   serverAddr,
		Path:   "/ws",
	}

	q := u.Query()
	q.Set("username", c.username)
	q.Set("room", c.room)
	u.RawQuery = q.Encode()

	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		return err
	}

	c.conn = conn
	return nil
}

func (c *Client) disconnect() {
	close(c.done)

	if c.conn != nil {
		c.conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
		c.conn.Close()
	}

	if c.ui != nil {
		c.ui.ShowGoodbye()
	}
}
