package main

import (
	"fmt"
	"strings"
)

// ANSI escape codes for colors
const (
	colorReset  = "\033[0m"
	colorBright = "\033[1m"
	colorDim    = "\033[2m"
	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorBlue   = "\033[34m"
	colorPurple = "\033[35m"
	colorCyan   = "\033[36m"
)

// Pretty print a Repo struct
func (r Repo) ToString() string {
	repoInfo := fmt.Sprintf("%sRepo: %s%s%s %s%s[%s]%s │ %sLocal: %s%s%s │ %sRemote: %s%s%s",
		colorBright, colorRed, r.Name, colorReset,
		colorBright, colorYellow, r.Config.BranchName, colorReset,
		colorBright, colorGreen, r.Root.Local + r.Name, colorReset,
		colorBright, colorCyan, r.Config.RemoteURL, colorReset,
	)

	// Create border
	border := strings.Repeat("─", 80)

	return fmt.Sprintf(
		"%s\n%s\n",
		border,
		repoInfo,
	)
}
