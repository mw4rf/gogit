package main

import (
	"fmt"
	"os"
	"path/filepath"
)



func main() {
	// Action that must be executed before loading the repositories

	// No argument
	if len(os.Args) < 2 {
		fmt.Println(ColorOutput(ColorYellow, "gogit - A simple git repository manager"))
		fmt.Println("Usage: gogit <command> [args]")
		fmt.Println(fmt.Sprintf("Use '%s' to see the list of available commands.", ColorOutput(ColorGreen, "gogit help")))
		os.Exit(0)
	}

	// Command: help
	if os.Args[1] == "help" {
		PrintHelp()
		os.Exit(0)
	}

	// Command: genrepos
	// Description: Generate and print a JSON string with the details of all git repositories in a given root folder
	// Example: gogit genrepos /path/to/root
	if os.Args[1] == "genrepos" {
		if len(os.Args) < 3 {
			fmt.Println(ColorOutput(ColorRed, "Error: Missing root folder argument"))
			fmt.Println("Usage: gogit genrepos /path/to/root")
			os.Exit(1)
		}
		root := os.Args[2]
		GenRepos(root)
	}

	// Load the repositories from the configuration file
	reposFile := filepath.Join(GetUserConfigDir(), "repos.json")
	repos, err := LoadReposFromJSON(reposFile)
	if err != nil {
		fmt.Println(ColorOutput(ColorRed, fmt.Sprintf("Error loading repositories: %s", err)))
		fmt.Println("Please make sure the configuration file exists and is valid.")
		fmt.Println("The configuration file should be a JSON file that contains an array of repositories.")
		fmt.Println("It should be located in the OS user's configuration directory, i.e. ~/.config/gogit/repos.json")
		os.Exit(1)
	}

	// Handle commands that require the repositories
	switch os.Args[1] {
	case "list":
		simpleOutput := true
		if len(os.Args) > 2 && os.Args[2] == "full" {
			simpleOutput = false
		}
		PrintReposList(repos, simpleOutput)
	default:
		fmt.Println(ColorOutput(ColorRed, fmt.Sprintf("Error: Unknown command '%s'", os.Args[1])))
		fmt.Println(fmt.Sprintf("Use '%s' to see the list of available commands.", ColorOutput(ColorGreen, "gogit help")))
	}

}
