package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"encoding/json"
	"io"
)

// Struct Repo describes a git repository
type Repo struct {
	Name   string            `json:"name"`
	Local  string            `json:"local"`
	Remote string            `json:"remote,omitempty"`
	Config map[string]interface{} `json:"config"`
}

func MakeRepo(dir string) (*Repo, error) {
	repo := &Repo{}

	// Check if the directory exists
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return nil, fmt.Errorf("Directory does not exist: %s", dir)
	}
	// Get the absolute path of the directory
	dir, err := filepath.Abs(dir)
	if err != nil {
		return nil, fmt.Errorf("Could not get absolute path of %s: %s", dir, err)
	}
	// Check if the directory is a git repository
	if _, err := os.Stat(filepath.Join(dir, ".git")); os.IsNotExist(err) {
		return nil, fmt.Errorf("Directory is not a git repository: %s", dir)
	}
	// Set the local path of the repository
	repo.Local = dir

	// Set the name of the directory
	repo.Name = filepath.Base(dir)

	// Read the .git/config file and store the key/values in the Config map
	repo.Config = make(map[string]interface{})
	configFile := filepath.Join(dir, ".git", "config")
	file, err := os.Open(configFile)
	if err != nil {
		return nil, fmt.Errorf("Could not open %s: %s", configFile, err)
	}
	defer file.Close()

	var currentSection string
	var currentSubsection string

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if len(line) == 0 || strings.HasPrefix(line, ";") || strings.HasPrefix(line, "#") {
			// Skip empty lines and comments
			continue
		}

		if strings.HasPrefix(line, "[") {
			// New section
			sectionParts := strings.Split(line[1:len(line)-1], " ")
			currentSection = sectionParts[0]
			if len(sectionParts) > 1 {
				currentSubsection = strings.Trim(sectionParts[1], "\"")
			} else {
				currentSubsection = ""
			}
		} else if strings.Contains(line, "=") {
			// Key-value pair
			parts := strings.SplitN(line, "=", 2)
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])

			key = strings.Trim(key, "\"")

			if currentSubsection != "" {
				if _, exists := repo.Config[currentSection]; !exists {
					repo.Config[currentSection] = make(map[string]interface{})
				}
				subsectionMap := repo.Config[currentSection].(map[string]interface{})
				if _, exists := subsectionMap[currentSubsection]; !exists {
					subsectionMap[currentSubsection] = make(map[string]string)
				}
				subsectionConfig := subsectionMap[currentSubsection].(map[string]string)
				subsectionConfig[key] = value
			} else {
				if _, exists := repo.Config[currentSection]; !exists {
					repo.Config[currentSection] = make(map[string]string)
				}
				sectionConfig := repo.Config[currentSection].(map[string]string)
				sectionConfig[key] = value
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("Error reading config file: %s", err)
	}

	// Set the remote URL
	repo.Remote, err = repo.GetConfigValue("remote.origin.url")
	if err != nil {
		return nil, fmt.Errorf("Could not get remote URL: %s", err)
	}

	return repo, nil
}

// Get the value of a key in the Config map
func (r *Repo) GetConfigValue(key string) (string, error) {
	parts := strings.Split(key, ".")
	if len(parts) < 2 {
		return "", fmt.Errorf("Invalid key format. Expected format: section.key or section.subsection.key")
	}

	section := parts[0]
	sectionMap, sectionExists := r.Config[section]
	if !sectionExists {
		return "", fmt.Errorf("Section %s not found", section)
	}

	if len(parts) == 2 {
		// No subsection, just section.key
		key := parts[1]
		value, keyExists := sectionMap.(map[string]string)[key]
		if !keyExists {
			return "", fmt.Errorf("Key %s not found in section %s", key, section)
		}
		return value, nil
	} else if len(parts) == 3 {
		// With subsection, section.subsection.key
		subsection := parts[1]
		key := parts[2]
		subsectionMap, subsectionExists := sectionMap.(map[string]interface{})[subsection]
		if !subsectionExists {
			return "", fmt.Errorf("Subsection %s not found in section %s", subsection, section)
		}
		value, keyExists := subsectionMap.(map[string]string)[key]
		if !keyExists {
			return "", fmt.Errorf("Key %s not found in subsection %s of section %s", key, subsection, section)
		}
		return value, nil
	}

	return "", fmt.Errorf("Invalid key format")
}

// Scan a root folder for repositories and return a Repos slice
func MakeReposFromRoot(root string) ([]Repo, error) {
	var repos []Repo

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			// Check if the directory is a git repository
			if _, err := os.Stat(filepath.Join(path, ".git")); err == nil {
				// Create a Repo from the directory
				repo, err := MakeRepo(path)
				if err != nil {
					fmt.Printf("Error creating repo from %s: %s\n", path, err)
					return nil // Continue walking
				}
				repos = append(repos, *repo)
				// Skip walking into the .git directory itself
				return filepath.SkipDir
			}
		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("Error walking the path %s: %s", root, err)
	}

	return repos, nil
}


// Export the Repos slice to a JSON string
// If includeConfig is false, the Config field is omitted
func ReposToJSON(repos []Repo, includeConfig bool) (string, error) {
	if includeConfig {
		jsonData, err := json.MarshalIndent(repos, "", "  ")
		if err != nil {
			return "", fmt.Errorf("Error marshalling repos to JSON: %s", err)
		}
		return string(jsonData), nil
	}

	type RepoWithoutConfig struct {
		Name   string `json:"name"`
		Local  string `json:"local"`
		Remote string `json:"remote,omitempty"`
	}

	reposWithoutConfig := make([]RepoWithoutConfig, len(repos))
	for i, repo := range repos {
		reposWithoutConfig[i] = RepoWithoutConfig{
			Name:   repo.Name,
			Local:  repo.Local,
			Remote: repo.Remote,
		}
	}

	jsonData, err := json.MarshalIndent(reposWithoutConfig, "", "  ")
	if err != nil {
		return "", fmt.Errorf("Error marshalling repos to JSON: %s", err)
	}

	return string(jsonData), nil
}

// Read a JSON string and return a slice of Repos
func ReposFromJSON(jsonData string) ([]Repo, error) {
	var repos []Repo
	err := json.Unmarshal([]byte(jsonData), &repos)
	if err != nil {
		return nil, fmt.Errorf("Error unmarshalling JSON to repos: %s", err)
	}
	return repos, nil
}

// Load the configuration file
// The configuration file is a JSON file that contains an array of repositories
// The file is located in the OS user's configuration directory, i.e. ~/.config/gogit/repos.json
func LoadReposFromJSON(file string) ([]Repo, error) {
	// Load the repositories from the configuration file
	f, err := os.Open(file)
	if err != nil {
		return nil, fmt.Errorf("Could not open %s: %s", file, err)
	}
	defer f.Close()

	bytes, err := io.ReadAll(f)
	if err != nil {
		return nil, fmt.Errorf("Error reading %s: %s", file, err)
	}

	var repos []Repo
	err = json.Unmarshal(bytes, &repos)
	if err != nil {
		return nil, fmt.Errorf("Error parsing JSON: %s", err)
	}

	// Fill the Config map for each repository
	for i, repo := range repos {
		repo, err := MakeRepo(repo.Local)
		if err != nil {
			return nil, fmt.Errorf("Error creating repo from %s: %s", repo.Local, err)
		}
		repos[i] = *repo
	}

	return repos, nil
}


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
