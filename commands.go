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
func PrintHelp(command ...string) {
	if len(command) == 0 {
		// General help
		fmt.Println(ColorOutput(ColorYellow, "Usage: gogit <command> [arguments]"))
		fmt.Println(ColorOutput(ColorYellow, "Commands:"))

		// Define the widths for each field
		commandWidth := 40

		fmt.Printf("  %-*s %s\n", commandWidth, ColorOutput(ColorCyan, "list"), ColorOutput(ColorWhite, "List the repositories in a simple and compact format"))
		fmt.Printf("  %-*s %s\n", commandWidth, ColorOutput(ColorCyan, "list full"), ColorOutput(ColorWhite, "List the repositories in a detailed format"))
		fmt.Printf("  %-*s %s\n", commandWidth, ColorOutput(ColorCyan, "run <command> [repository]"), ColorOutput(ColorWhite, "Execute a git command on a repository or on all repositories if no repository is provided"))
		fmt.Printf("  %-*s %s\n", commandWidth, ColorOutput(ColorCyan, "do <command> [repository]"), ColorOutput(ColorWhite, "Execute a predefined command on a repository or on all repositories if no repository is provided."))
		fmt.Printf("  %-*s %s\n", commandWidth, "", ColorOutput(ColorWhite, "To show all available commands, use 'gogit do help'"))
		fmt.Printf("  %-*s %s\n", commandWidth, ColorOutput(ColorCyan, "genrepos [root]"), ColorOutput(ColorWhite, "Generate and print a JSON string with the details of all git repositories in a given root folder"))
		fmt.Printf("  %-*s %s\n", commandWidth, ColorOutput(ColorCyan, "clone"), ColorOutput(ColorWhite, "Check all repositories and clone the ones that are missing"))
		fmt.Printf("  %-*s %s\n", commandWidth, ColorOutput(ColorCyan, "help [command]"), ColorOutput(ColorWhite, "Print this help message or detailed help for a specific command"))
	} else {
		// Detailed help for a specific command
		cmd := command[0]
		switch cmd {
		case "list":
			fmt.Println(ColorOutput(ColorYellow, "Usage: gogit list [full]"))
			fmt.Println(ColorOutput(ColorWhite, "List the repositories in a simple and compact format. Use 'full' to list in a detailed format."))
		case "run":
			fmt.Println(ColorOutput(ColorYellow, "Usage: gogit run <command> [repository]"))
			fmt.Println(ColorOutput(ColorWhite, "Execute a git command on a repository or on all repositories if no repository is provided."))
		case "do":
			fmt.Println(ColorOutput(ColorYellow, "Usage: gogit do <command> [repository]"))
			fmt.Println(ColorOutput(ColorWhite, "Show the details of a predefined command on a repository or on all repositories if no repository is provided."))
			fmt.Println(ColorOutput(ColorWhite, "Available predefined commands:"))
			for cmd, args := range predefinedCommands {
				// fmt.Printf("  %s\n", ColorOutput(ColorGreen, cmd))
				fmt.Printf("  %s => %s\n", ColorOutput(ColorGreen, cmd), ColorOutput(ColorBlue, strings.Join(args, " ")))
			}
		case "genrepos":
			fmt.Println(ColorOutput(ColorYellow, "Usage: gogit genrepos [root]"))
			fmt.Println(ColorOutput(ColorWhite, "Generate and print a JSON string with the details of all git repositories in a given root folder."))
		case "clone":
			fmt.Println(ColorOutput(ColorYellow, "Usage: gogit clone"))
			fmt.Println(ColorOutput(ColorWhite, "Check all repositories and clone the ones that are missing."))
		case "help":
			fmt.Println(ColorOutput(ColorYellow, "Usage: gogit help [command]"))
			fmt.Println(ColorOutput(ColorWhite, "Print this help message or detailed help for a specific command."))
		default:
			fmt.Println(ColorOutput(ColorRed, fmt.Sprintf("Error: Unknown command '%s'", cmd)))
			fmt.Println(ColorOutput(ColorWhite, "Use 'gogit help' to see the list of available commands."))
		}
	}
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

// Command: run
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

// Command: do
// Description: Show the details of a git command
// Example: gogit show status myrepo
var predefinedCommands = map[string][]string{
    // History and status
    "h":           {"log", "--oneline", "--decorate", "--graph", "--all"},
    "hf":          {"log", "--oneline", "--decorate", "--graph", "--all", "--simplify-by-decoration"},
    "ha":          {"log", "--author=<author>", "--oneline", "--decorate", "--graph"},
    "hg":          {"log", "--grep=<pattern>", "--oneline", "--decorate", "--graph"},
    "st":          {"status", "-s", "-b"},
    "status":      {"status"},

    // Tag and branch management
    "t":           {"tag", "-l"},
    "ta":          {"tag", "-a", "v1.0", "-m", "Version 1.0"},
    "td":          {"tag", "-d", "v1.0"},
    "b":           {"branch", "-a"},
    "bc":          {"branch", "--show-current"},
    "bn":          {"branch", "new-branch"},
    "bd":          {"branch", "-d", "old-branch"},
    "bf":          {"branch", "-D", "old-branch"},
    "bm":          {"branch", "-m", "old-branch", "new-branch"},

    // Remote management
    "r":           {"remote", "-v"},
    "ra":          {"remote", "add", "origin", "git@github.com:user/repo.git"},
    "rr":          {"remote", "remove", "origin"},
    "rs":          {"remote", "set-url", "origin", "git@github.com:new/repo.git"},
    "rp":          {"remote", "prune", "origin"},

    // Configuration
    "cfg":         {"config", "--list"},
    "cfge":        {"config", "--global", "--edit"},
    "cfgu":        {"config", "--global", "user.name", "Your Name"},
    "cfgm":        {"config", "--global", "user.email", "your.email@example.com"},
    "cfga":        {"config", "--global", "alias.co", "checkout"},

    // Diff commands
    "diff":        {"diff"},
    "diffs":       {"diff", "--stat"},
    "diffc":       {"diff", "--cached"},
    "diffn":       {"diff", "--name-only"},
    "diffsum":     {"diff", "--summary"},
    "diffcolor":   {"diff", "--color"},
    "diffw":       {"diff", "--word-diff"},
    "diffu":       {"diff", "-U10"},

    // Pull commands
    "pull":        {"pull"},
    "pullrb":      {"pull", "--rebase"},
    "pullff":      {"pull", "--ff-only"},
    "pullsq":      {"pull", "--squash"},
    "pullall":     {"pull", "--all"},

    // Push commands
    "push":        {"push"},
    "pushf":       {"push", "--force"},
    "pushfl":      {"push", "--force-with-lease"},
    "pushup":      {"push", "--set-upstream", "origin", "branch"},
    "pushtags":    {"push", "--tags"},
    "pushdel":     {"push", "origin", "--delete", "branch"},

    // Fetch commands
    "fetch":       {"fetch"},
    "fetchp":      {"fetch", "--prune"},
    "fetchall":    {"fetch", "--all"},
    "fetchtags":   {"fetch", "--tags"},

    // Merge and rebase
    "mg":          {"merge"},
    "mgab":        {"merge", "--abort"},
    "mgours":      {"merge", "--strategy=ours"},
    "mgm":         {"merge", "-m", "Merging branch"},
    "rb":          {"rebase"},
    "rbcon":       {"rebase", "--continue"},
    "rbab":        {"rebase", "--abort"},
    "rbmg":        {"rebase", "--merge"},
    "rbi":         {"rebase", "-i", "HEAD~5"},

    // Checkout and commit
    "ck":          {"checkout"},
    "ckn":         {"checkout", "-b", "new-branch"},
    "ckd":         {"checkout", "--detach"},
    "ckf":         {"checkout", "HEAD", "file.txt"},
    "cm":          {"commit"},
    "cmam":        {"commit", "--amend"},
    "cmm":         {"commit", "-m", "Your commit message"},
    "cmv":         {"commit", "--verbose"},
    "cma":         {"commit", "-a", "-m", "Commit all changes"},

    // Add, reset, and remove
    "add":         {"add"},
    "adda":        {"add", "--all"},
    "addi":        {"add", "-i"},
    "addp":        {"add", "-p"},
    "reset":       {"reset"},
    "resets":      {"reset", "--soft", "HEAD~1"},
    "reseth":      {"reset", "--hard", "HEAD~1"},
    "resetm":      {"reset", "--mixed", "HEAD~1"},
    "rm":          {"rm"},
    "rmc":         {"rm", "--cached"},
    "rmf":         {"rm", "-f"},

    // Logging and blame
    "log":         {"log"},
    "logg":        {"log", "--graph"},
    "logd":        {"log", "--decorate"},
    "loga":        {"log", "--author=<author>"},
    "loggrep":     {"log", "--grep=<pattern>"},
    "blame":       {"blame"},
    "blamei":      {"blame", "--incremental"},
    "blametime":   {"blame", "--since=2.weeks"},
    "blamel":      {"blame", "--line-porcelain"},

    // Cleaning and garbage collection
    "clean":       {"clean"},
    "cleanf":      {"clean", "-f"},
    "cleandx":     {"clean", "-d", "-x", "-f"},
    "gc":          {"gc", "--aggressive", "--prune=now"},
    "gcauto":      {"gc", "--auto"},

    // Submodule commands
    "subupd":      {"submodule", "update", "--init", "--recursive"},
    "substatus":   {"submodule", "status"},
    "subadd":      {"submodule", "add", "url", "path"},
    "subsync":     {"submodule", "sync", "--recursive"},

    // Other commands
    "shortlog":    {"shortlog", "-sn"},
    "ignore":      {"check-ignore", "*"},
    "revlist":     {"rev-list", "--all"},
    "reflog":      {"reflog"},
    "countobj":    {"count-objects", "-v"},
    "showbranch":  {"show-branch"},
    "verifypack":  {"verify-pack", "-v", ".git/objects/pack/*.pack"},
    "show":        {"show"},
    "grep":        {"grep", "--line-number", "TODO"},
    "archivelog":  {"archive", "--format=tar", "--output=log.tar", "HEAD"},
    "bundlecreate": {"bundle", "create", "repo.bundle", "HEAD"},
    "bundleverify": {"bundle", "verify", "repo.bundle"},
    "bundleheads": {"bundle", "list-heads", "repo.bundle"},
    "rangediff":   {"range-diff", "HEAD~5..HEAD", "origin/master"},
    "sparse":      {"sparse-checkout", "init", "--cone"},
    "worktreeadd": {"worktree", "add", "-b", "new-branch", "../path/to/new-worktree"},
    "fsck":        {"fsck", "--full"},
    "packrefs":    {"pack-refs", "--all"},
    "prune":       {"prune"},
    "bisectstart": {"bisect", "start"},
    "bisectbad":   {"bisect", "bad"},
    "bisectgood":  {"bisect", "good", "HEAD~10"},
    "bisectreset": {"bisect", "reset"},
    "repack":      {"repack", "-a", "-d", "--depth=250", "--window=250"},
    "verifytag":   {"verify-tag", "-v"},
    "verifycm":    {"verify-commit", "-v"},
    "lstree":      {"ls-tree", "-r", "HEAD"},
    "revparse":    {"rev-parse", "--verify", "HEAD"},
    "cherry":      {"cherry", "-v"},
    "cherrypick":  {"cherry-pick", "HEAD~3..HEAD"},
    "notes":       {"notes", "list"},
    "describetags": {"describe", "--tags", "--abbrev=0"},
    "checkoutindex": {"checkout-index", "-a", "-f"},
    "committree":  {"commit-tree", "HEAD^{tree}", "-m", "New commit"},
    "mergebase":   {"merge-base", "HEAD", "master"},
    "packobj":     {"pack-objects", "--all", ".git/objects/pack/pack"},
    "revparsehead": {"rev-parse", "--short", "HEAD"},
    "symbolicref": {"symbolic-ref", "HEAD"},
    "updateindex": {"update-index", "--refresh"},
    "updateref":   {"update-ref", "-d", "refs/heads/branch"},
    "whatchanged": {"whatchanged", "-p", "--abbrev-commit", "--pretty=medium"},
    "verifypackfiles": {"verify-pack", "-v", ".git/objects/pack/pack-*"},
    "unpackobj":   {"unpack-objects", ".git/objects/pack/*.pack"},
    "difftool":    {"difftool", "--dir-diff"},
    "mergetool":   {"mergetool"},
    "subtreesplit": {"subtree", "split", "--prefix=lib", "-b", "split-branch"},
    "filterbranch": {"filter-branch", "--index-filter", "git rm -r --cached --ignore-unmatch <path>", "HEAD"},
    "replace":     {"replace", "old-hash", "new-hash"},
    "showref":     {"show-ref"},
    "verifynotes": {"verify-notes", "-v"},
    "commitgraph": {"commit-graph", "write", "--reachable", "--changed-paths"},
    "worktreeprune": {"worktree", "prune"},
}

func DoCommand(repos []Repo, args []string, repoName string) {
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
