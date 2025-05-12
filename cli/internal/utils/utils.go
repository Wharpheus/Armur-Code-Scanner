package utils

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"

	"github.com/fatih/color"
	"github.com/olekukonko/tablewriter"
)

// Category descriptions
var categoryDescriptions = map[string]string{
	"security_issues":    "Potential security vulnerabilities.",
	"dead_code":          "Unused or unreachable code.",
	"duplicate_code":     "Duplicated code blocks.",
	"secret_detection":   "Possible hardcoded secrets.",
	"infra_security":     "Infrastructure-related security issues.",
	"sca":                "Software Composition Analysis: Vulnerabilities in dependencies.",
	"complex_functions":  "Functions with high cyclomatic complexity.",
	"docstring_absent":   "Missing docstrings in functions, classes, or modules.",
	"antipatterns_bugs":  "Code patterns that are considered bad practices or likely to cause bugs.",
}

// PrintResultsText prints the scan results in a tabular format.
func PrintResultsText(category string, issues []interface{}) {
	description, ok := categoryDescriptions[category]
	if !ok {
		description = "Unknown category"
	}

	color.New(color.FgCyan, color.Bold).Printf("Category: %s\n", category)
	fmt.Printf("Description: %s\n\n", description)

	if len(issues) == 0 {
		color.Yellow("No issues found in this category.\n")
		return
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"File", "Line", "Message", "Severity"})
	table.SetHeaderColor(
		tablewriter.Colors{tablewriter.Bold, tablewriter.FgHiCyanColor},
		tablewriter.Colors{tablewriter.Bold, tablewriter.FgHiCyanColor},
		tablewriter.Colors{tablewriter.Bold, tablewriter.FgHiCyanColor},
		tablewriter.Colors{tablewriter.Bold, tablewriter.FgHiCyanColor},
	)
	table.SetColumnColor(
		tablewriter.Colors{tablewriter.FgHiWhiteColor},
		tablewriter.Colors{tablewriter.FgHiWhiteColor},
		tablewriter.Colors{tablewriter.FgHiWhiteColor},
		tablewriter.Colors{}, // Severity will be colored based on its value
	)

	table.SetAutoWrapText(true) // Enable auto wrap for long text
	table.SetAutoFormatHeaders(true)
	table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetCenterSeparator("")
	table.SetColumnSeparator("")
	table.SetRowSeparator("-")

	// Use reflection to check if SetBorderStyle method exists (for backward compatibility)
	tableType := reflect.TypeOf(table)
	if _, ok := tableType.MethodByName("SetBorderStyle"); ok {
		// Use reflection to call SetBorderStyle dynamically
		setBorderStyle := reflect.ValueOf(table).MethodByName("SetBorderStyle")
		borderStyle := reflect.ValueOf(tablewriter.Border{Left: true, Top: false, Right: true, Bottom: false})
		setBorderStyle.Call([]reflect.Value{borderStyle})
	}

	table.SetTablePadding("  ") // Add padding
	table.SetNoWhiteSpace(true)

	for _, issue := range issues {
		issueMap, ok := issue.(map[string]interface{})
		if !ok {
			continue
		}

		file, line, message, severity := extractIssueDetails(category, issueMap)

		// Colorize severity
		switch severity {
		case "HIGH":
			severity = color.RedString(severity)
		case "MEDIUM":
			severity = color.YellowString(severity)
		case "LOW":
			severity = color.GreenString(severity)
		case "INFO":
			severity = color.CyanString(severity)
		}

		table.Append([]string{file, line, message, severity})

	}

	if table.NumLines() > 0 { // Check if the table has any rows
		table.Render()
		fmt.Println()
	} else {
		color.Yellow("No issues found in this category.\n")
	}
}

// extractIssueDetails extracts relevant information from an issue map based on the category.
func extractIssueDetails(category string, issueMap map[string]interface{}) (string, string, string, string) {
	var file, line, message, severity string

	file = "<nil>"
	line = "<nil>"
	message = "<nil>"
	severity = "<nil>"

	if val, ok := issueMap["path"]; ok && val != nil {
		file = fmt.Sprintf("%v", val)
	}
	if val, ok := issueMap["line"]; ok && val != nil {
		line = fmt.Sprintf("%v", val)
	}
	if val, ok := issueMap["message"]; ok && val != nil {
		message = fmt.Sprintf("%v", val)
	}
	if val, ok := issueMap["severity"]; ok && val != nil {
		severity = fmt.Sprintf("%v", val)
	}

	switch category {
	case "security_issues":
		// Handle cases where 'cwe' field is present but not a string
		if cweVal, ok := issueMap["cwe"]; ok {
			if cweStr, ok := cweVal.(string); ok {
				message = fmt.Sprintf("CWE: %s - %s", cweStr, message)
			} else if cweArr, ok := cweVal.([]interface{}); ok && len(cweArr) > 0 {
				if cweStr, ok := cweArr[0].(string); ok {
					message = fmt.Sprintf("CWE: %s - %s", cweStr, message)
				}
			}
		}
	case "dead_code":
		if issueMap["file"] != nil {
			file = fmt.Sprintf("%v", issueMap["file"])
		}
		severity = "INFO"
	case "duplicate_code":
		if issueMap["first_file"] != nil {
			file = fmt.Sprintf("%v", issueMap["first_file"])
		}
		if issueMap["second_file"] != nil && issueMap["start_in_second_file"] != nil && issueMap["end_in_second_file"] != nil {
			message = fmt.Sprintf("Duplication found with %v (lines %v-%v)", issueMap["second_file"], issueMap["start_in_second_file"], issueMap["end_in_second_file"])
		}
		severity = "INFO"
	case "secret_detection":
		// No changes needed
	case "infra_security":
		line = "N/A"
	case "sca":
		line = "N/A"
		message = fmt.Sprintf("Vulnerability in package: %v, Version: %v, Details: %v", issueMap["package"], issueMap["version"], issueMap["message"])
	case "complex_functions":
		if complexity, ok := issueMap["complexity"].(float64); ok {
			message = fmt.Sprintf("High cyclomatic complexity: %v", int(complexity))
		} else {
			message = "High cyclomatic complexity"
		}
		severity = "COMPLEX"

	case "docstring_absent":
		message = fmt.Sprintf("Missing Docstring: %v", issueMap["message"])
		severity = "INFO"

	case "antipatterns_bugs":
		severity = "INFO"

	default:
		severity = "UNKNOWN"
	}

	return file, line, message, severity
}

// PrintResultsJSON prints the scan results in JSON format.
func PrintResultsJSON(results map[string]interface{}) {
	jsonResult, err := json.MarshalIndent(results, "", "  ")
	if err != nil {
		color.Red("Error converting results to JSON: %v", err)
		return
	}

	fmt.Println(string(jsonResult))
}