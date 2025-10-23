package tools

import (
	utils "armur-codescanner/pkg"
	"encoding/json"
	"fmt"
	"log"
	"os/exec"
	"strings"
)

func RunSemgrep(directory string, rules string) map[string]interface{} {
	log.Println("Running Semgrep...")
	semgrepResults := runSemgrepOnRepo(directory, "--config=auto")
	categorizedResults := CategorizeSemgrepResults(semgrepResults, directory)
	newcatresult := utils.ConvertCategorizedResults(categorizedResults)
	fmt.Println(newcatresult)
	return newcatresult
}

func RunSemgrepSolidity(directory string, extraArgs ...string) map[string]interface{} {
	log.Println("Running Semgrep with Solidity rules...")

	args := []string{"--json", "--config=rule_config/semgrep/solidity"}
	if len(extraArgs) > 0 {
		args = append(args, extraArgs...)
	}
	args = append(args, directory)

	out, err := runSemgrepDockerSolidity(directory, args)
	if err != nil {
		out = runSemgrepNativeSolidity(args)
	}
	return categorizeSemgrepSolidity(out, directory)
}

func runSemgrepDockerSolidity(directory string, args []string) (string, error) {
	// Check docker availability
	if _, err := exec.LookPath("docker"); err != nil {
		return "", fmt.Errorf("docker not found")
	}
	// Compose docker command: mount directory to /src and run semgrep there
	dockerArgs := []string{"run", "--rm", "-v", fmt.Sprintf("%s:/src", directory), "-w", "/src", "returntocorp/semgrep:latest", "semgrep"}
	dockerArgs = append(dockerArgs, args...)
	cmd := exec.Command("docker", dockerArgs...)
	output, err := cmd.Output()
	return string(output), err
}

func runSemgrepNativeSolidity(args []string) string {
	cmd := exec.Command("semgrep", args...)
	output, _ := cmd.Output()
	return string(output)
}

func categorizeSemgrepSolidity(jsonOut string, directory string) map[string]interface{} {
	categorized := utils.InitCategorizedResults()
	if strings.TrimSpace(jsonOut) == "" {
		return utils.ConvertCategorizedResults(categorized)
	}

	var parsed map[string]interface{}
	if err := json.Unmarshal([]byte(jsonOut), &parsed); err != nil {
		log.Printf("Error parsing Semgrep Solidity output: %v", err)
		return utils.ConvertCategorizedResults(categorized)
	}

	results, _ := parsed["results"].([]interface{})
	for _, res := range results {
		result, ok := res.(map[string]interface{})
		if !ok {
			continue
		}

		path := safeString(result["path"])
		checkID := safeString(result["check_id"])
		message := safeString(result["extra"].(map[string]interface{})["message"])
		severity := formatSeverity(result)

		relPath := path
		if strings.HasPrefix(path, directory) {
			relPath = strings.TrimPrefix(path, directory+"/")
		}

		issue := map[string]interface{}{
			"path":     relPath,
			"line":     safeString(result["start"].(map[string]interface{})["line"]),
			"message":  message,
			"severity": severity,
			"rule":     checkID,
			"tool":     "semgrep-solidity",
		}

		categorized[utils.SECURITY_ISSUES] = append(categorized[utils.SECURITY_ISSUES], issue)
	}

	return utils.ConvertCategorizedResults(categorized)
}

func runSemgrepOnRepo(directory string, rules string) string {
	cmd := exec.Command("semgrep", rules, directory, "--json")
	output, _ := cmd.Output()
	return string(output)
}

func CategorizeSemgrepResults(results string, directory string) map[string][]interface{} {
	categorizedResults := utils.InitCategorizedResults()

	// Load CWE data
	data, err := utils.LoadCWEData("pkg/common/cwd.json")
	if err != nil {
		log.Printf("Error loading CWE data: %v", err)
		return categorizedResults
	}

	if results != "" {
		var parsedResults map[string]interface{}
		err := json.Unmarshal([]byte(results), &parsedResults)
		if err != nil {
			log.Printf("Error parsing Semgrep results: %v", err)
			log.Printf("Raw Semgrep results: %s", results) // Log raw input for debugging
			return categorizedResults
		}

		resultsArr, ok := parsedResults["results"].([]interface{})
		if !ok {
			log.Println("No 'results' array found in Semgrep output.")
			return categorizedResults
		}

		for _, res := range resultsArr {
			result, ok := res.(map[string]interface{})
			if !ok {
				log.Println("Error processing individual result, skipping.")
				continue
			}

			checkID, ok := result["check_id"].(string)
			if !ok {
				log.Println("Missing or invalid 'check_id', skipping result.")
				continue
			}
			path, ok := result["path"].(string)
			if !ok {
				log.Println("Missing or invalid 'path', skipping result.")
				continue
			}

			path = strings.Replace(path, directory, "", 1)

			switch {
			case strings.Contains(strings.ToLower(checkID), "docstring"):
				categorizedResults[DOCKSTRING_ABSENT] = append(categorizedResults[DOCKSTRING_ABSENT], result)

			case strings.Contains(strings.ToLower(checkID), "security"):
				securityResult := make(map[string]interface{})
				securityResult["path"] = path

				// Extract other fields with proper type assertions
				if start, ok := result["start"].(map[string]interface{}); ok {
					securityResult["line"] = start["line"]
					securityResult["column"] = start["col"]
				}
				if end, ok := result["end"].(map[string]interface{}); ok {
					securityResult["endLine"] = end["line"]
				}
				if extra, ok := result["extra"].(map[string]interface{}); ok {
					securityResult["message"] = extra["message"]
					securityResult["severity"] = formatSeverity(result)
					if metadata, ok := extra["metadata"].(map[string]interface{}); ok {
						securityResult["confidence"] = metadata["confidence"]
						securityResult["likelihood"] = metadata["likelihood"]
						securityResult["cwe"] = metadata["cwe"]
						securityResult["owasp"] = metadata["owasp"]
					}
				}

				// Update practices based on CWE
				if cwe, ok := securityResult["cwe"].([]interface{}); ok && len(cwe) > 0 {
					practices := utils.GetPracticesFromJSON(data, utils.DetectFileLanguage(path), cwe[0].(string))
					for key, value := range practices {
						securityResult[key] = value
					}
				}

				categorizedResults[SECURITY_ISSUES] = append(categorizedResults[SECURITY_ISSUES], securityResult)

			case strings.Contains(strings.ToLower(checkID), "complexity"):
				categorizedResults[COMPLEX_FUNCTIONS] = append(categorizedResults[COMPLEX_FUNCTIONS], result)

			default:
				// Handle antipattern bugs
				antipatternBug := map[string]interface{}{
					"check_id": result["check_id"],
					"column":   fmt.Sprintf("%v", result["start"].(map[string]interface{})["col"]),
					"line":     fmt.Sprintf("%v", result["start"].(map[string]interface{})["line"]),
					"message":  result["extra"].(map[string]interface{})["message"],
					"path":     path,
				}
				categorizedResults[ANTIPATTERNS_BUGS] = append(categorizedResults[ANTIPATTERNS_BUGS], antipatternBug)
			}
		}
	}
	return categorizedResults
}

func formatSeverity(result map[string]interface{}) string {
	extra, ok := result["extra"].(map[string]interface{})
	if !ok {
		return "UNKNOWN"
	}
	originalSeverity, ok := extra["severity"].(string)
	if !ok {
		return "UNKNOWN"
	}
	switch strings.ToUpper(originalSeverity) {
	case "INFO":
		return "LOW"
	case "WARNING":
		return "MEDIUM"
	case "ERROR":
		return "HIGH"
	default:
		return originalSeverity
	}
}
