package main

import (
	"flag"
	"fmt"
	"terminal-chat/client"
	"terminal-chat/utils"

	"github.com/pterm/pterm"
)

func main() {
	var host = flag.String("host", "localhost", "Host to connect to")
	var port = flag.String("port", "8080", "Port to connect to")
	flag.Parse()

	// Clear screen and show client banner
	utils.ClearScreen()
	showClientBanner()

	fmt.Printf("\n%s Connecting to chat server at %s:%s...\n",
		utils.ColorBlue("🔗"), *host, *port)

	client.StartClient(*host, *port)
}

func showClientBanner() {
	banner := pterm.DefaultBigText.WithLetters(
		pterm.NewLettersFromStringWithStyle("CHAT", pterm.NewStyle(pterm.FgCyan)),
		pterm.NewLettersFromStringWithStyle("CLIENT", pterm.NewStyle(pterm.FgBlue)),
	)
	banner.Render()

	pterm.DefaultCenter.WithCenterEachLineSeparately().Println(
		utils.ColorCyan("Terminal Chat Client\n") +
			utils.ColorBlue("Connect to Real-time Chat Rooms\n"))
}
