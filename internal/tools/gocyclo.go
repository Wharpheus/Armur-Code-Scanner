package internal

import (
	utils "armur-codescanner/pkg"
	"bytes"
	"log"
	"os/exec"
	"strings"
)

func RunGocyclo(directory string) map[string]interface{} {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Error while running Gocyclo: %v", r)
		}
	}()

	log.Println("Running Gocyclo")
	gocycloResults, err := RunGoCycloOnRepo(directory)
	if err != nil {
		log.Printf("Error while running Gocyclo: %v", err)
		return nil
	}

	categorizedResults := CategorizeGocycloResults(gocycloResults, directory)
	newcatresult := utils.ConvertCategorizedResults(categorizedResults)
	return newcatresult
}

func RunGoCycloOnRepo(directory string) (string, error) {
	cmd := exec.Command("gocyclo", directory)

	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	cmd.Run()
	return out.String(), nil
}

func CategorizeGocycloResults(results string, directory string) map[string][]interface{} {
	categorizedResults := utils.InitCategorizedResults()

	if results != "" {
		lines := strings.Split(results, "\n")
		for _, line := range lines {
			parts := strings.Fields(line)
			if len(parts) < 4 {
				continue
			}

			complexity := parts[0]
			pkg := parts[1]
			function := parts[2]

			locationParts := strings.Split(parts[3], ":")
			if len(locationParts) != 3 {
				log.Printf("Invalid location format: %s", parts[3])
				continue
			}

			filePath := strings.Replace(locationParts[0], directory, "", 1)
			lineNumber := locationParts[1]
			columnNumber := locationParts[2]

			resultEntry := map[string]interface{}{
				"complexity": complexity,
				"function":   function,
				"package":    pkg,
				"path":       filePath,
				"line":       lineNumber,
				"column":     columnNumber,
			}

			categorizedResults[utils.COMPLEX_FUNCTIONS] = append(categorizedResults[utils.COMPLEX_FUNCTIONS], resultEntry)
		}
	} else {
		log.Println("No results found from Gocyclo.")
	}

	return categorizedResults
}
