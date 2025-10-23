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

// RunMythril executes Mythril analysis on the given directory or file and converts the output to Armur categories.
// It prefers dockerized execution and falls back to native myth if unavailable.
func RunMythril(targetPath string, extraArgs ...string) map[string]interface{} {
	log.Println("Running Mythril...")

	// Build default args
	args := []string{"analyze", "-o", "json", "-x"}
	if len(extraArgs) > 0 {
		args = append(args, extraArgs...)
	}
	args = append(args, targetPath)

	// Try docker first
	out, err := runMythrilDocker(targetPath, args)
	if err != nil || strings.TrimSpace(out) == "" {
		out = runMythrilNative(args)
	}
	return categorizeMythril(out, filepath.Dir(targetPath))
}

func runMythrilNative(args []string) string {
	cmd := exec.Command("myth", args...)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	_ = cmd.Run()
	if s := stdout.String(); strings.TrimSpace(s) != "" {
		return s
	}
	return stderr.String()
}

func runMythrilDocker(targetPath string, args []string) (string, error) {
	if _, err := exec.LookPath("docker"); err != nil {
		return "", fmt.Errorf("docker not found")
	}
	hostDir := targetPath
	if !isDir(targetPath) {
		hostDir = filepath.Dir(targetPath)
	}
	relTarget := targetPath
	if r, err := filepath.Rel(hostDir, targetPath); err == nil {
		relTarget = r
	}
	dockerArgs := []string{"run", "--rm", "-v", fmt.Sprintf("%s:/src", hostDir), "-w", "/src", "mythril/myth", "myth"}
	dockerArgs = append(dockerArgs, replaceTarget(args, targetPath, relTarget)...)
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
	return "", fmt.Errorf("no output from mythril docker")
}

func replaceTarget(args []string, oldTarget, newTarget string) []string {
	res := make([]string, len(args))
	copy(res, args)
	if len(res) > 0 && res[len(res)-1] == oldTarget {
		res[len(res)-1] = newTarget
	}
	return res
}

func isDir(path string) bool {
	// simplistic check: if it ends with .sol treat as file
	return !strings.HasSuffix(strings.ToLower(path), ".sol")
}

// categorizeMythril parses Mythril JSON and maps findings to Armur categories.
func categorizeMythril(jsonOut string, baseDir string) map[string]interface{} {
	categorized := utils.InitCategorizedResults()
	if strings.TrimSpace(jsonOut) == "" {
		return utils.ConvertCategorizedResults(categorized)
	}

	var parsed map[string]interface{}
	if err := json.Unmarshal([]byte(jsonOut), &parsed); err != nil {
		// Mythril sometimes prints multiple JSON objects or extra logs; try to find last JSON block
		idx := strings.LastIndex(jsonOut, "}")
		if idx > 0 {
			trim := jsonOut[:idx+1]
			if err2 := json.Unmarshal([]byte(trim), &parsed); err2 != nil {
				log.Printf("Error parsing Mythril output: %v", err)
				return utils.ConvertCategorizedResults(categorized)
			}
		} else {
			log.Printf("Error parsing Mythril output: %v", err)
			return utils.ConvertCategorizedResults(categorized)
		}
	}

	issues, _ := parsed["issues"].([]interface{})
	for _, it := range issues {
		iss, _ := it.(map[string]interface{})
		title := safeString(iss["title"])
		severity := mapMythrilSeverity(safeString(iss["severity"]))
		description := safeString(iss["description"])

		var path, line string
		if locs, ok := iss["locations"].([]interface{}); ok && len(locs) > 0 {
			if loc, ok2 := locs[0].(map[string]interface{}); ok2 {
				file := safeString(loc["file"])
				if file != "" {
					if r, err := filepath.Rel(baseDir, file); err == nil {
						path = r
					} else {
						path = file
					}
				}
				line = safeString(loc["line"])
			}
		}

		issue := map[string]interface{}{
			"path":     path,
			"line":     line,
			"message":  fmt.Sprintf("%s: %s", title, description),
			"severity": severity,
			"rule":     title,
			"tool":     "mythril",
		}
		categorized[SECURITY_ISSUES] = append(categorized[SECURITY_ISSUES], issue)
	}

	return utils.ConvertCategorizedResults(categorized)
}

func mapMythrilSeverity(s string) string {
	sw := strings.ToUpper(strings.TrimSpace(s))
	switch sw {
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
