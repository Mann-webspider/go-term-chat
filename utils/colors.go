package utils

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"

	"github.com/fatih/color"
)

var (
	ColorRed     = color.New(color.FgRed).SprintFunc()
	ColorGreen   = color.New(color.FgGreen).SprintFunc()
	ColorYellow  = color.New(color.FgYellow).SprintFunc()
	ColorBlue    = color.New(color.FgBlue).SprintFunc()
	ColorMagenta = color.New(color.FgMagenta).SprintFunc()
	ColorCyan    = color.New(color.FgCyan).SprintFunc()
	ColorWhite   = color.New(color.FgWhite).SprintFunc()
	ColorBold    = color.New(color.Bold).SprintFunc()
)

// Color combinations
var (
	ColorSuccess = color.New(color.FgGreen, color.Bold).SprintFunc()
	ColorError   = color.New(color.FgRed, color.Bold).SprintFunc()
	ColorInfo    = color.New(color.FgBlue, color.Bold).SprintFunc()
	ColorWarning = color.New(color.FgYellow, color.Bold).SprintFunc()
)

// Background colors
var (
	BgRed     = color.New(color.BgRed, color.FgWhite).SprintFunc()
	BgGreen   = color.New(color.BgGreen, color.FgBlack).SprintFunc()
	BgYellow  = color.New(color.BgYellow, color.FgBlack).SprintFunc()
	BgBlue    = color.New(color.BgBlue, color.FgWhite).SprintFunc()
	BgMagenta = color.New(color.BgMagenta, color.FgWhite).SprintFunc()
	BgCyan    = color.New(color.BgCyan, color.FgBlack).SprintFunc()
)

// ClearScreen clears the terminal screen
func ClearScreen() {
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/c", "cls")
	} else {
		cmd = exec.Command("clear")
	}
	cmd.Stdout = os.Stdout
	cmd.Run()
}

// PrintBox prints text in a colored box
func PrintBox(text string, colorFunc func(...interface{}) string) {
	boxWidth := len(text) + 4
	border := "+" + string(rune('-')) + "+"
	for i := 0; i < boxWidth-2; i++ {
		border = border[:1] + "-" + border[2:]
	}

	fmt.Println(colorFunc(border))
	fmt.Println(colorFunc(fmt.Sprintf("| %s |", text)))
	fmt.Println(colorFunc(border))
}

// GetRandomColor returns a random color function
func GetRandomColor(index int) func(...interface{}) string {
	colors := []func(...interface{}) string{
		ColorRed, ColorGreen, ColorYellow, ColorBlue,
		ColorMagenta, ColorCyan,
	}
	return colors[index%len(colors)]
}
