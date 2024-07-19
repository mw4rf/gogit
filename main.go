package main

import (
	"flag"
	"fmt"
	"os"
)

// Struct Root describes a root path containing repositories
type Root struct {
	Local string
}

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




func main() {

}
