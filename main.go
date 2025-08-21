package main

import (
	"flag"
	"fmt"
	"log"
	"terminal-chat/client"
	"terminal-chat/server"
	"terminal-chat/utils"

	"github.com/pterm/pterm"
)

func main() {
	var mode = flag.String("mode", "", "Mode: server or client")
	var port = flag.String("port", "8080", "Port to run server on")
	var host = flag.String("host", "localhost", "Host to connect to")
	flag.Parse()

	// Clear screen and show banner
	utils.ClearScreen()
	showBanner()

	if *mode == "" {
		*mode = selectMode()
	}

	switch *mode {
	case "server":
		fmt.Printf("\n%s Starting server on port %s...\n",
			utils.ColorGreen("🚀"), *port)
		server.StartServer(*port)
	case "client":
		fmt.Printf("\n%s Connecting to %s:%s...\n",
			utils.ColorBlue("🔗"), *host, *port)
		client.StartClient(*host, *port)
	default:
		log.Fatal("Invalid mode. Use 'server' or 'client'")
	}
}

func showBanner() {
	banner := pterm.DefaultBigText.WithLetters(
		pterm.NewLettersFromStringWithStyle("CHAT", pterm.NewStyle(pterm.FgCyan)),
		pterm.NewLettersFromStringWithStyle("APP", pterm.NewStyle(pterm.FgMagenta)),
	)
	banner.Render()

	pterm.DefaultCenter.WithCenterEachLineSeparately().Println(
		utils.ColorCyan("Terminal-based WebSocket Chat Application\n") +
			utils.ColorYellow("Built with Go, WebSockets & Interactive CLI\n"))
}

func selectMode() string {
	options := []string{"🖥️  Start Server", "👤 Join as Client"}

	prompt := pterm.DefaultInteractiveSelect.
		WithDefaultText("Select mode:").
		WithOptions(options).
		WithDefaultOption("👤 Join as Client")

	selectedOption, _ := prompt.Show()

	if selectedOption == "🖥️  Start Server" {
		return "server"
	}
	return "client"
}
