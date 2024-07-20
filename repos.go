package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// Struct Repo describes a git repository
type Repo struct {
	Name   string            `json:"name"`
	Local  string            `json:"local"`
	Remote string            `json:"remote,omitempty"`
	Config map[string]interface{} `json:"config"`
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

// Load the configuration of a repository
// The configuration is stored in the .git/config file of the repository
// The function reads the file and stores the key-value pairs in the Config map
func (r *Repo) LoadConfig() error {
	configFile := filepath.Join(r.Local, ".git", "config")
	config, err := parseGitConfig(configFile)
	if err != nil {
		return err
	}

	r.Config = config

	// Set the remote URL
	r.Remote, err = r.GetConfigValue("remote.origin.url")
	if err != nil {
		return fmt.Errorf("Could not get remote URL: %s", err)
	}

	return nil
}

// Parse the .git/config file and return the configuration map
func parseGitConfig(configFile string) (map[string]interface{}, error) {
	config := make(map[string]interface{})

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
				if _, exists := config[currentSection]; !exists {
					config[currentSection] = make(map[string]interface{})
				}
				subsectionMap := config[currentSection].(map[string]interface{})
				if _, exists := subsectionMap[currentSubsection]; !exists {
					subsectionMap[currentSubsection] = make(map[string]string)
				}
				subsectionConfig := subsectionMap[currentSubsection].(map[string]string)
				subsectionConfig[key] = value
			} else {
				if _, exists := config[currentSection]; !exists {
					config[currentSection] = make(map[string]string)
				}
				sectionConfig := config[currentSection].(map[string]string)
				sectionConfig[key] = value
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("Error reading config file: %s", err)
	}

	return config, nil
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
	for i := range repos {
		repo := &repos[i]
		err = repo.LoadConfig()
		if err != nil {
			fmt.Println(ColorOutput(ColorYellow, fmt.Sprintf("[Warning]: %s -- Have you cloned this repository? Run <gogit clone>", err)))
			repo.Config = nil // Ensure the Config is nil if it could not be loaded
		}
	}

	return repos, nil
}

// Scan a root folder for repositories and return a Repos slice
// Used by the genrepos command to generate a JSON string with the details of all git repositories in a given root folder
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
				repo, err := MakeRepoFromLocal(path)
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

// Create a Repo from a directory
// Used by the MakeReposFromRoot function
// The function checks if the directory is a git repository and reads the .git/config file
// to get the remote URL and other configuration settings
func MakeRepoFromLocal(dir string) (*Repo, error) {
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
	config, err := parseGitConfig(filepath.Join(dir, ".git", "config"))
	if err != nil {
		return nil, err
	}

	repo.Config = config

	// Set the remote URL
	repo.Remote, err = repo.GetConfigValue("remote.origin.url")
	if err != nil {
		return nil, fmt.Errorf("Could not get remote URL: %s", err)
	}

	return repo, nil
}
