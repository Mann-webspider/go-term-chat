package server

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"terminal-chat/utils"

	"github.com/pterm/pterm"
)

// StartServer starts the WebSocket server
func StartServer(port string) {
	hub := NewHub()
	go hub.Run()

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		ServeWS(hub, w, r)
	})

	// Health check endpoint
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// Show server info
	showServerInfo(port)

	// CHANGED: Bind to all interfaces (0.0.0.0) instead of localhost
	addr := "0.0.0.0:" + port
	fmt.Printf("\n%s Server listening on %s\n",
		utils.ColorSuccess("🚀"), utils.ColorBold("all interfaces:"+port))
	fmt.Printf("%s WebSocket endpoint: %s\n",
		utils.ColorInfo("🔗"), utils.ColorBold("ws://<your-ip>:"+port+"/ws"))
	fmt.Printf("%s Local access: %s\n",
		utils.ColorInfo("🏠"), utils.ColorBold("ws://localhost:"+port+"/ws"))
	fmt.Printf("%s Press %s to stop the server\n\n",
		utils.ColorWarning("⚠️"), utils.ColorBold("Ctrl+C"))

	log.Fatal(http.ListenAndServe(addr, nil))
}

func showServerInfo(port string) {
	// Get local IP address
	localIP := getLocalIP()

	// Create server status panel
	panel := pterm.DefaultPanel.
		WithPanels([][]pterm.Panel{
			{
				{Data: pterm.DefaultCenter.Sprint("🖥️  SERVER STATUS")},
			},
			{
				{Data: fmt.Sprintf("Port: %s", utils.ColorGreen(port))},
				{Data: fmt.Sprintf("Status: %s", utils.ColorSuccess("RUNNING"))},
			},
			{
				{Data: fmt.Sprintf("Local IP: %s", utils.ColorBlue(localIP))},
				{Data: fmt.Sprintf("Network: %s", utils.ColorCyan("ENABLED"))},
			},
			{
				{Data: fmt.Sprintf("Local: %s", utils.ColorYellow("localhost:"+port))},
				{Data: fmt.Sprintf("Remote: %s", utils.ColorMagenta(localIP+":"+port))},
			},
		}).
		WithPadding(1)

	panel.Render()
}

// Add this helper function to get local IP
func getLocalIP() string {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return "Unable to determine"
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return localAddr.IP.String()
}
