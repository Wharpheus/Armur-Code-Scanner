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

// RunOyente executes Oyente symbolic execution analysis on the given directory.
// It uses Docker for execution and converts output to Armur categories.
func RunOyente(directory string, extraArgs ...string) map[string]interface{} {
	log.Println("Running Oyente...")

	// Build args: Oyente expects a single .sol file or directory
	args := []string{"-s", directory}
	if len(extraArgs) > 0 {
		args = append(args, extraArgs...)
	}

	// Run via Docker
	out, err := runOyenteDocker(directory, args)
	if err != nil {
		log.Printf("Oyente Docker failed: %v", err)
		return utils.ConvertCategorizedResults(utils.InitCategorizedResults())
	}

	return categorizeOyente(out, directory)
}

func runOyenteDocker(directory string, args []string) (string, error) {
	if _, err := exec.LookPath("docker"); err != nil {
		return "", fmt.Errorf("docker not found")
	}

	// Mount directory and run Oyente
	dockerArgs := []string{"run", "--rm", "-v", fmt.Sprintf("%s:/src", directory), "-w", "/src", "luongnguyen/oyente:latest"}
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
	return "", fmt.Errorf("no output from Oyente Docker")
}

// categorizeOyente parses Oyente output and maps to Armur categories.
// Oyente outputs JSON with vulnerabilities found.
func categorizeOyente(jsonOut string, directory string) map[string]interface{} {
	categorized := utils.InitCategorizedResults()
	if strings.TrimSpace(jsonOut) == "" {
		return utils.ConvertCategorizedResults(categorized)
	}

	var parsed map[string]interface{}
	if err := json.Unmarshal([]byte(jsonOut), &parsed); err != nil {
		log.Printf("Error parsing Oyente output: %v", err)
		return utils.ConvertCategorizedResults(categorized)
	}

	// Oyente output structure: {"vulnerabilities": {...}}
	vulns, _ := parsed["vulnerabilities"].(map[string]interface{})
	for vulnType, details := range vulns {
		detailMap, ok := details.(map[string]interface{})
		if !ok {
			continue
		}

		// Extract locations
		locations, _ := detailMap["locations"].([]interface{})
		for _, loc := range locations {
			locMap, ok := loc.(map[string]interface{})
			if !ok {
				continue
			}

			file := safeString(locMap["file"])
			line := safeString(locMap["line"])

			relFile := file
			if rel, err := filepath.Rel(directory, file); err == nil {
				relFile = rel
			}

			severity := mapOyenteSeverity(vulnType)
			issue := map[string]interface{}{
				"path":     relFile,
				"line":     line,
				"message":  fmt.Sprintf("Oyente detected %s vulnerability", vulnType),
				"severity": severity,
				"rule":     vulnType,
				"tool":     "oyente",
			}

			categorized[utils.SECURITY_ISSUES] = append(categorized[utils.SECURITY_ISSUES], issue)
		}
	}

	return utils.ConvertCategorizedResults(categorized)
}

func mapOyenteSeverity(vulnType string) string {
	switch strings.ToLower(vulnType) {
	case "reentrancy", "integer_overflow", "integer_underflow":
		return "HIGH"
	case "unchecked_call", "suicide":
		return "MEDIUM"
	default:
		return "LOW"
	}
}
