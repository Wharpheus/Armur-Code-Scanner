package cmd

import (
	"armur-cli/internal/api"
	"armur-cli/internal/config"
	"armur-cli/internal/utils"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/briandowns/spinner"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var scanCmd = &cobra.Command{
	Use:   "scan [target]",
	Short: "Scan a repository or file",
	Long: `Scan a Git repository (by providing the URL) or a local file/directory (by providing the path)
for security vulnerabilities using the Armur Code Scanner service.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.LoadConfig()
		if err != nil {
			color.Red("Error loading configuration: %v", err)
			os.Exit(1)
		}

		apiClient := api.NewClient(cfg.API.URL)
		target := args[0]
		language, _ := cmd.Flags().GetString("language")
		isAdvanced, _ := cmd.Flags().GetBool("advanced")
		outputFormat, _ := cmd.Flags().GetString("output")

		if (strings.HasPrefix(target, "http://") || strings.HasPrefix(target, "https://")) && language == "" {
			color.Red("Error: Language must be specified when scanning a repository. Use --language or -l flag.")
			os.Exit(1)
		}

		var taskID string
		var scanErr error

		if strings.HasPrefix(target, "http://") || strings.HasPrefix(target, "https://") {
			fmt.Println(color.CyanString("Initiating scan for repository: %s", target))
			taskID, scanErr = apiClient.ScanRepository(target, language, isAdvanced)
		} else {
			fmt.Println(color.CyanString("Initiating scan for local target: %s", target))
			taskID, scanErr = apiClient.ScanFile(target, isAdvanced)
		}

		if scanErr != nil {
			color.Red("Error initiating scan: %v", scanErr)
			os.Exit(1)
		}

		fmt.Println(color.GreenString("Scan initiated successfully. Task ID: %s", taskID))

		s := spinner.New(spinner.CharSets[14], 100*time.Millisecond)
		s.Suffix = " Waiting for scan to complete..."
		s.Start()

		var scanResult map[string]interface{}
		var status string
		for {
			status, scanResult, err = apiClient.GetTaskStatus(taskID)
			if err != nil {
				s.Stop()
				color.Red("Error getting task status: %v", err)
				os.Exit(1)
			}

			if status == "success" {
				s.Stop()
				break
			} else if status == "failed" {
				s.Stop()
				color.Red("Scan failed.")
				os.Exit(1)
			}
			time.Sleep(2 * time.Second)
		}

		if outputFormat == "json" {
			utils.PrintResultsJSON(scanResult)
		} else {
			printFormattedResults(scanResult)
		}
	},
}

func printFormattedResults(results map[string]interface{}) {
	categories := []string{
		"security_issues",
		"antipatterns_bugs",
		"complex_functions",
		"dockstring_absent",
	}

	for _, category := range categories {
		if categoryData, ok := results[category]; ok {
			fmt.Printf("\n%s\n", color.CyanString("Category: %s", category))
			printCategoryData(category, categoryData)
		}
	}
	fmt.Println(color.GreenString("\nScan completed. Results displayed."))
}

func printCategoryData(category string, data interface{}) {
	descriptions := map[string]string{
		"security_issues":    "Security vulnerabilities and potential threats",
		"antipatterns_bugs":  "Code patterns that are considered bad practices or likely to cause bugs",
		"complex_functions":  "Functions with high cyclomatic complexity",
		"dockstring_absent":  "Missing documentation strings",
	}

	fmt.Printf("Description: %s\n\n", descriptions[category])

	switch category {
	case "security_issues":
		printSecurityIssues(data)
	case "antipatterns_bugs":
		printAntipatternsBugs(data)
	case "complex_functions":
		printComplexFunctions(data)
	case "dockstring_absent":
		printDocstringIssues(data)
	}
}

func printSecurityIssues(data interface{}) {
	issues, ok := data.([]interface{})
	if !ok || len(issues) == 0 {
		color.Yellow("No security issues found.\n")
		return
	}

	fmt.Printf("%-20s %-10s %-40s %-15s %-20s\n", "FILE", "LINE", "MESSAGE", "SEVERITY", "CWE")
	fmt.Println(strings.Repeat("-", 105))

	for _, issue := range issues {
		if issueMap, ok := issue.(map[string]interface{}); ok {
			files, ok := issueMap["files"].([]interface{})
			if !ok {
				continue
			}

			cwe := getStringValue(issueMap, "cwe", "N/A")
			for _, file := range files {
				if fileMap, ok := file.(map[string]interface{}); ok {
					issues, ok := fileMap["issues"].([]interface{})
					if !ok {
						continue
					}

					for _, issueDetail := range issues {
						if detail, ok := issueDetail.(map[string]interface{}); ok {
							fmt.Printf("%-20s %-10s %-40s %-15s %-20s\n",
								truncateString(getStringValue(detail, "path", "-"), 20),
								getStringValue(detail, "line", "-"),
								truncateString(getStringValue(detail, "message", "-"), 40),
								getSeverityColored(getStringValue(detail, "severity", "INFO")),
								cwe,
							)
						}
					}
				}
			}
		}
	}
	fmt.Println()
}

func printAntipatternsBugs(data interface{}) {
	issues, ok := data.([]interface{})
	if !ok || len(issues) == 0 {
		color.Yellow("No antipatterns or bugs found.\n")
		return
	}

	fmt.Printf("%-30s %-10s %-50s\n", "FILE", "LINE", "MESSAGE")
	fmt.Println(strings.Repeat("-", 90))

	for _, issue := range issues {
		if issueMap, ok := issue.(map[string]interface{}); ok {
			files, ok := issueMap["files"].([]interface{})
			if !ok {
				continue
			}

			for _, file := range files {
				if fileMap, ok := file.(map[string]interface{}); ok {
					path := getStringValue(fileMap, "path", "-")
					issues, ok := fileMap["issues"].([]interface{})
					if !ok {
						continue
					}

					for _, issueDetail := range issues {
						if detail, ok := issueDetail.(map[string]interface{}); ok {
							fmt.Printf("%-30s %-10s %-50s\n",
								truncateString(path, 30),
								getStringValue(detail, "line", "-"),
								truncateString(getStringValue(detail, "message", "-"), 50),
							)
						}
					}
				}
			}
		}
	}
	fmt.Println()
}

func printComplexFunctions(data interface{}) {
	issues, ok := data.([]interface{})
	if !ok || len(issues) == 0 {
		color.Yellow("No complex functions found.\n")
		return
	}

	fmt.Printf("%-40s %-10s %-20s\n", "FILE", "LINE", "COMPLEXITY")
	fmt.Println(strings.Repeat("-", 70))

	for _, issue := range issues {
		if issueMap, ok := issue.(map[string]interface{}); ok {
			path := getStringValue(issueMap, "path", "-")
			issues, ok := issueMap["issues"].([]interface{})
			if !ok {
				continue
			}

			for _, issueDetail := range issues {
				if detail, ok := issueDetail.(map[string]interface{}); ok {
					fmt.Printf("%-40s %-10s %-20s\n",
						truncateString(path, 40),
						getStringValue(detail, "line", "-"),
						getStringValue(detail, "complexity", "-"),
					)
				}
			}
		}
	}
	fmt.Println()
}

func printDocstringIssues(data interface{}) {
	issues, ok := data.([]interface{})
	if !ok || len(issues) == 0 {
		color.Yellow("No docstring issues found.\n")
		return
	}

	fmt.Printf("%-40s %-10s %-40s\n", "FILE", "LINE", "MESSAGE")
	fmt.Println(strings.Repeat("-", 90))

	for _, issue := range issues {
		if issueMap, ok := issue.(map[string]interface{}); ok {
			path := getStringValue(issueMap, "path", "-")
			fmt.Printf("%-40s %-10s %-40s\n",
				truncateString(path, 40),
				getStringValue(issueMap, "line", "-"),
				truncateString(getStringValue(issueMap, "message", "Missing docstring"), 40),
			)
		}
	}
	fmt.Println()
}

func getStringValue(m map[string]interface{}, key, defaultValue string) string {
	if val, ok := m[key]; ok && val != nil {
		switch v := val.(type) {
		case string:
			return v
		case float64:
			return fmt.Sprintf("%d", int(v))
		case int:
			return fmt.Sprintf("%d", v)
		default:
			return defaultValue
		}
	}
	return defaultValue
}

func getSeverityColored(severity string) string {
	switch strings.ToUpper(severity) {
	case "HIGH":
		return color.RedString(severity)
	case "MEDIUM":
		return color.YellowString(severity)
	case "LOW":
		return color.GreenString(severity)
	case "INFO":
		return color.CyanString(severity)
	default:
		return severity
	}
}

func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}

func init() {
	rootCmd.AddCommand(scanCmd)
	scanCmd.Flags().StringP("language", "l", "", "Specify the programming language")
	scanCmd.Flags().BoolP("simple", "s", true, "Perform a simple scan")
	scanCmd.Flags().BoolP("advanced", "a", false, "Perform an advanced scan")
	scanCmd.Flags().StringP("output", "o", "text", "Output format (text, json)")
}