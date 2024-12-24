package cmd

import (
    "fmt"
    "os"
    "os/exec"

    "github.com/fatih/color"
    "github.com/spf13/cobra"
)

var dockerCmd = &cobra.Command{
    Use:   "docker",
    Short: "Manage Docker Compose",
    Long:  `Manage the Docker Compose setup for the Armur Code Scanner service.`,
}

var dockerUpCmd = &cobra.Command{
    Use:   "up",
    Short: "Start Docker Compose services",
    Long:  `Starts the Docker Compose services for Armur Code Scanner.`,
    Run: func(cmd *cobra.Command, args []string) {
        fmt.Println(color.CyanString("Starting Docker Compose services..."))
        if err := runDockerComposeCommand("up", "-d"); err != nil {
            color.Red("Error starting Docker Compose services: %v", err)
            os.Exit(1)
        }
        fmt.Println(color.GreenString("Docker Compose services started successfully."))
    },
}

var dockerDownCmd = &cobra.Command{
    Use:   "down",
    Short: "Stop Docker Compose services",
    Long:  `Stops the Docker Compose services for Armur Code Scanner.`,
    Run: func(cmd *cobra.Command, args []string) {
        fmt.Println(color.CyanString("Stopping Docker Compose services..."))
        if err := runDockerComposeCommand("down"); err != nil {
            color.Red("Error stopping Docker Compose services: %v", err)
            os.Exit(1)
        }
        fmt.Println(color.GreenString("Docker Compose services stopped successfully."))
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
    dockerCmd.AddCommand(dockerUpCmd)
    dockerCmd.AddCommand(dockerDownCmd)
}