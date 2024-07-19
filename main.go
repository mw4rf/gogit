package main

import (
	"flag"
	"fmt"
	"os"
)

// Global variable to control debug mode
var debugMode bool



// logDebug prints debug messages if debugMode is enabled
func LogDebug(message string) {
	if debugMode {
		fmt.Println("DEBUG:", message)
	}
}


func main() {
	// Define and parse command line "--flags"
	debugFlag := flag.Bool("debug", false, "Enable debug mode")
	flag.Parse()

	// Set debug mode based on command line flag
	debugMode = *debugFlag

	LogDebug("Debug mode enabled")

	// Load repositories from config.toml file
	repos := LoadRepos()

	// Regular arguments after "--flags"
	command := flag.Arg(0)

	switch command {
		case "init":
			fmt.Println("init command")
		case "status":
			PrintRepos(repos)
		default:
			fmt.Println("gogit v0.1")
			fmt.Println("Usage: gogit [flags] <commands>")
			fmt.Println("Flags:")
			fmt.Println("  --debug Enable debug mode")
			fmt.Println("Commands:")
			fmt.Println("  init   Initialize a new repository")
			fmt.Println("  status Show the status of the repository")
			os.Exit(1)
	}

}
