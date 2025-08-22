package main

import (
	"flag"
	"fmt"
	"terminal-chat/server"
	"terminal-chat/utils"

	"github.com/pterm/pterm"
)

func main() {
	var port = flag.String("port", "8080", "Port to run server on")
	flag.Parse()

	// Clear screen and show server banner
	utils.ClearScreen()
	showServerBanner()

	fmt.Printf("\n%s Starting terminal chat server on port %s...\n",
		utils.ColorGreen("🚀"), *port)

	server.StartServer(*port)
}

func showServerBanner() {
	banner := pterm.DefaultBigText.WithLetters(
		pterm.NewLettersFromStringWithStyle("CHAT", pterm.NewStyle(pterm.FgCyan)),
		pterm.NewLettersFromStringWithStyle("SERVER", pterm.NewStyle(pterm.FgGreen)),
	)
	banner.Render()

	pterm.DefaultCenter.WithCenterEachLineSeparately().Println(
		utils.ColorCyan("Terminal Chat Server\n") +
			utils.ColorGreen("WebSocket-based Real-time Communication\n"))
}
