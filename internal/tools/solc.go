package tools

import (
	utils "armur-codescanner/pkg"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"
	"time"
)

type SolcDiagnostic struct {
	Path     string `json:"path"`
	Line     string `json:"line"`
	Message  string `json:"message"`
	Severity string `json:"severity"`
	Tool     string `json:"tool"`
}

// RunSolcCheck performs a lightweight compilation to surface syntax/build issues early.
// It attempts dockerized solc first with a specified version, then falls back to local solc.
func RunSolcCheck(root string, version string, remappings []string) map[string]interface{} {
	categorized := utils.InitCategorizedResults()

	args := []string{"--ast-compact-json"}
	for _, r := range remappings {
		args = append(args, "--allow-paths", ".", "--include-path", ".", "--remap-imports", r)
	}
	// Collect all .sol files via solc recursive flag
	args = append(args, "--pretty-json")
	args = append(args, "--base-path", ".")
	args = append(args, "--")
	args = append(args, "$(find . -type f -name '*.sol' -print)")

	// Prefer dockerized solc if version provided
	var out string
	var err error
	if version != "" {
		out, err = runSolcDocker(root, version, args)
	}
	if err != nil || strings.TrimSpace(out) == "" {
		out = runSolcNative(root, args)
	}

	// solc prints errors to stderr; we only convert non-empty stderr lines to diagnostics if JSON failed.
	// If output is not valid JSON, treat as compiler errors.
	var parsed interface{}
	if json.Unmarshal([]byte(out), &parsed) != nil {
		for _, line := range strings.Split(out, "\n") {
			l := strings.TrimSpace(line)
			if l == "" {
				continue
			}
			categorized[utils.BUILD_ISSUES] = append(categorized[utils.BUILD_ISSUES], map[string]interface{}{
				"path":     "",
				"line":     "-",
				"message":  l,
				"severity": "ERROR",
				"tool":     "solc",
			})
		}
	}
	return utils.ConvertCategorizedResults(categorized)
}

func runSolcNative(workdir string, args []string) string {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, "bash", "-lc", "solc "+strings.Join(args, " "))
	cmd.Dir = workdir
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	_ = cmd.Run()
	if s := stdout.String(); strings.TrimSpace(s) != "" {
		return s
	}
	return stderr.String()
}

func runSolcDocker(root, version string, args []string) (string, error) {
	if _, err := exec.LookPath("docker"); err != nil {
		return "", fmt.Errorf("docker not found")
	}
	dockerArgs := []string{"run", "--rm", "-v", fmt.Sprintf("%s:/src", root), "-w", "/src", fmt.Sprintf("ethereum/solc:%s", version), "bash", "-lc", "solc " + strings.Join(args, " ")}
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
	return "", fmt.Errorf("no output from solc docker")
}
