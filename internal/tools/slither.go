package tools

import (
	utils "armur-codescanner/pkg"
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

const (
	GAS_ISSUES = "gas_issues"
)

// RunSlither executes Slither on the given directory and converts the output to Armur categories.
// It prefers dockerized execution and falls back to native slither if unavailable.
// Enhanced with dependency checks for known vulnerable contracts.
func RunSlither(directory string, extraArgs ...string) map[string]interface{} {
	log.Println("Running Slither...")

	// Build common args
	args := []string{"--json", "-", "--detect", "all"}
	if len(extraArgs) > 0 {
		args = append(args, extraArgs...)
	}
	args = append(args, directory)

	// Try docker first
	out, err := runSlitherDocker(directory, args)
	if err != nil || strings.TrimSpace(out) == "" {
		// Fallback to native
		out = runSlitherNative(args)
	}
	return categorizeSlither(out, directory)
}

func RunSlitherDependencies(directory string) map[string]interface{} {
	log.Println("Running Slither dependency checks...")

	args := []string{"--json", "-", "--check-known-thefts"}

	out, err := runSlitherDocker(directory, args)
	if err != nil || strings.TrimSpace(out) == "" {
		out = runSlitherNative(args)
	}
	return categorizeSlitherDependencies(out, directory)
}

func runSlitherNative(args []string) string {
	cmd := exec.Command("slither", args...)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	_ = cmd.Run()
	if s := stdout.String(); strings.TrimSpace(s) != "" {
		return s
	}
	return stderr.String()
}

func runSlitherDocker(directory string, args []string) (string, error) {
	// Check docker availability
	if _, err := exec.LookPath("docker"); err != nil {
		return "", fmt.Errorf("docker not found")
	}
	// Compose docker command: mount directory to /src and run slither there
	dockerArgs := []string{"run", "--rm", "-v", fmt.Sprintf("%s:/src", directory), "-w", "/src", "trailofbits/slither:latest", "slither"}
	dockerArgs = append(dockerArgs, utils.ReplaceTarget(args, directory, ".")...)
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
	return "", fmt.Errorf("no output from slither docker")
}

// categorizeSlither parses Slither JSON and maps findings to Armur categories.
func categorizeSlither(jsonOut string, directory string) map[string]interface{} {
	categorized := utils.InitCategorizedResults()
	if strings.TrimSpace(jsonOut) == "" {
		return utils.ConvertCategorizedResults(categorized)
	}

	var parsed map[string]interface{}
	if err := json.Unmarshal([]byte(jsonOut), &parsed); err != nil {
		log.Printf("Error parsing Slither output: %v", err)
		return utils.ConvertCategorizedResults(categorized)
	}

	results, _ := parsed["results"].(map[string]interface{})
	detectors, _ := results["detectors"].([]interface{})
	for _, d := range detectors {
		det, _ := d.(map[string]interface{})
		impact := strings.ToUpper(safeString(det["impact"]))
		rule := safeString(det["check"])              // rule id
		description := safeString(det["description"]) // message
		for _, elem := range toSlice(det["elements"]) {
			e := elem.(map[string]interface{})
			sourceMapping, _ := e["source_mapping"].(map[string]interface{})
			filename := safeString(sourceMapping["filename_absolute"])
			if filename == "" {
				filename = safeString(sourceMapping["filename"])
			}
			start, _ := extractLines(sourceMapping)
			line := fmt.Sprintf("%d", start) // Slither can provide an array or first line

			rel := filename
			if rel != "" {
				if r, err := filepath.Rel(directory, filename); err == nil {
					rel = r
				}
			}

			severity := mapSeverity(impact)
			issue := map[string]interface{}{
				"path":     rel,
				"line":     line,
				"message":  description,
				"severity": severity,
				"rule":     rule,
				"tool":     "slither",
			}

			// Classify gas vs security roughly by rule prefix
			bucket := utils.SECURITY_ISSUES
			if strings.Contains(strings.ToLower(rule), "gas") || strings.Contains(strings.ToLower(description), "gas") {
				bucket = GAS_ISSUES
			}
			categorized[bucket] = append(categorized[bucket], issue)
		}
	}

	return utils.ConvertCategorizedResults(categorized)
}

func categorizeSlitherDependencies(jsonOut string, directory string) map[string]interface{} {
	categorized := utils.InitCategorizedResults()
	if strings.TrimSpace(jsonOut) == "" {
		return utils.ConvertCategorizedResults(categorized)
	}

	var parsed map[string]interface{}
	if err := json.Unmarshal([]byte(jsonOut), &parsed); err != nil {
		log.Printf("Error parsing Slither dependencies output: %v", err)
		return utils.ConvertCategorizedResults(categorized)
	}

	// Slither dependency check output structure
	results, _ := parsed["results"].(map[string]interface{})
	detectors, _ := results["detectors"].([]interface{})
	for _, d := range detectors {
		det, _ := d.(map[string]interface{})
		description := safeString(det["description"])

		for _, elem := range toSlice(det["elements"]) {
			e := elem.(map[string]interface{})
			sourceMapping, _ := e["source_mapping"].(map[string]interface{})
			filename := safeString(sourceMapping["filename_absolute"])
			if filename == "" {
				filename = safeString(sourceMapping["filename"])
			}
			line := safeIntString(sourceMapping["lines"])

			rel := filename
			if rel != "" {
				if r, err := filepath.Rel(directory, filename); err == nil {
					rel = r
				}
			}

			issue := map[string]interface{}{
				"path":     rel,
				"line":     line,
				"message":  description,
				"severity": "HIGH", // Dependency issues are typically high severity
				"rule":     "known-vulnerable-dependency",
				"tool":     "slither-dependencies",
			}

			categorized[utils.SECURITY_ISSUES] = append(categorized[utils.SECURITY_ISSUES], issue)
		}
	}

	return utils.ConvertCategorizedResults(categorized)
}

func mapSeverity(impact string) string {
	switch strings.ToUpper(impact) {
	case "CRITICAL":
		return "CRITICAL"
	case "HIGH":
		return "HIGH"
	case "MEDIUM":
		return "MEDIUM"
	case "LOW":
		return "LOW"
	default:
		return "INFO"
	}
}

func safeString(v interface{}) string {
	if v == nil {
		return ""
	}
	switch t := v.(type) {
	case string:
		return t
	default:
		b, _ := json.Marshal(t)
		return string(b)
	}
}

func safeIntString(v interface{}) string {
	switch t := v.(type) {
	case float64:
		return strings.TrimSuffix(strings.TrimPrefix(strings.TrimSpace(jsonNumber(t)), "\""), "\"")
	case []interface{}:
		if len(t) > 0 {
			return safeIntString(t[0])
		}
	}
	return "-"
}

func jsonNumber(f float64) string {
	b, _ := json.Marshal(f)
	return string(b)
}

func toSlice(v interface{}) []interface{} {
	if s, ok := v.([]interface{}); ok {
		return s
	}
	return []interface{}{}
}

// extractLines attempts to return start and end line numbers from slither source_mapping.
func extractLines(sm map[string]interface{}) (int, int) {
	// Preferred: "lines" is an array of numbers; take first and last
	if lines, ok := sm["lines"].([]interface{}); ok && len(lines) > 0 {
		first := toInt(lines[0])
		last := toInt(lines[len(lines)-1])
		if last < first {
			last = first
		}
		return first, last
	}
	// Fallback: "lines" might be single number
	if l, ok := sm["lines"].(float64); ok {
		v := int(l)
		return v, v
	}
	// As a last resort, return 0,0 meaning unknown
	return 0, 0
}

func toInt(v interface{}) int {
	switch t := v.(type) {
	case float64:
		return int(t)
	case string:
		if n, err := strconv.Atoi(t); err == nil {
			return n
		}
	}
	return 0
}
