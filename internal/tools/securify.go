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

// RunSecurify executes Securify property verification on the given directory.
// It uses Docker and converts output to Armur categories.
func RunSecurify(directory string, extraArgs ...string) map[string]interface{} {
	log.Println("Running Securify...")

	// Securify analyzes bytecode, so we need to compile first or assume .sol files
	args := []string{"--input", directory}
	if len(extraArgs) > 0 {
		args = append(args, extraArgs...)
	}

	out, err := runSecurifyDocker(directory, args)
	if err != nil {
		log.Printf("Securify Docker failed: %v", err)
		return utils.ConvertCategorizedResults(utils.InitCategorizedResults())
	}

	return categorizeSecurify(out, directory)
}

func runSecurifyDocker(directory string, args []string) (string, error) {
	if _, err := exec.LookPath("docker"); err != nil {
		return "", fmt.Errorf("docker not found")
	}

	dockerArgs := []string{"run", "--rm", "-v", fmt.Sprintf("%s:/src", directory), "-w", "/src", "trailofbits/securify:latest"}
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
	return "", fmt.Errorf("no output from Securify Docker")
}

// categorizeSecurify parses Securify JSON output and maps to categories.
// Securify outputs patterns and violations.
func categorizeSecurify(jsonOut string, directory string) map[string]interface{} {
	categorized := utils.InitCategorizedResults()
	if strings.TrimSpace(jsonOut) == "" {
		return utils.ConvertCategorizedResults(categorized)
	}

	var parsed map[string]interface{}
	if err := json.Unmarshal([]byte(jsonOut), &parsed); err != nil {
		log.Printf("Error parsing Securify output: %v", err)
		return utils.ConvertCategorizedResults(categorized)
	}

	// Securify structure: {"patterns": {...}, "violations": {...}}
	violations, _ := parsed["violations"].(map[string]interface{})
	for pattern, details := range violations {
		detailMap, ok := details.(map[string]interface{})
		if !ok {
			continue
		}

		// Extract locations if available
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

			severity := mapSecurifySeverity(pattern)
			issue := map[string]interface{}{
				"path":     relFile,
				"line":     line,
				"message":  fmt.Sprintf("Securify violation: %s", pattern),
				"severity": severity,
				"rule":     pattern,
				"tool":     "securify",
			}

			categorized[utils.SECURITY_ISSUES] = append(categorized[utils.SECURITY_ISSUES], issue)
		}
	}

	return utils.ConvertCategorizedResults(categorized)
}

func mapSecurifySeverity(pattern string) string {
	switch strings.ToLower(pattern) {
	case "reentrancy", "unchecked_send", "suicidal":
		return "HIGH"
	case "locked_ether", "dao":
		return "MEDIUM"
	default:
		return "LOW"
	}
}
