package main

import (
	"fmt"
	"os"
	"strings"
	"sync"
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
	fmt.Printf("  %-*s %s\n", commandWidth, ColorOutput(ColorCyan, "do <command> [repository]"), "Execute a git command on a repository or on all repositories if no repository is provided")
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

// Command: do
// Description: Execute a git command on all repositories
// This function runs the git command in parallel for each repository with goroutines
// Example: gogit do pull
func ExecGitCommand(repos []Repo, args []string, repoName string) {
    if len(repos) == 0 {
        fmt.Println(ColorOutput(ColorYellow, "No repositories found"))
        os.Exit(0)
    }
    if len(args) == 0 {
        fmt.Println(ColorOutput(ColorRed, "Error: Missing command to execute"))
        fmt.Println(ColorOutput(ColorYellow, "Usage: gogit do <command> [args] [repo_name]"))
        os.Exit(1)
    }

    var wg sync.WaitGroup
    var mu sync.Mutex

    argsStr := strings.Join(args, " ")

    // Filter repositories if a specific repository name is provided
    filteredRepos := repos
    if repoName != "" {
        filteredRepos = []Repo{}
        for _, repo := range repos {
            if repo.Name == repoName {
                filteredRepos = append(filteredRepos, repo)
                break
            }
        }
        if len(filteredRepos) == 0 {
            fmt.Println(ColorOutput(ColorRed, fmt.Sprintf("Error: Repository '%s' not found", repoName)))
            os.Exit(1)
        }
    }

    for _, repo := range filteredRepos {
        wg.Add(1)
        go func(repo Repo) {
            defer wg.Done()

            mu.Lock()
            fmt.Println(ColorOutput(ColorCyan, "======================================="))
            fmt.Println(ColorOutput(ColorCyan, fmt.Sprintf("Executing '%s' in %s", argsStr, repo.Local)))
            fmt.Println(ColorOutput(ColorCyan, "---------------------------------------"))
            mu.Unlock()

            // Run the Git command
            err := repo.RunGitCommand(args)

            mu.Lock()
            if err != nil {
                fmt.Println(ColorOutput(ColorRed, fmt.Sprintf("Error executing command in %s: %s", repo.Name, err)))
            } else {
                fmt.Println(ColorOutput(ColorGreen, fmt.Sprintf("Successfully executed command in %s", repo.Name)))
            }
            fmt.Println(ColorOutput(ColorCyan, "=======================================\n"))
            mu.Unlock()
        }(repo)
    }

    wg.Wait()
    os.Exit(0)
}

// Command: show
// Description: Show the details of a git command
// Example: gogit show status myrepo
var predefinedCommands = map[string][]string{
    "history": {"log", "--oneline", "--decorate", "--graph", "--all"},
    // Add more predefined commands here
}
func ShowCommand(repos []Repo, args []string, repoName string) {
    if len(repos) == 0 {
        fmt.Println(ColorOutput(ColorYellow, "No repositories found"))
        os.Exit(0)
    }
    if len(args) == 0 {
        fmt.Println(ColorOutput(ColorRed, "Error: Missing command to execute"))
        fmt.Println(ColorOutput(ColorYellow, "Usage: gogit show <command> [repo_name]"))
        os.Exit(1)
    }

    // Get the predefined command
    cmdArgs, exists := predefinedCommands[args[0]]
    if !exists {
        fmt.Println(ColorOutput(ColorRed, fmt.Sprintf("Error: Unknown command '%s'", args[0])))
        os.Exit(1)
    }

    // Filter repositories if a specific repository name is provided
    filteredRepos := repos
    if repoName != "" {
        filteredRepos = []Repo{}
        for _, repo := range repos {
            if repo.Name == repoName {
                filteredRepos = append(filteredRepos, repo)
                break
            }
        }
        if len(filteredRepos) == 0 {
            fmt.Println(ColorOutput(ColorRed, fmt.Sprintf("Error: Repository '%s' not found", repoName)))
            os.Exit(1)
        }
    }

    var wg sync.WaitGroup
    var mu sync.Mutex

    for _, repo := range filteredRepos {
        wg.Add(1)
        go func(repo Repo) {
            defer wg.Done()

            mu.Lock()
            fmt.Println(ColorOutput(ColorCyan, "======================================="))
            fmt.Println(ColorOutput(ColorCyan, fmt.Sprintf("Details for %s", repo.Name)))
            fmt.Println(ColorOutput(ColorCyan, "---------------------------------------"))
            mu.Unlock()

            // Run the Git command
            err := repo.RunGitCommand(cmdArgs)

            mu.Lock()
            if err != nil {
                fmt.Println(ColorOutput(ColorRed, fmt.Sprintf("Error executing command in %s: %s", repo.Name, err)))
            }
            fmt.Println(ColorOutput(ColorCyan, "=======================================\n"))
            mu.Unlock()
        }(repo)
    }

    wg.Wait()
    os.Exit(0)
}
