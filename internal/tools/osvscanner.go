package tools

import (
	utils "armur-codescanner/pkg"
	"encoding/json"
	"log"
	"os/exec"
	"strings"
)

func RunOSVScanner(directory string) (map[string]interface{}, error) {
	log.Println("Running OSV-Scanner...")
	result, err := runOSVScannerOnRepo(directory)
	if err != nil {
		return nil, err
	}
	ans := categorizeOSVResults(result, directory)
	//fmt.Println(ans)
	newcat := utils.ConvertCategorizedResults(ans)
	return newcat, nil
}

func runOSVScannerOnRepo(directory string) (string, error) {
	cmd := exec.Command("osv-scanner", "--format", "json", directory)
	output, _ := cmd.Output()
	return string(output), nil
}

func categorizeOSVResults(results string, directory string) map[string][]interface{} {
	categorizedResults := utils.InitAdvancedCategorizedResults()

	if results == "" {
		return categorizedResults
	}

	var osvResults struct {
		Results []struct {
			Source struct {
				Path string `json:"path"`
			} `json:"source"`
			Packages []struct {
				Package struct {
					Name    string `json:"name"`
					Version string `json:"version"`
				} `json:"package"`
				Vulnerabilities []struct {
					ID       string      `json:"id"`
					Summary  string      `json:"summary"`
					Severity interface{} `json:"severity"`
				} `json:"vulnerabilities"`
			} `json:"packages"`
		} `json:"results"`
	}

	err := json.Unmarshal([]byte(results), &osvResults)
	if err != nil {
		log.Printf("Failed to parse OSV-Scanner output: %v", err)
		return categorizedResults
	}

	for _, result := range osvResults.Results {
		for _, packageData := range result.Packages {
			for _, vulnerability := range packageData.Vulnerabilities {
				categorizedVuln := map[string]interface{}{
					"path":     strings.Replace(result.Source.Path, directory, "", 1),
					"package":  packageData.Package.Name,
					"version":  packageData.Package.Version,
					"check_id": vulnerability.ID,
					"message":  vulnerability.Summary,
				}

				switch severity := vulnerability.Severity.(type) {
				case string:
					categorizedVuln["severity"] = severity
				case []interface{}:
					severities := make([]string, len(severity))
					for i, v := range severity {
						if s, ok := v.(string); ok {
							severities[i] = s
						}
					}
					categorizedVuln["severity"] = severities
				default:
					categorizedVuln["severity"] = "Unknown"
				}

				categorizedResults[utils.SCA] = append(categorizedResults[utils.SCA], categorizedVuln)
			}
		}
	}

	if len(categorizedResults[utils.SCA]) == 0 {
		categorizedResults[utils.SCA] = []interface{}{}
	}

	return categorizedResults
}
