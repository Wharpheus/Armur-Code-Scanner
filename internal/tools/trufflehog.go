package tools

import (
	// Replace with the actual import path for utils
	utils "armur-codescanner/pkg"
	"encoding/json"
	"log"
	"os/exec"
	"strconv"
	"strings"
)

type Secret struct {
	Path string      `json:"path"`
	Line interface{} `json:"line"`
	Rule struct {
		Message  string `json:"message"`
		Severity string `json:"severity"`
	} `json:"rule"`
	Secret string `json:"secret"`
}

// Run TruffleHog on the provided directory
func RunTrufflehog(directory string) map[string]interface{} {
	log.Println("Running TruffleHog")
	trufflehogResults, err := runTrufflehogOnRepo(directory)
	if err != nil {
		log.Printf("Error while running TruffleHog: %s", err)
		return nil
	}
	categorizedResults := categorizeTrufflehogResults(trufflehogResults, directory)
	newcategorizedResults := utils.ConvertCategorizedResults(categorizedResults)
	return newcategorizedResults
}

// Run TruffleHog on a repository, returns the results as a string
func runTrufflehogOnRepo(directory string) (string, error) {
	// Add the directory to the safe.directory configuration
	err := setGitSafeDirectory(directory)
	if err != nil {
		return "", err
	}
	defer unsetGitSafeDirectory(directory)

	// Run TruffleHog
	cmd := exec.Command("trufflehog3", "--no-entropy", "--format", "JSON", directory)
	output, _ := cmd.CombinedOutput()
	return string(output), nil
}

// Set the git directory as a safe directory
func setGitSafeDirectory(directory string) error {
	cmd := exec.Command("git", "config", "--global", "--add", "safe.directory", directory)
	_, err := cmd.CombinedOutput()
	return err
}

// Unset the git safe directory configuration
func unsetGitSafeDirectory(directory string) error {
	cmd := exec.Command("git", "config", "--global", "--unset", "safe.directory", directory)
	_, err := cmd.CombinedOutput()
	return err
}

// Categorize the TruffleHog results
func categorizeTrufflehogResults(results string, directory string) map[string][]interface{} {
	categorizedResults := utils.InitAdvancedCategorizedResults()

	if results != "" {
		var secrets []Secret
		err := json.Unmarshal([]byte(results), &secrets)
		if err != nil {
			log.Printf("Failed to parse TruffleHog results: %v", err)
			return categorizedResults
		}

		// Format each secret and categorize it
		for _, secret := range secrets {
			// Handle the line field, convert if it's a string
			line := secret.Line
			if lineStr, ok := line.(string); ok {
				// Try to convert the string to an int
				if convertedLine, err := strconv.Atoi(lineStr); err == nil {
					line = convertedLine
				}
			}

			// Add the formatted secret issue to the categorized results
			categorizedResults[utils.SECRET_DETECTION] = append(categorizedResults[utils.SECRET_DETECTION], formatSecretIssue(secret, directory))
		}
	}

	return categorizedResults
}

// Format each secret issue into a map
func formatSecretIssue(secret Secret, directory string) map[string]interface{} {
	return map[string]interface{}{
		"path":     strings.Replace(secret.Path, directory, "", 1),
		"line":     secret.Line,
		"message":  secret.Rule.Message,
		"severity": secret.Rule.Severity,
		"secret":   secret.Secret,
	}
}
