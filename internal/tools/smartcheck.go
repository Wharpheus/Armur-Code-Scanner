package tools

import (
	utils "armur-codescanner/pkg"
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os/exec"
	"path/filepath"
	"strings"
)

// RunSmartCheck executes SmartCheck static analysis on the given directory.
// It uses Docker and converts output to Armur categories.
func RunSmartCheck(directory string, extraArgs ...string) map[string]interface{} {
	log.Println("Running SmartCheck...")

	args := []string{"-p", directory}
	if len(extraArgs) > 0 {
		args = append(args, extraArgs...)
	}

	out, err := runSmartCheckDocker(directory, args)
	if err != nil {
		log.Printf("SmartCheck Docker failed: %v", err)
		return utils.ConvertCategorizedResults(utils.InitCategorizedResults())
	}

	return categorizeSmartCheck(out, directory)
}

func runSmartCheckDocker(directory string, args []string) (string, error) {
	if _, err := exec.LookPath("docker"); err != nil {
		return "", fmt.Errorf("docker not found")
	}

	dockerArgs := []string{"run", "--rm", "-v", fmt.Sprintf("%s:/src", directory), "-w", "/src", "smartcontractsecurity/smartcheck:latest"}
	dockerArgs = append(dockerArgs, args...)

	cmd := exec.Command("docker", dockerArgs...)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	_ = cmd.Run()

	if s := stdout.String(); strings.TrimSpace(s) != "" {
		return s, nil
	}
	if t := stderr.String(); strings.TrimSpace(t) != "" {
		return t, nil
	}
	return "", fmt.Errorf("no output from SmartCheck Docker")
}

// categorizeSmartCheck parses SmartCheck SARIF or JSON output and maps to categories.
func categorizeSmartCheck(jsonOut string, directory string) map[string]interface{} {
	categorized := utils.InitCategorizedResults()
	if strings.TrimSpace(jsonOut) == "" {
		return utils.ConvertCategorizedResults(categorized)
	}

	var parsed map[string]interface{}
	if err := json.Unmarshal([]byte(jsonOut), &parsed); err != nil {
		log.Printf("Error parsing SmartCheck output: %v", err)
		return utils.ConvertCategorizedResults(categorized)
	}

	// SmartCheck SARIF structure: {"runs": [{"results": [...]}]}
	runs, _ := parsed["runs"].([]interface{})
	for _, run := range runs {
		runMap, ok := run.(map[string]interface{})
		if !ok {
			continue
		}

		results, _ := runMap["results"].([]interface{})
		for _, res := range results {
			resMap, ok := res.(map[string]interface{})
			if !ok {
				continue
			}

			ruleID := safeString(resMap["ruleId"])
			message := safeString(resMap["message"])
			severity := safeString(resMap["level"])

			// Extract locations
			locations, _ := resMap["locations"].([]interface{})
			for _, loc := range locations {
				locMap, ok := loc.(map[string]interface{})
				if !ok {
					continue
				}

				physicalLocation, _ := locMap["physicalLocation"].(map[string]interface{})
				artifactLocation, _ := physicalLocation["artifactLocation"].(map[string]interface{})
				file := safeString(artifactLocation["uri"])

				region, _ := physicalLocation["region"].(map[string]interface{})
				line := safeString(region["startLine"])

				relFile := file
				if rel, err := filepath.Rel(directory, file); err == nil {
					relFile = rel
				}

				issue := map[string]interface{}{
					"path":     relFile,
					"line":     line,
					"message":  message,
					"severity": mapSmartCheckSeverity(severity),
					"rule":     ruleID,
					"tool":     "smartcheck",
				}

				categorized[utils.SECURITY_ISSUES] = append(categorized[utils.SECURITY_ISSUES], issue)
			}
		}
	}

	return utils.ConvertCategorizedResults(categorized)
}

func mapSmartCheckSeverity(severity string) string {
	switch strings.ToLower(severity) {
	case "error":
		return "HIGH"
	case "warning":
		return "MEDIUM"
	case "info":
		return "LOW"
	default:
		return "INFO"
	}
}
