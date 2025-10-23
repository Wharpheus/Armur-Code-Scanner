package tools

import (
	utils "armur-codescanner/pkg"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
	"path/filepath"
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
// Now includes AST analysis for deeper insights.
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
	} else {
		// Parse AST for additional analysis
		astIssues := analyzeAST(out, root)
		for _, issue := range astIssues {
			categorized[utils.SECURITY_ISSUES] = append(categorized[utils.SECURITY_ISSUES], issue)
		}
	}
	return utils.ConvertCategorizedResults(categorized)
}

// analyzeAST performs deeper analysis on the Solidity AST to extract insights like function complexity, inheritance, etc.
func analyzeAST(parsed interface{}, root string) []map[string]interface{} {
	var issues []map[string]interface{}

	astMap, ok := parsed.(map[string]interface{})
	if !ok {
		return issues
	}

	sources, ok := astMap["sources"].(map[string]interface{})
	if !ok {
		return issues
	}

	for filePath, sourceData := range sources {
		sourceMap, ok := sourceData.(map[string]interface{})
		if !ok {
			continue
		}

		ast, ok := sourceMap["AST"].(map[string]interface{})
		if !ok {
			continue
		}

		relPath, _ := filepath.Rel(root, filePath)

		// Analyze contract nodes
		contractIssues := analyzeContractNodes(ast, relPath)
		issues = append(issues, contractIssues...)
	}

	return issues
}

// analyzeContractNodes extracts contract-level insights from AST.
func analyzeContractNodes(ast map[string]interface{}, relPath string) []map[string]interface{} {
	var issues []map[string]interface{}

	nodes, ok := ast["nodes"].([]interface{})
	if !ok {
		return issues
	}

	for _, node := range nodes {
		nodeMap, ok := node.(map[string]interface{})
		if !ok {
			continue
		}

		nodeType, _ := nodeMap["nodeType"].(string)
		if nodeType == "ContractDefinition" {
			contractIssues := analyzeContractDefinition(nodeMap, relPath)
			issues = append(issues, contractIssues...)
		}
	}

	return issues
}

// analyzeContractDefinition analyzes a single contract for complexity, inheritance, etc.
func analyzeContractDefinition(contract map[string]interface{}, relPath string) []map[string]interface{} {
	var issues []map[string]interface{}

	name, _ := contract["name"].(string)
	baseContracts, _ := contract["baseContracts"].([]interface{})

	// Check for multiple inheritance (potential complexity)
	if len(baseContracts) > 2 {
		issues = append(issues, map[string]interface{}{
			"path":     relPath,
			"line":     "-",
			"message":  fmt.Sprintf("Contract %s has multiple inheritance (%d bases), consider simplifying", name, len(baseContracts)),
			"severity": "MEDIUM",
			"tool":     "solc-ast",
		})
	}

	// Analyze functions
	subNodes, _ := contract["nodes"].([]interface{})
	functionCount := 0
	for _, subNode := range subNodes {
		subMap, ok := subNode.(map[string]interface{})
		if !ok {
			continue
		}

		if subType, _ := subMap["nodeType"].(string); subType == "FunctionDefinition" {
			functionCount++
			funcIssues := analyzeFunctionDefinition(subMap, relPath, name)
			issues = append(issues, funcIssues...)
		}
	}

	// Flag contracts with too many functions
	if functionCount > 20 {
		issues = append(issues, map[string]interface{}{
			"path":     relPath,
			"line":     "-",
			"message":  fmt.Sprintf("Contract %s has %d functions, consider breaking into smaller contracts", name, functionCount),
			"severity": "LOW",
			"tool":     "solc-ast",
		})
	}

	return issues
}

// analyzeFunctionDefinition analyzes a function for complexity.
func analyzeFunctionDefinition(funcDef map[string]interface{}, relPath, contractName string) []map[string]interface{} {
	var issues []map[string]interface{}

	name, _ := funcDef["name"].(string)
	body, _ := funcDef["body"].(map[string]interface{})

	if body == nil {
		return issues // Abstract function
	}

	// Simple complexity metric: count statements
	statements := countStatements(body)
	if statements > 50 {
		issues = append(issues, map[string]interface{}{
			"path":     relPath,
			"line":     "-",
			"message":  fmt.Sprintf("Function %s in contract %s has %d statements, consider refactoring", name, contractName, statements),
			"severity": "MEDIUM",
			"tool":     "solc-ast",
		})
	}

	return issues
}

// countStatements recursively counts statements in a block.
func countStatements(node map[string]interface{}) int {
	count := 0
	if statements, ok := node["statements"].([]interface{}); ok {
		count += len(statements)
		for _, stmt := range statements {
			if stmtMap, ok := stmt.(map[string]interface{}); ok {
				count += countStatements(stmtMap)
			}
		}
	}
	return count
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
