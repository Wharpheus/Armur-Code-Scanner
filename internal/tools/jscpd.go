package internal

import (
	utils "armur-codescanner/pkg"
	"encoding/json"
	"log"
	"os"
	"os/exec"
	"strings"
)

const (
	DUPLICATE_CODE_LINE_THRESHOLD = 10
)

func RunJSCPD(directory string) map[string]interface{} {
	log.Println("Running jscpd")
	duplicates, err := RunJSCPDOnRepo(directory)
	if err != nil {
		log.Printf("Error while running jscpd: %v", err)
		return nil
	}
	results := CategorizeJSCPDResults(duplicates, directory)
	newcatresult := utils.ConvertCategorizedResults(results)
	return newcatresult
}

func RunJSCPDOnRepo(directory string) ([]map[string]interface{}, error) {
	cmd := exec.Command("jscpd", directory, "-r", "json", "-o", ".")
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("Failed to run jscpd: %v, Output: %s", err, string(output))
		return nil, err
	}

	outputFile := "jscpd-report.json"
	if _, err := os.Stat(outputFile); err == nil {
		file, err := os.Open(outputFile)
		if err != nil {
			log.Printf("Failed to open report file: %v", err)
			return nil, err
		}
		defer file.Close()

		var report map[string]interface{}
		if err := json.NewDecoder(file).Decode(&report); err != nil {
			log.Printf("Failed to parse JSON report: %v", err)
			return nil, err
		}

		duplicates, ok := report["duplicates"].([]interface{})
		if !ok {
			log.Println("No duplicates found in the report")
			return nil, nil
		}

		var result []map[string]interface{}
		for _, d := range duplicates {
			if dup, ok := d.(map[string]interface{}); ok {
				result = append(result, dup)
			}
		}
		return result, nil
	}

	log.Println("No report generated or file not found")
	return nil, nil
}

func CategorizeJSCPDResults(results []map[string]interface{}, directory string) map[string][]interface{} {
	categorizedResults := utils.InitAdvancedCategorizedResults()

	for _, duplicate := range results {
		firstFile, _ := duplicate["firstFile"].(map[string]interface{})
		secondFile, _ := duplicate["secondFile"].(map[string]interface{})

		firstFileName, _ := firstFile["name"].(string)
		secondFileName, _ := secondFile["name"].(string)
		startInFirstFile, _ := firstFile["start"].(float64)
		endInFirstFile, _ := firstFile["end"].(float64)
		startInSecondFile, _ := secondFile["start"].(float64)
		endInSecondFile, _ := secondFile["end"].(float64)

		cleanDirectory := strings.Replace(directory, "/armur/", "", -1)

		categorizedResults[utils.DUPLICATE_CODE] = append(categorizedResults[utils.DUPLICATE_CODE], map[string]interface{}{
			"first_file":           strings.Replace(firstFileName, cleanDirectory, "", -1),
			"second_file":          strings.Replace(secondFileName, cleanDirectory, "", -1),
			"start_in_first_file":  startInFirstFile,
			"end_in_first_file":    endInFirstFile,
			"start_in_second_file": startInSecondFile,
			"end_in_second_file":   endInSecondFile,
			"severity":             determineSeverity(endInFirstFile - startInFirstFile),
		})
	}

	return categorizedResults
}

func determineSeverity(lines float64) string {
	if lines < DUPLICATE_CODE_LINE_THRESHOLD {
		return "Minor"
	}
	return "Major"
}
