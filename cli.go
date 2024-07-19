package main

import (
	"fmt"
)

const (
	ColorRed = "31"
	ColorGreen = "32"
	ColorYellow = "33"
	ColorBlue = "34"
	ColorMagenta = "35"
	ColorCyan = "36"
)


func ColorOutput(color string, message string) string {
	return fmt.Sprintf("\033[%sm%s\033[0m", color, message)
}
