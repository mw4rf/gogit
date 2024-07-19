package main

import (
	"fmt"
	"os"
)

// Command: help
// (Default if no command is provided)
// Print the help message
func PrintHelp() {
	fmt.Println("Usage: gogit <command> [args]")
	fmt.Println("Commands:")
	fmt.Println("  genrepos <root> - Generate and print a JSON string with the details of all git repositories in a given root folder")
	fmt.Println("  help - Print this help message")
}

// Command: genrepos
// Description: Generate and print a JSON string with the details of all git repositories in a given root folder
// Example: gogit genrepos /path/to/root
func GenRepos(root string) {
	repos, err := MakeReposFromRoot(root)
	if err != nil {
		fmt.Println(ColorOutput(ColorRed, fmt.Sprintf("Error generating repositories: %s", err)))
		os.Exit(1)
	}
	// Print the JSON string with the details of the repositories
	jsonData, err := ReposToJSON(repos, false)
	if err != nil {
		fmt.Println(ColorOutput(ColorRed, fmt.Sprintf("Error generating JSON: %s", err)))
		os.Exit(1)
	}
	fmt.Println(jsonData)
	os.Exit(0)
}

// Command: list
// Description: List the repositories
// Example: gogit list
func PrintReposList(repos []Repo, simpleOutput bool) {
	if len(repos) == 0 {
		fmt.Println(ColorOutput(ColorYellow, "No repositories found"))
		os.Exit(0)
	}
	for _, repo := range repos {
		if simpleOutput {
			PrintRepoSimple(&repo)
		} else {
			separator := ColorOutput(ColorRed, "----------------------------------------")
			fmt.Println(separator)
			PrintRepo(&repo)
		}
	}
	os.Exit(0)
}
