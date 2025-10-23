package cmd

import (
	"armur-cli/internal/api"
	"armur-cli/internal/config"
	"fmt"
	"os"

	"github.com/charmbracelet/huh"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Get the status of a scan task",
	Long:  `Check the status of a specific scan task by providing its ID.`,
	Run: func(cmd *cobra.Command, args []string) {
		var taskID string

		form := huh.NewForm(
			huh.NewGroup(
				huh.NewInput().Title("Task ID").Value(&taskID),
			),
		)

		if err := form.Run(); err != nil {
			fmt.Println("Prompt canceled.")
			return
		}

		cfg, err := config.LoadConfig()
		if err != nil {
			color.Red("Error loading configuration: %v", err)
			os.Exit(1)
		}

		apiClient := api.NewClient(cfg.API.URL)
		status, result, err := apiClient.GetTaskStatus(taskID)
		if err != nil {
			color.Red("Error getting task status: %v", err)
			os.Exit(1)
		}

		fmt.Println(color.CyanString("Task Status:"))
		fmt.Printf("  Task ID: %s\n", taskID)
		fmt.Printf("  Status: %s\n", status)

		switch status {
		case "success":
			fmt.Println(color.GreenString("  Scan Results:"))
			for category, details := range result {
				fmt.Printf("    Category: %s\n", category)
				if issues, ok := details.([]interface{}); ok {
					for _, issue := range issues {
						if issueMap, ok := issue.(map[string]interface{}); ok {
							for key, value := range issueMap {
								fmt.Printf("      %s: %v\n", key, value)
							}
							fmt.Println("      ---")
						}
					}
				} else {
					fmt.Printf("      Details: %v\n", details)
				}
			}
		case "failed":
			color.Red("  Scan failed.")
		default:
			fmt.Printf("  Status: %s\n", status)
		}
	},
}

func init() {
	rootCmd.AddCommand(statusCmd)
}
