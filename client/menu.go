package client

import (
	"fmt"
	"os"
	"strings"
	"terminal-chat/utils"

	"github.com/manifoldco/promptui"
	"github.com/pterm/pterm"
)

// MenuOption represents a menu option
type MenuOption struct {
	Label       string
	Description string
	Action      func() error
}

// ShowMainMenu displays the main application menu
func ShowMainMenu() (*MenuOption, error) {
	utils.ClearScreen()

	// Show welcome banner
	pterm.DefaultBigText.WithLetters(
		pterm.NewLettersFromStringWithStyle("MENU", pterm.NewStyle(pterm.FgCyan)),
	).Render()

	options := []*MenuOption{
		{
			Label:       "🚀 Join Chat Room",
			Description: "Connect to a chat room and start messaging",
		},
		{
			Label:       "📋 View Available Rooms",
			Description: "See list of active chat rooms",
		},
		{
			Label:       "⚙️  Settings",
			Description: "Configure application settings",
		},
		{
			Label:       "❓ Help",
			Description: "View help and instructions",
		},
		{
			Label:       "🚪 Exit",
			Description: "Exit the application",
		},
	}

	templates := &promptui.SelectTemplates{
		Label:    "{{ . | cyan | bold }}",
		Active:   "▶ {{ .Label | green | bold }} - {{ .Description | faint }}",
		Inactive: "  {{ .Label | faint }} - {{ .Description | faint }}",
		Selected: "✓ {{ .Label | green | bold }}",
	}

	prompt := promptui.Select{
		Label:     "Select an option",
		Items:     options,
		Templates: templates,
		Size:      5,
	}

	index, _, err := prompt.Run()
	if err != nil {
		return nil, err
	}

	return options[index], nil
}

// GetUserInput prompts for user input with validation
func GetUserInput() (string, string, error) {
	// Username input with validation
	usernamePrompt := promptui.Prompt{
		Label: "Enter your username",
		Validate: func(input string) error {
			if len(input) < 1 {
				return fmt.Errorf("username cannot be empty")
			}
			if len(input) > 20 {
				return fmt.Errorf("username too long (max 20 characters)")
			}
			return nil
		},
		Templates: &promptui.PromptTemplates{
			Prompt:  "{{ . | bold | blue }}{{ \":\" | bold | blue }} ",
			Valid:   "{{ . | bold | blue }}{{ \":\" | bold | blue }} {{ . | green }}",
			Invalid: "{{ . | bold | blue }}{{ \":\" | bold | blue }} {{ . | red }}",
			Success: "{{ . | bold | blue }}{{ \":\" | bold | blue }} {{ . | green | bold }}",
		},
	}

	username, err := usernamePrompt.Run()
	if err != nil {
		return "", "", err
	}

	// Clean room selection without emojis that might cause encoding issues
	rooms := []string{
		"general",
		"tech",
		"gaming",
		"books",
		"music",
		"random",
		"Create new room",
	}

	roomPrompt := promptui.Select{
		Label: "Select a chat room",
		Items: rooms,
		Templates: &promptui.SelectTemplates{
			Label:    "{{ . | cyan | bold }}",
			Active:   "▶ {{ . | green | bold }}",
			Inactive: "  {{ . | faint }}",
			Selected: "✓ {{ . | green | bold }}",
		},
	}

	_, selectedRoom, err := roomPrompt.Run()
	if err != nil {
		return "", "", err
	}

	room := selectedRoom // No emoji prefix to remove

	// Handle custom room creation...
	return username, room, nil
}

// ShowConnectionMenu displays connection options
// ShowConnectionMenu displays connection options
func ShowConnectionMenu() (string, error) {
	options := []string{
		"🏠 localhost:8080 (Local)",
		"🌐 LAN Server (Enter IP)",
		"🔗 Custom server",
	}

	prompt := promptui.Select{
		Label: "Select server to connect",
		Items: options,
		Templates: &promptui.SelectTemplates{
			Label:    "{{ . | cyan | bold }}",
			Active:   "▶ {{ . | green | bold }}",
			Inactive: "  {{ . | faint }}",
			Selected: "✓ {{ . | green | bold }}",
		},
	}

	index, _, err := prompt.Run()
	if err != nil {
		return "", err
	}

	switch index {
	case 0:
		return "localhost:8080", nil
	case 1:
		return getLANServerAddress()
	case 2:
		return getCustomServerAddress()
	}

	return "", fmt.Errorf("invalid selection")
}

// getLANServerAddress prompts for LAN server IP
func getLANServerAddress() (string, error) {
	ipPrompt := promptui.Prompt{
		Label: "Enter server IP address",
		Validate: func(input string) error {
			if len(input) < 1 {
				return fmt.Errorf("IP address cannot be empty")
			}
			// Basic IP validation
			parts := strings.Split(input, ".")
			if len(parts) != 4 {
				return fmt.Errorf("invalid IP format (use x.x.x.x)")
			}
			return nil
		},
	}

	ip, err := ipPrompt.Run()
	if err != nil {
		return "", err
	}

	portPrompt := promptui.Prompt{
		Label:   "Enter port (default: 8080)",
		Default: "8080",
	}

	port, err := portPrompt.Run()
	if err != nil {
		return "", err
	}

	return ip + ":" + port, nil
}

// getCustomServerAddress prompts for custom server
func getCustomServerAddress() (string, error) {
	serverPrompt := promptui.Prompt{
		Label: "Enter server address (host:port)",
		Validate: func(input string) error {
			if len(input) < 1 {
				return fmt.Errorf("server address cannot be empty")
			}
			if !strings.Contains(input, ":") {
				return fmt.Errorf("address must include port (host:port)")
			}
			return nil
		},
	}

	return serverPrompt.Run()
}

// ShowHelp displays help information
func ShowHelp() {
	utils.ClearScreen()

	helpText := `
🎯 CHAT APPLICATION HELP

📝 COMMANDS:
  • Type messages and press Enter to send
  • /quit or /exit - Leave the chat
  • /help - Show this help
  • /users - List online users
  • /clear - Clear the screen

🎨 FEATURES:
  • Real-time messaging
  • Multiple chat rooms
  • Colorful user interface
  • User presence indicators
  • Message timestamps

🔧 SHORTCUTS:
  • Ctrl+C - Exit application
  • Arrow keys - Navigate menus
  • Enter - Select option
  • Tab - Auto-complete (where available)

💡 TIPS:
  • Choose unique usernames
  • Be respectful to other users
  • Use appropriate room names
  • Keep messages concise

Press Enter to return to main menu...
`

	fmt.Print(utils.ColorCyan(helpText))
	fmt.Scanln()
}

// ShowSettings displays settings menu
func ShowSettings() error {
	options := []string{
		"🎨 Change Color Theme",
		"🔔 Notification Settings",
		"📱 Display Preferences",
		"🔙 Back to Main Menu",
	}

	prompt := promptui.Select{
		Label: "Settings",
		Items: options,
	}

	index, _, err := prompt.Run()
	if err != nil {
		return err
	}

	switch index {
	case 0:
		fmt.Println(utils.ColorInfo("🎨 Color themes will be available in future updates!"))
	case 1:
		fmt.Println(utils.ColorInfo("🔔 Notification settings coming soon!"))
	case 2:
		fmt.Println(utils.ColorInfo("📱 Display preferences will be added!"))
	case 3:
		return nil
	}

	fmt.Println("\nPress Enter to continue...")
	fmt.Scanln()
	return nil
}

// ConfirmExit asks user to confirm exit
func ConfirmExit() bool {
	prompt := promptui.Prompt{
		Label:     "Are you sure you want to exit? (y/N)",
		IsConfirm: true,
	}

	result, err := prompt.Run()
	if err != nil {
		return false
	}

	return result == "y" || result == "Y"
}

// HandleMenuSelection handles the selected menu option
func HandleMenuSelection(option *MenuOption) error {
	switch option.Label {
	case "🚀 Join Chat Room":
		return nil // Will be handled in main client logic
	case "📋 View Available Rooms":
		fmt.Println(utils.ColorInfo("📋 Available rooms feature coming soon!"))
		fmt.Println("Press Enter to continue...")
		fmt.Scanln()
		return fmt.Errorf("back_to_menu")
	case "⚙️  Settings":
		ShowSettings()
		return fmt.Errorf("back_to_menu")
	case "❓ Help":
		ShowHelp()
		return fmt.Errorf("back_to_menu")
	case "🚪 Exit":
		if ConfirmExit() {
			fmt.Println(utils.ColorSuccess("👋 Goodbye!"))
			os.Exit(0)
		}
		return fmt.Errorf("back_to_menu")
	}
	return nil
}
