package internal

import (
	"bytes"
	"encoding/json"
	"log"
	"os/exec"
	"strings"

	"armur-codescanner/pkg" // Replace with your actual utils package import
)

func RunPylint(directory string) map[string]interface{} {
	log.Println("Running Pylint...")
	pylintOutput := RunPylintOnRepo(directory)
	categorizedResults := CategorizePylintResults(pylintOutput, directory)
	newcatresult := utils.ConvertCategorizedResults(categorizedResults)
	return newcatresult
}

func RunPylintOnRepo(directory string) string {
	cmd := exec.Command("pylint", "--recursive=y", directory, "-f", "json")
	var out, stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	cmd.Run()

	return out.String()
}

func CategorizePylintResults(results string, directory string) map[string][]interface{} {
	// Initialize categorized results
	categorizedResults := utils.InitCategorizedResults()

	if results != "" {
		// Parse the JSON results from Pylint
		var parsedResults []map[string]interface{}
		err := json.Unmarshal([]byte(results), &parsedResults)
		if err != nil {
			log.Printf("Error parsing Pylint results: %v\n", err)
			return categorizedResults
		}

		// Normalize file paths and convert to []interface{}
		directory = strings.Replace(directory, "/armur/", "", -1)
		for _, result := range parsedResults {
			if path, ok := result["path"].(string); ok {
				result["path"] = strings.Replace(path, directory, "", -1)
			}
			// Append each map[string]interface{} as interface{}
			categorizedResults[ANTIPATTERNS_BUGS] = append(categorizedResults[ANTIPATTERNS_BUGS], result)
		}
	}

	return categorizedResults
}
