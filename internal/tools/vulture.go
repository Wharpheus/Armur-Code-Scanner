package internal

import (
	"armur-codescanner/pkg"
	"os/exec"
	"strings"
)

func RunVulture(directory string) (map[string]interface{}, error) {
	vultureResults, err := runVultureOnRepo(directory)
	if err != nil {
		return nil, err
	}
	ans := categorizeVultureResults(vultureResults, directory)
	newcat := utils.ConvertCategorizedResults(ans)
	return newcat, nil
}

func runVultureOnRepo(directory string) (string, error) {
	cmd := exec.Command("vulture", directory)
	output, _ := cmd.Output()
	return string(output), nil
}

func categorizeVultureResults(vultureResults string, directory string) map[string][]interface{} {
	categorizedResults := utils.InitAdvancedCategorizedResults()

	if vultureResults != "" {
		var results []map[string]interface{}
		lines := strings.Split(vultureResults, "\n")

		for _, line := range lines {
			match := strings.SplitN(line, ":", 3)
			if len(match) == 3 {
				// Create a result map for each line
				result := map[string]interface{}{
					"file":    strings.Replace(match[0], directory, "", 1),
					"line":    match[1],
					"message": strings.TrimSpace(match[2]),
				}
				results = append(results, result)
			}
		}

		categorizedResults[utils.DEAD_CODE] = make([]interface{}, len(results))
		for i, res := range results {
			categorizedResults[utils.DEAD_CODE][i] = res
		}
	}

	return categorizedResults
}
