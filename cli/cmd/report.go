package cmd

import (
	"armur-cli/internal/api"
	"armur-cli/internal/config"
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var reportCmd = &cobra.Command{
	Use:   "report [taskID]",
	Short: "Generate a report for a completed scan task",
	Long:  `Generate an OWASP or SANS report for a completed scan task by providing its ID.`,
	Args:  cobra.ExactArgs(1), // Expect exactly one argument (task ID)
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.LoadConfig()
		if err != nil {
			color.Red("Error loading configuration: %v", err)
			os.Exit(1)
		}

		apiClient := api.NewClient(cfg.API.URL)
		taskID := args[0]
		reportType, _ := cmd.Flags().GetString("type")

		var report interface{}
		var reportErr error

		switch reportType {
		case "owasp":
			report, reportErr = apiClient.GetOwaspReport(taskID)
		case "sans":
			report, reportErr = apiClient.GetSansReport(taskID)
		default:
			color.Red("Invalid report type. Use 'owasp' or 'sans'.")
			os.Exit(1)
		}

		if reportErr != nil {
			color.Red("Error generating report: %v", reportErr)
			os.Exit(1)
		}

		// Display the report (you might want to format this more nicely)
		fmt.Println(color.CyanString("%s Report for Task ID %s:", strings.ToUpper(reportType), taskID))
		fmt.Println(report)
	},
}

func init() {
	rootCmd.AddCommand(reportCmd)
	reportCmd.Flags().StringP("type", "t", "owasp", "Type of report to generate (owasp, sans)")
}