package tools

import (
	utils "armur-codescanner/pkg"
	"bytes"
	"log"
	"os/exec"
	"regexp"
	"strings"
)

func RunPydocstyle(directory string) map[string]interface{} {
	log.Println("Running Pydocstyle...")
	pydocstyleResults := RunPydocstyleOnRepo(directory)
	categorizedResults := CategorizePydocstyleResults(pydocstyleResults, directory)
	newcattu := utils.ConvertCategorizedResults(categorizedResults)
	return newcattu
}

func RunPydocstyleOnRepo(directory string) string {
	cmd := exec.Command("pydocstyle", directory)
	var out, stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	cmd.Run()
	return out.String()
}

func CategorizePydocstyleResults(pydocstyleResults string, directory string) map[string][]interface{} {
	categorizedResults := utils.InitCategorizedResults()
	if pydocstyleResults != "" {
		re := regexp.MustCompile(`(?P<filename>.*?):(?P<line>\d+)\s(?P<level>.*):\s+D(?P<error_code>\d{3}):\s(?P<message>.*)`)

		matches := re.FindAllStringSubmatch(pydocstyleResults, -1)

		for _, match := range matches {
			if len(match) == 6 {
				result := map[string]interface{}{
					"path":       strings.Replace(match[1], directory, "", 1),
					"line":       match[2],
					"level":      match[3],
					"error_code": "D" + match[4],
					"message":    match[5],
				}
				categorizedResults[DOCKSTRING_ABSENT] = append(categorizedResults[DOCKSTRING_ABSENT], result)
			}
		}
	}

	return categorizedResults
}
