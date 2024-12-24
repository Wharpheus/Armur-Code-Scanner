package cmd

import (
	"armur-cli/internal/api"
	"armur-cli/internal/config"
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var statusCmd = &cobra.Command{
	Use:   "status [taskID]",
	Short: "Get the status of a scan task",
	Long:  `Check the status of a specific scan task by providing its ID.`,
	Args:  cobra.ExactArgs(1), // Expect exactly one argument (task ID)
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.LoadConfig()
		if err != nil {
			color.Red("Error loading configuration: %v", err)
			os.Exit(1)
		}

		apiClient := api.NewClient(cfg.API.URL)
		taskID := args[0]

		status, result, err := apiClient.GetTaskStatus(taskID)
		if err != nil {
			color.Red("Error getting task status: %v", err)
			os.Exit(1)
		}

		fmt.Println(color.CyanString("Task Status:"))
		fmt.Printf("  Task ID: %s\n", taskID)
		fmt.Printf("  Status: %s\n", status)

		if status == "success" {
			fmt.Println(color.GreenString("  Scan Results:"))
			// Iterate through the results and print them
			for category, details := range result {
				fmt.Printf("    Category: %s\n", category)
				// Assuming the details are in a slice, iterate through them
				if issues, ok := details.([]interface{}); ok {
					for _, issue := range issues {
						// Assuming each issue is a map, print its details
						if issueMap, ok := issue.(map[string]interface{}); ok {
							for key, value := range issueMap {
								fmt.Printf("      %s: %v\n", key, value)
							}
							fmt.Println("      ---") // Separator between issues
						}
					}
				} else {
					// Handle the case where details are not a slice
					fmt.Printf("      Details: %v\n", details)
				}
			}
		} else if status == "failed" {
			color.Red("  Scan failed.")
		}
	},
}

func init() {
	rootCmd.AddCommand(statusCmd)
}