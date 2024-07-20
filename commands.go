package main

import (
	"fmt"
	"os"
)

// Command: help
// (Default if no command is provided)
// Print the help message
// PrintHelp prints the usage and commands with aligned columns and colors
// PrintHelp prints the usage and commands with aligned columns and colors
func PrintHelp() {
	fmt.Println(ColorOutput(ColorYellow, "Usage: gogit <command> [arguments]"))
	fmt.Println(ColorOutput(ColorYellow, "Commands:"))

	// Define the widths for each field
	commandWidth := 30

	fmt.Printf("  %-*s %s\n", commandWidth, ColorOutput(ColorCyan, "list"), "List the repositories in a simple and compact format")
	fmt.Printf("  %-*s %s\n", commandWidth, ColorOutput(ColorCyan, "list full"), "List the repositories in a detailed format")
	fmt.Printf("  %-*s %s\n", commandWidth, ColorOutput(ColorCyan, "genrepos [root]"), "Generate and print a JSON string with the details of all git repositories in a given root folder")
	fmt.Printf("  %-*s %s\n", commandWidth, ColorOutput(ColorCyan, "clone"), "Check all repositories and clone the ones that are missing")
	fmt.Printf("  %-*s %s\n", commandWidth, ColorOutput(ColorCyan, "help"), "Print this help message")
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

// Command: clone
// Description: Check all repositories and clone the ones that are missing
// Example: gogit clone
func CloneRepos(repos []Repo) {
	if len(repos) == 0 {
		fmt.Println(ColorOutput(ColorYellow, "No repositories found"))
		os.Exit(0)
	}
	for _, repo := range repos {
		if _, err := os.Stat(repo.Local); os.IsNotExist(err) {
			fmt.Printf("Cloning %s into %s\n", repo.Remote, repo.Local)
			err := repo.Clone()
			if err != nil {
				fmt.Println(ColorOutput(ColorRed, fmt.Sprintf("Error cloning %s: %s", repo.Name, err)))
			}
		} else {
			fmt.Println(ColorOutput(ColorYellow, fmt.Sprintf("Skipping %s: repository already exists", repo.Name)))
		}
	}
	os.Exit(0)
}
