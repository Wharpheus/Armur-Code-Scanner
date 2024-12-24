package internal

import (
	utils "armur-codescanner/pkg"
	"fmt"
	"log"
	"os/exec"
	"strings"
)

func RunGoDeadcode(directory string) map[string]interface{} {
	log.Println("Running deadcode")
	results, err := RunDecodeOnRepO(directory)
	if err != nil {
		log.Printf("Error while running deadcode: %v", err)
		return nil
	}
	catresult := categorizeDeadCodeResults(results, directory)
	newcatresult := utils.ConvertCategorizedResults(catresult)
	return newcatresult
}

func RunDecodeOnRepO(directory string) (string, error) {
	cmd := exec.Command("deadcode", directory)
	output, err := cmd.CombinedOutput()
	fmt.Println("Go deadcode Error: ", err)
	return strings.TrimSpace(string(output)), nil
}

func categorizeDeadCodeResults(results string, directory string) map[string][]interface{} {
	categorizedResults := utils.InitAdvancedCategorizedResults()

	// If results are empty or only contain errors, return clean empty results
	if results == "" || !strings.Contains(results, ":") { // Adjust to detect valid output format
		return categorizedResults
	}

	// Process results line by line
	lines := strings.Split(results, "\n")
	for _, line := range lines {
		if strings.TrimSpace(line) == "" { // Skip empty lines
			continue
		}
		formattedIssue := formatDeadCodeIssue(line, directory)
		if len(formattedIssue) > 0 {
			categorizedResults[utils.DEAD_CODE] = append(categorizedResults[utils.DEAD_CODE], formattedIssue)
		}
	}

	return categorizedResults
}

func formatDeadCodeIssue(issue string, directory string) map[string]interface{} {
	parts := strings.SplitN(issue, ":", 3)
	if len(parts) < 3 {
		return nil // Return nil if the issue format is invalid
	}
	return map[string]interface{}{
		"path":    strings.TrimPrefix(parts[0], directory),
		"line":    strings.TrimSpace(parts[1]),
		"message": strings.TrimSpace(parts[2]),
	}
}
