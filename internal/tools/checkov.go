package internal

import (
	utils "armur-codescanner/pkg"
	"encoding/json"
	"fmt"
	"log"
	"os/exec"
	"strings"
)

func RunCheckov(directory string) map[string]interface{} {
	log.Println("Running Checkov")
	checkovResults := runCheckovOnRepo(directory)
	categorizedResults := categorizeCheckovResults(checkovResults, directory)
	newcatresult := utils.ConvertCategorizedResults(categorizedResults)
	return newcatresult
}

func runCheckovOnRepo(directory string) string {
	cmd := exec.Command("checkov", "-d", directory, "--quiet", "--compact", "-o", "json")
	result, _ := cmd.CombinedOutput()
	return string(result)
}
func categorizeCheckovResults(results string, directory string) map[string][]interface{} {
	categorizedResults := utils.InitAdvancedCategorizedResults()

	if results != "" {
		var checkovResults interface{}
		err := json.Unmarshal([]byte(results), &checkovResults)
		if err != nil {
			log.Printf("Failed to parse Checkov output: %v", err)
			return categorizedResults
		}

		switch checkovResults := checkovResults.(type) {
		case []interface{}:
			for _, result := range checkovResults {
				resultMap, ok := result.(map[string]interface{})
				if !ok {
					log.Printf("Skipping result as it is not a map: %v", result)
					continue
				}
				if resultsMap, ok := resultMap["results"].(map[string]interface{}); ok {
					if failedChecks, ok := resultsMap["failed_checks"].([]interface{}); ok {
						for _, failedCheck := range failedChecks {
							check, ok := failedCheck.(map[string]interface{})
							if !ok {
								log.Printf("Skipping failed check as it is not a map: %v", failedCheck)
								continue
							}
							categorizedResults[utils.INFRA_SECURITY] = append(
								categorizedResults[utils.INFRA_SECURITY],
								formatCheckovIssue(check, directory),
							)
						}
					} else {
						log.Printf("No failed_checks found in results: %v", resultsMap)
					}
				} else {
					log.Printf("No results map found in result: %v", resultMap)
				}
			}
		case map[string]interface{}:
			if resultsMap, ok := checkovResults["results"].(map[string]interface{}); ok {
				if failedChecks, ok := resultsMap["failed_checks"].([]interface{}); ok {
					for _, failedCheck := range failedChecks {
						check, ok := failedCheck.(map[string]interface{})
						if !ok {
							log.Printf("Skipping failed check as it is not a map: %v", failedCheck)
							continue
						}
						categorizedResults[utils.INFRA_SECURITY] = append(
							categorizedResults[utils.INFRA_SECURITY],
							formatCheckovIssue(check, directory),
						)
					}
				} else {
					log.Printf("No failed_checks found in results: %v", resultsMap)
				}
			} else {
				log.Printf("No results map found in checkovResults: %v", checkovResults)
			}
		default:
			log.Printf("Unexpected structure of Checkov results: %v", checkovResults)
		}
	}

	return categorizedResults
}

func formatCheckovIssue(issue map[string]interface{}, directory string) map[string]interface{} {
	// Safely access the fields using type assertions
	filePath, _ := issue["file_path"].(string)
	checkID, _ := issue["check_id"].(string)
	severity, _ := issue["severity"].(string)
	checkName, _ := issue["check_name"].(string)
	fileLineRange, _ := issue["file_line_range"].([]interface{}) // Access file_line_range

	// Convert file_line_range to the format "start:end"
	var fileLineRangeStr string
	if len(fileLineRange) == 2 {
		// Ensure to properly convert the elements to integers
		startLine, ok1 := fileLineRange[0].(float64) // Use float64 to handle JSON numbers
		endLine, ok2 := fileLineRange[1].(float64)
		if ok1 && ok2 {
			// Convert to int and format
			fileLineRangeStr = fmt.Sprintf("%d:%d", int(startLine), int(endLine))
		}
	}

	// Format the issue and return
	return map[string]interface{}{
		"path":            strings.Replace(filePath, directory, "", 1),
		"check_id":        checkID,
		"severity":        severity,
		"message":         checkName,
		"file_line_range": fileLineRangeStr, // Include the formatted file line range
	}
}
