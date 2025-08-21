package gifs

import (
	"fmt"
	"strings"
	"terminal-chat/utils"
	"time"
)

// GIFFrame represents a single frame of ASCII art
type GIFFrame struct {
	Content  string
	Duration time.Duration
}

// AnimatedGIF represents a complete ASCII GIF
type AnimatedGIF struct {
	Name   string
	Frames []GIFFrame
	Loop   bool
}

// Popular ASCII GIFs
var GIFLibrary = map[string]AnimatedGIF{
	"wave": {
		Name: "wave",
		Loop: true,
		Frames: []GIFFrame{
			{Content: "👋", Duration: 500 * time.Millisecond},
			{Content: " 👋", Duration: 500 * time.Millisecond},
			{Content: "  👋", Duration: 500 * time.Millisecond},
			{Content: " 👋", Duration: 500 * time.Millisecond},
		},
	},
	"typing": {
		Name: "typing",
		Loop: true,
		Frames: []GIFFrame{
			{Content: "💬", Duration: 800 * time.Millisecond},
			{Content: "💬.", Duration: 300 * time.Millisecond},
			{Content: "💬..", Duration: 300 * time.Millisecond},
			{Content: "💬...", Duration: 300 * time.Millisecond},
		},
	},
	"loading": {
		Name: "loading",
		Loop: true,
		Frames: []GIFFrame{
			{Content: "⠋", Duration: 80 * time.Millisecond},
			{Content: "⠙", Duration: 80 * time.Millisecond},
			{Content: "⠹", Duration: 80 * time.Millisecond},
			{Content: "⠸", Duration: 80 * time.Millisecond},
			{Content: "⠼", Duration: 80 * time.Millisecond},
			{Content: "⠴", Duration: 80 * time.Millisecond},
			{Content: "⠦", Duration: 80 * time.Millisecond},
			{Content: "⠧", Duration: 80 * time.Millisecond},
			{Content: "⠇", Duration: 80 * time.Millisecond},
			{Content: "⠏", Duration: 80 * time.Millisecond},
		},
	},
	"celebration": {
		Name: "celebration",
		Loop: false,
		Frames: []GIFFrame{
			{Content: "🎉", Duration: 200 * time.Millisecond},
			{Content: "✨🎉✨", Duration: 200 * time.Millisecond},
			{Content: "🎊✨🎉✨🎊", Duration: 200 * time.Millisecond},
			{Content: "🎉🎊✨🎉✨🎊🎉", Duration: 300 * time.Millisecond},
			{Content: "✨🎉✨", Duration: 200 * time.Millisecond},
			{Content: "🎉", Duration: 200 * time.Millisecond},
		},
	},
	"cat": {
		Name: "cat",
		Loop: true,
		Frames: []GIFFrame{
			{Content: "/\\_/\\  \n( ^.^ )\n_) (_", Duration: 1000 * time.Millisecond},
			{Content: "/\\_/\\  \n( -.- )\n_) (_", Duration: 500 * time.Millisecond},
			{Content: "/\\_/\\  \n( ^.^ )\n_) (_", Duration: 1000 * time.Millisecond},
			{Content: "/\\_/\\  \n( o.o )\n_) (_", Duration: 300 * time.Millisecond},
		},
	},
	"dance": {
		Name: "dance",
		Loop: true,
		Frames: []GIFFrame{
			{Content: "♪┏(°.°)┛", Duration: 400 * time.Millisecond},
			{Content: "♪┗(°.°)┓", Duration: 400 * time.Millisecond},
			{Content: "♪┏(°.°)┛", Duration: 400 * time.Millisecond},
			{Content: "♪┗(°.°)┓", Duration: 400 * time.Millisecond},
		},
	},
}

// GetGIF returns a GIF by name
func GetGIF(name string) (AnimatedGIF, bool) {
	gif, exists := GIFLibrary[name]
	return gif, exists
}

// GetAvailableGIFs returns list of available GIF names
func GetAvailableGIFs() []string {
	var names []string
	for name := range GIFLibrary {
		names = append(names, name)
	}
	return names
}

// FormatGIFFrame formats a GIF frame for display
func FormatGIFFrame(frame string, username string, timestamp string) string {
	lines := strings.Split(frame, "\n")
	if len(lines) == 1 {
		// Single line frame
		return fmt.Sprintf("%s│ %s %s │ %s %s",
			utils.ColorCyan(""), timestamp, utils.ColorMagenta(username),
			utils.ColorYellow(frame), utils.ColorCyan(""))
	} else {
		// Multi-line frame
		result := ""
		for i, line := range lines {
			if i == 0 {
				result += fmt.Sprintf("%s│ %s %s │ %s%s",
					utils.ColorCyan(""), timestamp, utils.ColorMagenta(username),
					utils.ColorYellow(line), utils.ColorCyan(""))
			} else {
				result += fmt.Sprintf("\n%s│%s │ %s%s",
					utils.ColorCyan(""), strings.Repeat(" ", 15),
					utils.ColorYellow(line), utils.ColorCyan(""))
			}
		}
		return result
	}
}
