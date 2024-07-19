package main

import (
	"bufio"
	"fmt"
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


func main() {
	repo, err := MakeRepo(".")
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("Repository: %s\n", repo.Name)
	fmt.Printf("Local Path: %s\n", repo.Local)
	fmt.Printf("Remote URL: %s\n", repo.Remote)
	fmt.Printf("Config: %+v\n", repo.Config)
}
