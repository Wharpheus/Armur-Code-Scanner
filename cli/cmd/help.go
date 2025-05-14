package cmd

import (
	"fmt"
	"os"

	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
)

var (
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FF5F87")).
			Padding(0, 1)

	commandStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#00FF9F"))

	descriptionStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#A9B1D6"))
)

type CommandInfo struct {
	Name        string
	Description string
	Usage       string
}

var commands = []CommandInfo{
	{
		Name:        "scan",
		Description: "Use the scan commmand to scan your codebase for security vulnerabilities",
		Usage:       "armur scan [flags]",
	},
	{
		Name:        "report",
		Description: "Generate and view security reports from a scan",
		Usage:       "armur report [flags]",
	},
	{
		Name:        "status",
		Description: "Check the status of your scans",
		Usage:       "armur status [flags]",
	},
	{
		Name:        "config",
		Description: "Configure Armur CLI settings which includes api_url, redis_url and api_key",
		Usage:       "armur config [flags]",
	},
	{
		Name:        "docker",
		Description: "For ease of use of the cli, easily start the project in docker using the available commands",
		Usage:       "armur docker [flags]",
	},
	{
		Name:        "apiaudit",
		Description: "Use the API directly from Armur to audit your files",
		Usage:       "armur apiaudit [flags]",
	},
}

func ShowInteractiveHelp() {
	var selectedCommand string

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title(titleStyle.Render("Armur CLI Commands")).
				Description("Select a command to view its details").
				Options(
					huh.NewOption("Scan", "scan"),
					huh.NewOption("Report", "report"),
					huh.NewOption("Status", "status"),
					huh.NewOption("Config", "config"),
					huh.NewOption("Docker", "docker"),
					huh.NewOption("API Audit", "apiaudit"),
				).
				Value(&selectedCommand),
		),
	)

	err := form.Run()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	// Find the selected command info
	var selectedInfo CommandInfo
	for _, cmd := range commands {
		if cmd.Name == selectedCommand {
			selectedInfo = cmd
			break
		}
	}

	// Display command details
	fmt.Printf("\n%s\n", titleStyle.Render("Command Details"))
	fmt.Printf("%s: %s\n", commandStyle.Render("Command"), selectedInfo.Name)
	fmt.Printf("%s: %s\n", descriptionStyle.Render("Description"), selectedInfo.Description)
	fmt.Printf("%s: %s\n", commandStyle.Render("Usage"), selectedInfo.Usage)
	fmt.Println("\nPress Enter to exit...")
	fmt.Scanln()
}

func init() {
	rootCmd.Run = func(cmd *cobra.Command, args []string) {
		ShowInteractiveHelp()
	}
}
