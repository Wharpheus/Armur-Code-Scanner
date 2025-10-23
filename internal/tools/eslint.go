package tools

import (
	utils "armur-codescanner/pkg"
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os/exec"
	"strings"
)

type Issue struct {
	Path    string `json:"path"`
	Line    int    `json:"line"`
	Column  int    `json:"column,omitempty"`
	Message string `json:"message"`
	RuleID  string `json:"ruleId,omitempty"`
}

func RunESLintAdvanced(directory string) map[string]interface{} {
	log.Println("Running ESLint Advanced")
	results, err := RunESLintAdvancedOnRepo(directory)
	if err != nil {
		log.Printf("Error while running ESLint Advanced: %v\n", err)
		return map[string]interface{}{}
	}
	newcatresult := utils.ConvertCategorizedResults(results)
	return newcatresult
}

func RunESLintOnRepo(repoPath string) map[string]interface{} {
	categorizedResults := utils.InitCategorizedResults()

	// Docstring check
	log.Println("Running ESLint for docstrings")
	docResults, _ := RunESLint(repoPath, "eslint_jsdoc.config.js")
	categorizedResults[DOCKSTRING_ABSENT] = FormatIssues(docResults, repoPath)

	// Security check
	log.Println("Running ESLint for security issues")
	securityResults, _ := RunESLint(repoPath, "eslint_security.config.js")
	categorizedResults[SECURITY_ISSUES] = FormatIssues(securityResults, repoPath)

	// Complex functions and antipatterns
	log.Println("Running ESLint for complex functions and antipatterns")
	complexResults, _ := RunESLint(repoPath, "eslint.config.js")
	categorizedComplex := CategorizeESLintResults(complexResults, repoPath)
	categorizedResults[COMPLEX_FUNCTIONS] = categorizedComplex[COMPLEX_FUNCTIONS]
	categorizedResults[ANTIPATTERNS_BUGS] = categorizedComplex[ANTIPATTERNS_BUGS]
	newcatresult := utils.ConvertCategorizedResults(categorizedResults)
	return newcatresult

}

func RunESLintAdvancedOnRepo(repoPath string) (map[string][]interface{}, error) {
	categorizedResults := utils.InitAdvancedCategorizedResults()

	// Run ESLint for dead code
	cmd := exec.Command("eslint", "--format", "json", "--config", "eslint_deadcode.config.js", repoPath)
	output, _ := cmd.Output()
	fmt.Println(string(output))
	var eslintOutput []map[string]interface{}
	err := json.Unmarshal(output, &eslintOutput)
	if err != nil {
		return nil, err
	}

	categorized := CategorizeESLintAdvancedResults(eslintOutput, repoPath)
	categorizedResults[utils.DEAD_CODE] = categorized[utils.DEAD_CODE]

	return categorizedResults, nil
}

func CategorizeESLintAdvancedResults(eslintResults []map[string]interface{}, directory string) map[string][]interface{} {
	categorizedResults := utils.InitAdvancedCategorizedResults()

	deadCodeESLintRuleIDs := []string{
		"no-unused-vars",
		"no-unreachable",
		"no-constant-condition",
		"no-unused-expressions",
		"no-unused-private-class-members",
		"no-useless-assignment",
	}

	for _, fileResult := range eslintResults {
		for _, message := range fileResult["messages"].([]interface{}) {
			msg := message.(map[string]interface{})
			if checkID, ok := msg["ruleId"].(string); ok {
				msg["check_id"] = checkID
				delete(msg, "ruleId")
				msg["file"] = strings.ReplaceAll(fileResult["filePath"].(string), directory, "")

				if contains(deadCodeESLintRuleIDs, checkID) {
					categorizedResults[utils.DEAD_CODE] = append(categorizedResults[utils.DEAD_CODE], msg)
				}
			}
		}
	}
	return categorizedResults
}

func contains(slice []string, item string) bool {
	for _, val := range slice {
		if val == item {
			return true
		}
	}
	return false
}

func RunESLint(directory, configFile string) ([]map[string]interface{}, error) {
	cmd := exec.Command("eslint", "--format", "json", "--config", configFile, directory)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		log.Println("error running ESLint: ", err)
	}
	var results []map[string]interface{}
	err = json.Unmarshal(stdout.Bytes(), &results)
	if err != nil {
		log.Printf("Error parsing ESLint results: %v", err)
		return nil, err
	}

	return results, nil
}

func FormatIssues(issues []map[string]interface{}, directory string) []interface{} {
	var formattedIssues []interface{}
	for _, issue := range issues {
		filePath := strings.ReplaceAll(issue["filePath"].(string), directory, "")
		messages := issue["messages"].([]interface{})

		for _, msg := range messages {
			message := msg.(map[string]interface{})
			formattedIssues = append(formattedIssues, map[string]interface{}{
				"path":    filePath,
				"line":    int(message["line"].(float64)),
				"column":  int(message["column"].(float64)),
				"message": message["message"].(string),
				"ruleId":  message["ruleId"],
			})
		}
	}
	return formattedIssues
}

func CategorizeESLintResults(results []map[string]interface{}, directory string) map[string][]interface{} {
	categorized := utils.InitCategorizedResults()

	antipatternsRuleIDs := map[string]bool{
		"no-console":           true,
		"strict":               true,
		"no-magic-numbers":     true,
		"consistent-return":    true,
		"no-else-return":       true,
		"no-use-before-define": true,
		"prefer-const":         true,
		"no-restricted-syntax": true,
		"eqeqeq":               true,
		"complexity":           true,
	}

	for _, fileResult := range results {
		filePath := strings.ReplaceAll(fileResult["filePath"].(string), directory, "")
		messages := fileResult["messages"].([]interface{})

		for _, msg := range messages {
			message := msg.(map[string]interface{})
			ruleID, ok := message["ruleId"].(string)
			if !ok {
				continue
			}

			issue := map[string]interface{}{
				"path":    filePath,
				"line":    int(message["line"].(float64)),
				"message": message["message"].(string),
				"ruleId":  ruleID,
			}

			if antipatternsRuleIDs[ruleID] {
				categorized[ANTIPATTERNS_BUGS] = append(categorized[ANTIPATTERNS_BUGS], issue)
			} else if ruleID == "complexity" {
				categorized[COMPLEX_FUNCTIONS] = append(categorized[COMPLEX_FUNCTIONS], issue)
			}
		}
	}
	return categorized
}
