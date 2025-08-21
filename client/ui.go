package client

import (
	"fmt"
	"strings"
	"terminal-chat/gifs"
	"terminal-chat/models"
	"terminal-chat/utils"
	"time"

	"github.com/pterm/pterm"
)

// Add GIF animation tracking to UI struct
type UI struct {
	username       string
	room           string
	messages       []models.Message
	users          []string
	colorMap       map[string]func(...interface{}) string
	messageArea    *pterm.AreaPrinter
	chatHeight     int
	terminalWidth  int
	terminalHeight int
	activeGIFs     map[string]*GIFAnimation // Track active GIF animations
}

// GIFAnimation tracks an active GIF animation
type GIFAnimation struct {
	GIF          gifs.AnimatedGIF
	CurrentFrame int
	LastUpdate   time.Time
	Position     int // Line position in chat
	Username     string
	Timestamp    string
}

// Update NewUI to initialize activeGIFs
func NewUI(username, room string) *UI {
	width, height, err := pterm.GetTerminalSize()
	if err != nil {
		width = 80
		height = 24
	}

	return &UI{
		username:       username,
		room:           room,
		messages:       make([]models.Message, 0),
		users:          make([]string, 0),
		colorMap:       make(map[string]func(...interface{}) string),
		chatHeight:     height - 8,
		terminalWidth:  width,
		terminalHeight: height,
		activeGIFs:     make(map[string]*GIFAnimation), // Initialize GIF tracking
	}
}

// InitScreen initializes the chat screen with fixed input bar
func (ui *UI) InitScreen() {
	utils.ClearScreen()
	ui.showChatHeader()
	ui.showInputBar()
	ui.positionCursorForChat()
}

// showChatHeader displays the chat header
func (ui *UI) showChatHeader() {
	// Chat header
	headerPanel := pterm.DefaultPanel.
		WithPanels([][]pterm.Panel{
			{
				{Data: pterm.DefaultCenter.Sprint(fmt.Sprintf("💬 CHAT ROOM: %s", strings.ToUpper(ui.room)))},
			},
			{
				{Data: fmt.Sprintf("User: %s", utils.ColorGreen(ui.username))},
				{Data: fmt.Sprintf("Status: %s", utils.ColorSuccess("CONNECTED"))},
			},
		}).
		WithPadding(1)

	headerPanel.Render()

	// Chat messages header with colorful border
	chatBorder := strings.Repeat("─", ui.terminalWidth-2)
	fmt.Printf("%s┌%s┐%s\n",
		utils.ColorCyan(""),
		utils.ColorCyan(chatBorder),
		utils.ColorCyan(""))
	fmt.Printf("%s│%s CHAT MESSAGES %s│%s\n",
		utils.ColorCyan(""),
		utils.ColorYellow(""),
		strings.Repeat(" ", ui.terminalWidth-17),
		utils.ColorCyan(""))
}

// showInputBar displays the fixed input bar at the bottom
func (ui *UI) showInputBar() {
	// Move cursor to bottom of terminal
	fmt.Printf("\033[%d;1H", ui.terminalHeight-2)

	// Input bar with colorful border
	inputBorder := strings.Repeat("─", ui.terminalWidth-2)

	// Top border of input bar
	fmt.Printf("%s┌%s┐%s\n",
		utils.ColorMagenta(""),
		utils.ColorMagenta(inputBorder),
		utils.ColorMagenta(""))

	// Input line
	fmt.Printf("%s│ %s> %s%s│%s",
		utils.ColorMagenta(""),
		utils.ColorGreen(""),
		utils.ColorWhite(""),
		strings.Repeat(" ", ui.terminalWidth-6),
		utils.ColorMagenta(""))
}

// Update DisplayMessage to handle GIF messages properly
func (ui *UI) DisplayMessage(msg models.Message) {
	ui.messages = append(ui.messages, msg)

	if ui.colorMap[msg.Username] == nil {
		ui.colorMap[msg.Username] = utils.GetRandomColor(len(ui.colorMap))
	}

	userColor := ui.colorMap[msg.Username]
	timestamp := utils.ColorWhite(msg.FormatTime())

	var output string

	switch msg.Type {
	case models.MessageTypeGIF:
		// Handle GIF message - USE the userColor variable here
		if gif, exists := gifs.GetGIF(msg.GIFName); exists {
			ui.startGIFAnimation(gif, msg.Username, timestamp, len(ui.messages))
			username := userColor(fmt.Sprintf("%-12s", msg.Username)) // Use userColor
			output = fmt.Sprintf("%s│ %s %s │ %s (GIF: %s)%s",
				utils.ColorCyan(""), timestamp, username,
				utils.ColorMagenta("🎬"), msg.GIFName,
				utils.ColorCyan(""))
		}

	case models.MessageTypeChat:
		username := userColor(fmt.Sprintf("%-12s", msg.Username))
		content := utils.ColorWhite(msg.Content)
		output = fmt.Sprintf("%s│ %s %s │ %s%s│%s",
			utils.ColorCyan(""), timestamp, username, content,
			strings.Repeat(" ", ui.terminalWidth-len(msg.Content)-len(msg.Username)-15),
			utils.ColorCyan(""))

	case models.MessageTypeJoin:
		joinMsg := fmt.Sprintf("%s joined the chat", userColor(msg.Username))
		output = fmt.Sprintf("%s│ %s %s %s%s│%s",
			utils.ColorCyan(""), timestamp,
			utils.ColorGreen("→"), joinMsg,
			strings.Repeat(" ", ui.terminalWidth-len(joinMsg)-12),
			utils.ColorCyan(""))

	case models.MessageTypeLeave:
		leaveMsg := fmt.Sprintf("%s left the chat", userColor(msg.Username))
		output = fmt.Sprintf("%s│ %s %s %s%s│%s",
			utils.ColorCyan(""), timestamp,
			utils.ColorRed("←"), leaveMsg,
			strings.Repeat(" ", ui.terminalWidth-len(leaveMsg)-12),
			utils.ColorCyan(""))

	case models.MessageTypeSystem:
		systemMsg := utils.ColorYellow(msg.Content)
		output = fmt.Sprintf("%s│ %s %s %s%s│%s",
			utils.ColorCyan(""), timestamp,
			utils.ColorYellow("ℹ"), systemMsg,
			strings.Repeat(" ", ui.terminalWidth-len(msg.Content)-12),
			utils.ColorCyan(""))
	}

	if output != "" {
		fmt.Print("\033[s") // Save cursor position
		ui.printMessageAtPosition(output)
		fmt.Print("\033[u") // Restore cursor position
	}
}

// startGIFAnimation starts animating a GIF
func (ui *UI) startGIFAnimation(gif gifs.AnimatedGIF, username, timestamp string, position int) {
	animationID := fmt.Sprintf("%s_%d", username, position)

	animation := &GIFAnimation{
		GIF:          gif,
		CurrentFrame: 0,
		LastUpdate:   time.Now(),
		Position:     position,
		Username:     username,
		Timestamp:    timestamp,
	}

	ui.activeGIFs[animationID] = animation

	// Start animation goroutine
	go ui.animateGIF(animationID)
}

// animateGIF runs the GIF animation
func (ui *UI) animateGIF(animationID string) {
	animation, exists := ui.activeGIFs[animationID]
	if !exists {
		return
	}

	cycles := 0
	maxCycles := 3 // Limit animation cycles

	for {
		if !animation.GIF.Loop && cycles >= 1 {
			break
		}
		if animation.GIF.Loop && cycles >= maxCycles {
			break
		}

		for frameIndex, frame := range animation.GIF.Frames {
			time.Sleep(frame.Duration)

			// Update the frame in the chat
			ui.updateGIFFrame(animationID, frameIndex, frame.Content)

			// Check if animation should continue
			if _, exists := ui.activeGIFs[animationID]; !exists {
				return
			}
		}
		cycles++
	}

	// Clean up animation
	delete(ui.activeGIFs, animationID)
}

// updateGIFFrame updates a specific GIF frame in the chat
func (ui *UI) updateGIFFrame(animationID string, frameIndex int, frameContent string) {
	animation, exists := ui.activeGIFs[animationID]
	if !exists {
		return
	}

	// Save cursor position
	fmt.Print("\033[s")

	// Calculate the line position
	line := 6 + animation.Position

	// Ensure user color exists in colorMap
	if ui.colorMap[animation.Username] == nil {
		ui.colorMap[animation.Username] = utils.GetRandomColor(len(ui.colorMap))
	}

	// Format and display the frame (remove unused userColor variable)
	formattedFrame := gifs.FormatGIFFrame(frameContent, animation.Username, animation.Timestamp)

	// Move to the specific line and update
	fmt.Printf("\033[%d;1H\033[K", line) // Move to line and clear it
	fmt.Print(formattedFrame)

	// Restore cursor position
	fmt.Print("\033[u")
}

// printMessageAtPosition prints a message in the chat area
func (ui *UI) printMessageAtPosition(message string) {
	// Calculate next message line (scrolling if needed)
	messageCount := len(ui.messages)
	maxVisibleMessages := ui.chatHeight - 3

	if messageCount <= maxVisibleMessages {
		// Print at next available line
		line := 6 + messageCount
		fmt.Printf("\033[%d;1H%s\n", line, message)
	} else {
		// Scroll chat area
		ui.scrollChatArea()
		fmt.Printf("\033[%d;1H%s\n", ui.chatHeight+2, message)
	}
}

// scrollChatArea scrolls the chat messages up
func (ui *UI) scrollChatArea() {
	startLine := 7
	endLine := ui.chatHeight + 2

	// Move each line up by one
	for line := startLine; line < endLine; line++ {
		fmt.Printf("\033[%d;1H\033[K", line) // Clear line
		if line < endLine-1 {
			// This would need more complex logic to actually move text up
			// For now, we'll just clear and let new messages appear
		}
	}
}

// positionCursorForChat positions cursor in chat input area
func (ui *UI) positionCursorForChat() {
	// Position cursor in input bar
	fmt.Printf("\033[%d;5H", ui.terminalHeight-1)
}

// ClearInput clears the input bar
func (ui *UI) ClearInput() {
	// Clear input line
	fmt.Printf("\033[%d;5H", ui.terminalHeight-1)
	fmt.Printf("%s%s",
		utils.ColorWhite(""),
		strings.Repeat(" ", ui.terminalWidth-6))
	// Reposition cursor
	fmt.Printf("\033[%d;5H", ui.terminalHeight-1)
}

// ShowUserList displays online users in a side panel
func (ui *UI) ShowUserList() {
	// Save cursor, show users, restore cursor
	fmt.Print("\033[s")

	fmt.Printf("\033[%d;%dH", 7, ui.terminalWidth-25)
	userBox := pterm.DefaultBox.
		WithTitle("👥 Online").
		WithTitleTopCenter().
		WithBoxStyle(pterm.NewStyle(pterm.FgGreen))

	userList := ""
	for i, user := range ui.users {
		if ui.colorMap[user] == nil {
			ui.colorMap[user] = utils.GetRandomColor(len(ui.colorMap))
		}
		userColor := ui.colorMap[user]

		status := "●"
		if user == ui.username {
			status = utils.ColorGreen("● (you)")
		} else {
			status = utils.ColorGreen("●")
		}

		userList += fmt.Sprintf("%s %s", status, userColor(user))
		if i < len(ui.users)-1 {
			userList += "\n"
		}
	}

	if userList == "" {
		userList = "No users online"
	}

	// Print user box (this is simplified - in real implementation you'd position it properly)
	fmt.Print(userBox.Sprint(userList))

	fmt.Print("\033[u") // Restore cursor
}

// Update ShowHelp in client/ui.go
func (ui *UI) ShowHelp() {
	helpMsg := models.Message{
		Type:      models.MessageTypeSystem,
		Username:  "system",
		Content:   "Commands: /quit, /exit, /help, /users, /clear, /gif <name>, /gifs",
		Timestamp: time.Now(),
	}
	ui.DisplayMessage(helpMsg)
}

// ClearChat clears the chat area
func (ui *UI) ClearChat() {
	// Clear chat area but keep header and input bar
	for line := 7; line <= ui.chatHeight+2; line++ {
		fmt.Printf("\033[%d;1H\033[K", line)
	}

	// Reset messages
	ui.messages = make([]models.Message, 0)

	// Redraw chat border
	fmt.Printf("\033[7;1H%s│%s%s│%s\n",
		utils.ColorCyan(""),
		strings.Repeat(" ", ui.terminalWidth-2),
		utils.ColorCyan(""))
}

// ShowDisconnected shows disconnection message
func (ui *UI) ShowDisconnected() {
	disconnectMsg := models.Message{
		Type:      models.MessageTypeSystem,
		Username:  "system",
		Content:   "⚠️ Connection lost! Type /quit to exit.",
		Timestamp: time.Now(),
	}
	ui.DisplayMessage(disconnectMsg)
}

// ShowGoodbye shows farewell message
func (ui *UI) ShowGoodbye() {
	utils.ClearScreen()
	goodbyeBox := pterm.DefaultBox.
		WithTitle("👋 Goodbye").
		WithTitleTopCenter().
		WithBoxStyle(pterm.NewStyle(pterm.FgGreen))

	fmt.Println(goodbyeBox.Sprint(
		fmt.Sprintf("Thanks for chatting, %s!\n", ui.username) +
			"You have left the chat room.\n" +
			"See you next time! 🌟"))
}

// UpdateInputBar updates the input bar with current text
func (ui *UI) UpdateInputBar(text string) {
	fmt.Printf("\033[%d;5H", ui.terminalHeight-1)
	fmt.Printf("%s%s%s",
		utils.ColorWhite(text),
		strings.Repeat(" ", ui.terminalWidth-6-len(text)),
		"")
	fmt.Printf("\033[%d;%dH", ui.terminalHeight-1, 5+len(text))
}

// UpdateUserList updates the online users list
func (ui *UI) UpdateUserList(users []string) {
	ui.users = users
}
