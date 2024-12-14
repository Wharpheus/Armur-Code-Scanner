package internal

import (
	utils "armur-codescanner/pkg"
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

type CWEData struct {
	CWE                 string            `json:"cwe"`
	GoodPracticeExample map[string]string `json:"good_practice_example"`
	BadPracticeExample  map[string]string `json:"bad_practice_example"`
}

func GetPracticesFromJSON(data []CWEData, language string, cwe string) map[string]string {
	for _, item := range data {
		if item.CWE == cwe {
			goodPractice, goodFound := item.GoodPracticeExample[language]
			badPractice, badFound := item.BadPracticeExample[language]

			if goodFound && badFound {
				return map[string]string{
					"good_practice": goodPractice,
					"bad_practice":  badPractice,
				}
			}
		}
	}
	return map[string]string{}
}

func RunGosec(directory string) map[string]interface{} {
	log.Println("Running Gosec")
	results, _ := RunGosecOnRepo(directory)

	categorizedResults := CategorizeGosecResults(results, directory)
	newcategorizedResults := utils.ConvertCategorizedResults(categorizedResults)
	return newcategorizedResults
}

func RunGosecOnRepo(directory string) (string, error) {
	cmd := exec.Command("gosec", "-fmt=json", "./...")
	cmd.Dir = directory

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	cmd.Run()
	fmt.Println(stdout.String())
	return stdout.String(), nil
}

func CategorizeGosecResults(results string, directory string) map[string][]interface{} {
	categorizedResults := utils.InitCategorizedResults()

	// Read the CWE JSON file
	cweDataFile, err := os.ReadFile("pkg/common/cwd.json")
	if err != nil {
		log.Printf("Error reading CWE JSON: %v", err)
		return categorizedResults
	}

	var cweData []CWEData
	if err := json.Unmarshal(cweDataFile, &cweData); err != nil {
		log.Printf("Error parsing CWE JSON: %v", err)
		return categorizedResults
	}

	if results != "" {
		var parsedResults map[string]interface{}
		if err := json.Unmarshal([]byte(results), &parsedResults); err != nil {
			log.Printf("Error parsing Gosec results: %v", err)
			return categorizedResults
		}

		// Process Issues
		if issues, ok := parsedResults["Issues"].([]interface{}); ok {
			for _, issue := range issues {
				issueMap, mapOk := issue.(map[string]interface{})
				if !mapOk {
					log.Printf("Invalid issue format: %v", issue)
					continue
				}

				// Extract CWE ID
				cweID := "unknown"
				if cweObj, ok := issueMap["cwe"].(map[string]interface{}); ok {
					if id, ok := cweObj["id"].(string); ok {
						cweID = fmt.Sprintf("CWE-%s", id)
					}
				}

				formattedIssue := map[string]interface{}{
					"path":       strings.Replace(issueMap["file"].(string), directory, "", 1),
					"line":       issueMap["line"],
					"column":     issueMap["column"],
					"message":    issueMap["details"],
					"severity":   strings.ToUpper(issueMap["severity"].(string)),
					"check_id":   issueMap["rule_id"],
					"confidence": strings.ToUpper(issueMap["confidence"].(string)),
					"likelihood": "NA",
					"cwe":        []string{cweID},
					"owasp":      []string{},
					"endLine":    issueMap["line"],
				}

				// Update practices
				practices := GetPracticesFromJSON(cweData, "go", cweID)
				for key, value := range practices {
					formattedIssue[key] = value
				}

				// Append formatted issue to SECURITY_ISSUES
				categorizedResults[SECURITY_ISSUES] = append(categorizedResults[SECURITY_ISSUES], formattedIssue)
			}
		}

		// Process Golang Errors
		if golangErrors, ok := parsedResults["Golang errors"].(map[string]interface{}); ok {
			for path, errors := range golangErrors {
				errorList, listOk := errors.([]interface{})
				if !listOk {
					log.Printf("Invalid Golang errors format for path %s: %v", path, errors)
					continue
				}

				for _, err := range errorList {
					errorMap, mapOk := err.(map[string]interface{})
					if !mapOk {
						log.Printf("Invalid error format: %v", err)
						continue
					}

					antipatternsBug := map[string]interface{}{
						"path":    strings.Replace(path, directory, "", 1),
						"line":    fmt.Sprintf("%v", errorMap["line"]),
						"column":  fmt.Sprintf("%v", errorMap["column"]),
						"message": fmt.Sprintf("%v", errorMap["error"]),
					}

					// Append antipatterns bug to ANTIPATTERNS_BUGS
					categorizedResults[ANTIPATTERNS_BUGS] = append(categorizedResults[ANTIPATTERNS_BUGS], antipatternsBug)
				}
			}
		}
	}

	return categorizedResults
}
