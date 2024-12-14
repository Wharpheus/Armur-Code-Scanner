package internal

import (
	utils "armur-codescanner/pkg"
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func RunStaticCheck(directory string) map[string]interface{} {

	defer func() {
		if r := recover(); r != nil {
			log.Printf("Error while running Staticcheck: %v", r)
		}
	}()

	log.Println("Running Staticcheck")
	staticcheckResults, err := RunStaticcheckOnRepo(directory)
	if err != nil {
		log.Printf("Error while running Staticcheck: %v", err)
		return nil
	}
	newcategorisedresult := CategorizeStaticcheckResults(staticcheckResults, directory)
	newcatresult := utils.ConvertCategorizedResults(newcategorisedresult)
	return newcatresult
}

func RunStaticcheckOnRepo(directory string) (string, error) {
	var cmd *exec.Cmd

	if _, err := os.Stat(filepath.Join(directory, "go.mod")); err == nil {
		cmd = exec.Command("staticcheck", "-f", "json", "./...")
		cmd.Dir = directory
	} else {
		files := []string{}
		err := filepath.Walk(directory, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if strings.HasSuffix(info.Name(), ".go") {
				files = append(files, path)
			}
			return nil
		})
		if err != nil {
			return "", err
		}
		cmd = exec.Command("staticcheck", "-f", "json")
		cmd.Args = append(cmd.Args, files...)
		cmd.Dir = directory
	}

	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out

	cmd.Run()
	return strings.TrimSpace(out.String()), nil
}

func CategorizeStaticcheckResults(results string, directory string) map[string][]interface{} {
	categorizedResults := utils.InitCategorizedResults()

	if results != "" {
		issues := strings.Split(results, "\n")

		for _, issue := range issues {
			issue = strings.TrimSpace(issue)
			if issue == "" {
				continue
			}

			var jsonIssue map[string]interface{}
			err := json.Unmarshal([]byte(issue), &jsonIssue)
			if err != nil {
				log.Printf("Error parsing issue JSON: %v", err)
				continue
			}

			code, ok := jsonIssue["code"].(string)
			if !ok {
				log.Printf("Invalid issue code: %+v", jsonIssue)
				continue
			}
			fmt.Println(code)
			if strings.HasPrefix(code, "ST") {
				categorizedResults[utils.DOCKSTRING_ABSENT] = append(categorizedResults[utils.DOCKSTRING_ABSENT], FormatIssueForStatic(jsonIssue, directory))
			} else if strings.HasPrefix(code, "SA2") || strings.HasPrefix(code, "SA1") {
				categorizedResults[utils.COMPLEX_FUNCTIONS] = append(categorizedResults[utils.COMPLEX_FUNCTIONS], FormatIssueForStatic(jsonIssue, directory))
			} else {
				log.Printf("Unmatched issue: %+v", jsonIssue)
			}
		}
	} else {
		log.Println("No results found from Staticcheck.")
	}
	fmt.Println(categorizedResults)
	return categorizedResults
}

func FormatIssueForStatic(issue map[string]interface{}, directory string) map[string]interface{} {
	location, ok := issue["location"].(map[string]interface{})
	if !ok {
		log.Printf("Invalid issue location: %+v", issue)
		return nil
	}

	filePath, _ := location["file"].(string)
	line, _ := location["line"].(float64)
	column, _ := location["column"].(float64)
	message, _ := issue["message"].(string)

	filePath = strings.Replace(filePath, directory, "", 1)

	return map[string]interface{}{
		"path":    filePath,
		"line":    int(line),
		"column":  int(column),
		"message": message,
	}
}
