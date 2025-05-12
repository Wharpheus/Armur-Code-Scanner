package cmd

import (
	"armur-cli/internal/api"
	"armur-cli/internal/config"
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var reportCmd = &cobra.Command{
	Use:   "report",
	Short: "Generate a report for a completed scan task",
	Long:  `Generate an OWASP or SANS report for a completed scan task.`,
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.LoadConfig()
		if err != nil {
			color.Red("Error loading configuration: %v", err)
			os.Exit(1)
		}

		var taskID, reportType string
		form := huh.NewForm(
			huh.NewGroup(
				huh.NewInput().Title("Task ID").Value(&taskID),
				huh.NewSelect[string]().
					Title("Select report type").
					Options(
						huh.NewOption("OWASP", "owasp"),
						huh.NewOption("SANS", "sans"),
					).
					Value(&reportType),
			),
		)

		if err := form.Run(); err != nil {
			fmt.Println("Prompt canceled.")
			return
		}

		apiClient := api.NewClient(cfg.API.URL)

		var report interface{}
		var reportErr error

		switch reportType {
		case "owasp":
			report, reportErr = apiClient.GetOwaspReport(taskID)
		case "sans":
			report, reportErr = apiClient.GetSansReport(taskID)
		default:
			color.Red("Invalid report type.")
			os.Exit(1)
		}

		if reportErr != nil {
			color.Red("Error generating report: %v", reportErr)
			os.Exit(1)
		}

		fmt.Println(color.CyanString("%s Report for Task ID %s:", strings.ToUpper(reportType), taskID))
		fmt.Println(report)
	},
}

func init() {
	rootCmd.AddCommand(reportCmd)
}
