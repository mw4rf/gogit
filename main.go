package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/pelletier/go-toml"
	"gopkg.in/ini.v1"
)

// Global variable to control debug mode
var debugMode bool

// Struct Repo describes a git repository
type Repo struct {
	Name   string
	Remote string
	Config Config
	Root   Root
}

// Struct Config describes the configuration of a git repository
type Config struct {
	RemoteName   string
	RemoteURL    string
	BranchName   string
	BranchRemote string
}

// Struct Root describes a root path containing repositories
type Root struct {
	Local string
}

// Loads the repositories from the config.toml file and returns a list of Repo structs
func LoadRepos() []Repo {
	// Open the config.toml file
	configDir, err := getDefaultDirectory()
	if err != nil {
		log.Fatalf("Error getting default directory: %v", err)
	}
	configFilePath := filepath.Join(configDir, "config.toml")
	file, err := os.Open(configFilePath)
	if err != nil {
		log.Fatalf("Error opening config.toml file: %v", err)
	}
	defer file.Close()

	// Read the contents of the file
	fileContents, err := io.ReadAll(file)
	if err != nil {
		log.Fatalf("Error reading config.toml file: %v", err)
	}

	// Parse the contents of the file
	var data map[string]interface{}
	err = toml.Unmarshal(fileContents, &data)
	if err != nil {
		log.Fatalf("Error parsing config.toml file: %v", err)
	}

	// Create a slice to hold Repo structs
	var repoList []Repo

	// Track roots
	for rootName, rootData := range data {
		rootMap, ok := rootData.(map[string]interface{})
		if !ok {
			logDebug(fmt.Sprintf("Error parsing root data for key: %s", rootName))
			continue
		}
		localPath, ok := rootMap["local"].(string)
		if !ok {
			logDebug(fmt.Sprintf("Local path not found for root: %s", rootName))
			continue
		}
		root := Root{Local: localPath}
		logDebug(fmt.Sprintf("Found root: %s with local path: %s", rootName, localPath))

		for subKey, subValue := range rootMap {
			if subKey == "local" {
				continue
			}
			repoData, ok := subValue.(map[string]interface{})
			if !ok {
				logDebug(fmt.Sprintf("Error parsing repo data for key: %s", subKey))
				continue
			}
			remoteURL, ok := repoData["remote"].(string)
			if !ok {
				logDebug(fmt.Sprintf("Remote URL not found for repo: %s", subKey))
				continue
			}
			repo := Repo{
				Name:   subKey,
				Remote: remoteURL,
				Root:   root,
			}
			repo.Config = GetConfig(repo)
			repoList = append(repoList, repo)
			logDebug(fmt.Sprintf("Found repo: %s under root: %s with remote: %s", subKey, rootName, remoteURL))
		}
	}

	return repoList
}

// Fetches the configuration of a git repository from the .git/config file and returns a Config struct
func GetConfig(repo Repo) Config {
	repoPath := filepath.Join(repo.Root.Local, repo.Name)
	logDebug(fmt.Sprintf("Fetching config for repo: %s at path: %s", repo.Name, repoPath))

	// Open the .git/config file
	file, err := ini.Load(filepath.Join(repoPath, ".git", "config"))
	if err != nil {
		log.Fatalf("Error opening .git/config file for %s: %v", repoPath, err)
	}

	var config Config

	// Get the first remote section
	for _, section := range file.Sections() {
		if strings.HasPrefix(section.Name(), "remote ") {
			config.RemoteName = strings.Trim(section.Name()[7:], "\"")
			config.RemoteURL = section.Key("url").String()
			break
		}
	}

	// Get the first branch section
	for _, section := range file.Sections() {
		if strings.HasPrefix(section.Name(), "branch ") {
			config.BranchName = strings.Trim(section.Name()[7:], "\"")
			config.BranchRemote = section.Key("remote").String()
			break
		}
	}

	return config
}

// logDebug prints debug messages if debugMode is enabled
func logDebug(message string) {
	if debugMode {
		fmt.Println("DEBUG:", message)
	}
}


// Returns the default directory for the application,
// according to the operating system
// e.g. /home/user/.config/gogit
func getDefaultDirectory() (string, error) {
	dir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	dir = filepath.Join(dir, "gogit")

	// Create the directory if it doesn't exist
	err = os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		return "", err
	}

	return dir, nil
}

func main() {
	// Define and parse command line flags
	debugFlag := flag.Bool("debug", false, "Enable debug mode")
	flag.Parse()

	// Set debug mode based on command line flag
	debugMode = *debugFlag

	fmt.Println("gogit v0.1")

	repos := LoadRepos()
	for _, repo := range repos {
		logDebug(fmt.Sprintf("Repo: %s", repo.Name))
		logDebug(fmt.Sprintf("Remote: %s", repo.Remote))
		logDebug(fmt.Sprintf("Local: %s", repo.Root.Local))
		logDebug(fmt.Sprintf("Remote Name: %s", repo.Config.RemoteName))
		logDebug(fmt.Sprintf("Remote URL: %s", repo.Config.RemoteURL))
		logDebug(fmt.Sprintf("Branch Name: %s", repo.Config.BranchName))
		logDebug(fmt.Sprintf("Branch Remote: %s", repo.Config.BranchRemote))
		fmt.Println()
		fmt.Println(repo.ToString())
	}
}
