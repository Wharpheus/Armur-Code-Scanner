package internal

import (
	"bytes"
	"encoding/json"
	"log"
	"os/exec"
	"strings"

	utils "armur-codescanner/pkg"
)

func RunRadon(directory string) map[string]interface{} {
	log.Println("Running Radon...")
	radonResults := RunRadonOnRepo(directory)
	categorizedResults := CategorizeRadonResults(radonResults, directory)
	newcattu := utils.ConvertCategorizedResults(categorizedResults)
	return newcattu
}

func RunRadonOnRepo(directory string) string {
	cmd := exec.Command("radon", "cc", "-j", directory)
	var out, stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		log.Printf("Error running Radon: %v\n", err)
		log.Printf("stderr: %s\n", stderr.String())
		return ""
	}

	return out.String()
}

func CategorizeRadonResults(results string, directory string) map[string][]interface{} {
	// Initialize categorized results
	categorizedResults := utils.InitCategorizedResults()

	if results != "" {
		// Parse the JSON results from Radon
		var parsedResults map[string][]map[string]interface{}
		err := json.Unmarshal([]byte(results), &parsedResults)
		if err != nil {
			log.Printf("Error parsing Radon results: %v\n", err)
			return categorizedResults
		}

		// Process each file and its issues
		for filePath, issues := range parsedResults {
			relativePath := strings.Replace(filePath, directory, "", 1)
			for _, issue := range issues {
				// Add the path to the issue and append it directly to COMPLEX_FUNCTIONS
				issue["path"] = relativePath
				categorizedResults[utils.COMPLEX_FUNCTIONS] = append(categorizedResults[utils.COMPLEX_FUNCTIONS], issue)
			}
		}
	}

	return categorizedResults
}
