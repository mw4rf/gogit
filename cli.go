package main

import (
	"fmt"
)

const (
	ColorRed = "31"
	ColorGreen = "32"
	ColorYellow = "33"
	ColorBlue = "34"
	ColorMagenta = "35"
	ColorCyan = "36"
	ColorWhite = "37"
)


func ColorOutput(color string, message string) string {
	return fmt.Sprintf("\033[%sm%s\033[0m", color, message)
}

// Print the details of a Repo on a single line without the Config field
func PrintRepoSimple(repo *Repo) {
	repoNameWidth := 30
	localPathWidth := 50
	remoteURLWidth := 50

	fmt.Printf("%s: %-*s %s: %-*s %s: %-*s\n",
		ColorOutput(ColorCyan, "Repo"), repoNameWidth, ColorOutput(ColorGreen, repo.Name),
		ColorOutput(ColorCyan, "Local"), localPathWidth, ColorOutput(ColorGreen, repo.Local),
		ColorOutput(ColorCyan, "Remote"), remoteURLWidth, ColorOutput(ColorGreen, repo.Remote))
}


// Print the details of a Repo
// PrintRepo prints the details of a Repo with colors and separations
func PrintRepo(repo *Repo) {
	fmt.Printf("%s: %s\n", ColorOutput(ColorCyan, "Repository"), ColorOutput(ColorGreen, repo.Name))
	fmt.Printf("%s: %s\n", ColorOutput(ColorCyan, "Local Path"), ColorOutput(ColorGreen, repo.Local))
	if repo.Remote != "" {
		fmt.Printf("%s: %s\n", ColorOutput(ColorCyan, "Remote URL"), ColorOutput(ColorGreen, repo.Remote))
	}
	fmt.Println(ColorOutput(ColorCyan, "Config:"))
	for section, settings := range repo.Config {
		fmt.Printf("  %s\n", ColorOutput(ColorMagenta, section))
		switch settings := settings.(type) {
		case map[string]string:
			for key, value := range settings {
				fmt.Printf("    %s: %s\n", ColorOutput(ColorYellow, key), ColorOutput(ColorGreen, value))
			}
		case map[string]interface{}:
			for subsection, subsettings := range settings {
				fmt.Printf("    %s\n", ColorOutput(ColorBlue, subsection))
				if subsettings, ok := subsettings.(map[string]string); ok {
					for key, value := range subsettings {
						fmt.Printf("      %s: %s\n", ColorOutput(ColorYellow, key), ColorOutput(ColorGreen, value))
					}
				}
			}
		}
	}
}
