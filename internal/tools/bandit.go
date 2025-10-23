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

const (
	SECURITY_ISSUES   = "security_issues"
	DOCKSTRING_ABSENT = "dockstring_absent"
	COMPLEX_FUNCTIONS = "complex_functions"
	ANTIPATTERNS_BUGS = "antipatterns_bugs"
)

type BanditIssue struct {
	Path       string   `json:"path"`
	Line       int      `json:"line"`
	Message    string   `json:"message"`
	Severity   string   `json:"severity"`
	TestID     string   `json:"test_id"`
	Confidence string   `json:"confidence"`
	Likelihood string   `json:"likelihood"`
	CWE        []string `json:"cwe"`
	Owasp      []string `json:"owasp"`
	EndLine    int      `json:"endLine"`
}

func RunBandit(directory string) map[string]interface{} {
	log.Println("Running Bandit...")
	results := RunBanditOnRepo(directory)
	categorizedResults := CategorizeBanditResults(results)
	newcatresult := utils.ConvertCategorizedResults(categorizedResults)
	return newcatresult
}

func RunBanditOnRepo(directory string) string {
	cmd := exec.Command("bandit", "-r", directory, "-f", "json")
	var out, stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	cmd.Run()
	return out.String()
}

func CategorizeBanditResults(results string) map[string][]interface{} {
	categorizedResults := utils.InitCategorizedResults()

	// Return early if the results string is empty or whitespace
	if strings.TrimSpace(results) == "" {
		log.Println("No results to categorize.")
		return categorizedResults
	}

	// Parse JSON input
	var parsedResults map[string]interface{}
	err := json.Unmarshal([]byte(results), &parsedResults)
	if err != nil {
		log.Printf("Failed to parse Bandit results: %v\n", err)
		return categorizedResults
	}

	// Process each issue in the results
	if issues, ok := parsedResults["results"].([]interface{}); ok {
		for _, rawIssue := range issues {
			if issueMap, ok := rawIssue.(map[string]interface{}); ok {
				// Extract issue details with type assertions
				fullPath, _ := issueMap["filename"].(string)
				fileName := filepath.Base(fullPath)
				line, _ := issueMap["line_number"].(float64)
				message, _ := issueMap["issue_text"].(string)
				severity, _ := issueMap["issue_severity"].(string)
				testID, _ := issueMap["test_id"].(string)
				confidence, _ := issueMap["issue_confidence"].(string)

				// Convert CWE if available
				cwe := []string{}
				if testID != "" {
					cwe = append(cwe, fmt.Sprintf("CWE-%s: %s", testID, message))
				}

				// Create BanditIssue
				issue := BanditIssue{
					Path:       fileName,
					Line:       int(line),
					Message:    message,
					Severity:   strings.ToUpper(severity),
					TestID:     testID,
					Confidence: strings.ToUpper(confidence),
					Likelihood: "LOW", // Default value
					CWE:        cwe,
					Owasp:      []string{"A07:2017 - Cross-Site Scripting (XSS)", "A03:2021 - Injection"},
					EndLine:    int(line),
				}

				// Convert BanditIssue to map[string]interface{}
				issueMapConverted := map[string]interface{}{
					"path":       issue.Path,
					"line":       issue.Line,
					"message":    issue.Message,
					"severity":   issue.Severity,
					"test_id":    issue.TestID,
					"confidence": issue.Confidence,
					"likelihood": issue.Likelihood,
					"cwe":        issue.CWE,
					"owasp":      issue.Owasp,
					"endLine":    issue.EndLine,
				}

				// Add to the SECURITY_ISSUES category
				categorizedResults[SECURITY_ISSUES] = append(categorizedResults[SECURITY_ISSUES], issueMapConverted)
			}
		}
	}

	return categorizedResults
}
