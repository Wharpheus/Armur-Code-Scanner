package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/charmbracelet/huh"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var dockerCmd = &cobra.Command{
	Use:   "docker",
	Short: "Manage Docker Compose",
	Long:  `Manage the Docker Compose setup for the Armur Code Scanner service.`,
	Run: func(cmd *cobra.Command, args []string) {
		var action string
		form := huh.NewForm(
			huh.NewGroup(
				huh.NewSelect[string]().
					Title("What do you want to do?").
					Options(
						huh.NewOption("Start Docker Services", "up"),
						huh.NewOption("Stop Docker Services", "down"),
					).
					Value(&action),
			),
		)

		if err := form.Run(); err != nil {
			fmt.Println("Prompt canceled.")
			return
		}

		switch action {
		case "up":
			fmt.Println(color.CyanString("Starting Docker Compose services..."))
			if err := runDockerComposeCommand("up", "-d"); err != nil {
				color.Red("Error starting Docker Compose services: %v", err)
				os.Exit(1)
			}
			fmt.Println(color.GreenString("Docker Compose services started."))
		case "down":
			fmt.Println(color.CyanString("Stopping Docker Compose services..."))
			if err := runDockerComposeCommand("down"); err != nil {
				color.Red("Error stopping Docker Compose services: %v", err)
				os.Exit(1)
			}
			fmt.Println(color.GreenString("Docker Compose services stopped."))
		}
	},
}

func runDockerComposeCommand(args ...string) error {
	cmd := exec.Command("docker-compose", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func init() {
	rootCmd.AddCommand(dockerCmd)
}
