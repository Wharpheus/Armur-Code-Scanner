package tasks

import (
	tools "armur-codescanner/internal/tools"
	"armur-codescanner/pkg"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

func RunScanTask(repositoryURL, language string) map[string]interface{} {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Error while running scan: %v", r)
		}
	}()

	// Clone the repository
	dirPath, err := utils.CloneRepo(repositoryURL)
	if err != nil {
		log.Println("Error cloning repository:", err)
		return map[string]interface{}{
			"status": "failed",
			"error":  err.Error(),
		}
	}

	if language == "" {
		language = utils.DetectRepoLanguage(dirPath)
		log.Println("Language detected:", language)
	} else {
		// Remove non-relevant files based on the language
		if err := utils.RemoveNonRelevantFiles(dirPath, language); err != nil {
			log.Println("Error removing non-relevant files:", err)
			return map[string]interface{}{
				"status": "failed",
				"error":  err.Error(),
			}
		}
	}

	// Run the scan
	categorizedResults, err_ := RunSimpleScan(dirPath, language)
	if err_ != nil {
		return map[string]interface{}{
			"status": "failed",
		}
	}

	return categorizedResults
}
func AdvancedScanRepositoryTask(repositoryURL, language string) map[string]interface{} {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Error while running scan: %v", r)
		}
	}()

	// Clone the repository
	dirPath, err := utils.CloneRepo(repositoryURL)
	if err != nil {
		log.Println("Error cloning repository:", err)
		return map[string]interface{}{
			"status": "failed",
			"error":  err.Error(),
		}
	}

	// Detect the language if not provided
	if language == "" {
		language = utils.DetectRepoLanguage(dirPath)
		log.Println("Language detected:", language)
	} else {
		// Remove non-relevant files based on the language
		if err := utils.RemoveNonRelevantFiles(dirPath, language); err != nil {
			log.Println("Error removing non-relevant files:", err)
			return map[string]interface{}{
				"status": "failed",
				"error":  err.Error(),
			}
		}
	}

	// Run the advanced scans
	categorizedResults, err := RunAdvancedScans(dirPath, language)
	if err != nil {
		log.Println("Error running advanced scans:", err)
		return map[string]interface{}{
			"status": "failed",
			"error":  err.Error(),
		}
	}
	return categorizedResults
}

func RunSimpleScan(dirPath string, language string) (map[string]interface{}, error) {
	categorizedResults := utils.InitCategorizedResults()
	semgrepResult := tools.RunSemgrep(dirPath, "--config=auto")
	mergeResultss(categorizedResults, semgrepResult)

	// If the language is Go, we run Go-specific tools
	if language == "go" {
		gosecResults := tools.RunGosec(dirPath)
		mergeResultss(categorizedResults, gosecResults)

		// Run Golint
		golintResults := tools.RunGolint(dirPath)
		mergeResultss(categorizedResults, golintResults)

		// Run Govet
		govetResults := tools.RunGovet(dirPath)
		mergeResultss(categorizedResults, govetResults)

		// Run Staticcheck
		staticcheckResults := tools.RunStaticCheck(dirPath)
		mergeResultss(categorizedResults, staticcheckResults)

		// Run Gocyclo
		gocyloResults := tools.RunGocyclo(dirPath)
		mergeResultss(categorizedResults, gocyloResults)
	} else if language == "py" {
		// Run Bandit for Python
		banditResults := tools.RunBandit(dirPath)
		mergeResultss(categorizedResults, banditResults)

		// Run Pydocstyle
		pydocstyleResults := tools.RunPydocstyle(dirPath)
		mergeResultss(categorizedResults, pydocstyleResults)

		// Run Radon
		radonResults := tools.RunRadon(dirPath)
		mergeResultss(categorizedResults, radonResults)

		// Run Pylint
		pylintResults := tools.RunPylint(dirPath)
		mergeResultss(categorizedResults, pylintResults)
	} else if language == "js" {
		eslintResult := tools.RunESLintOnRepo(dirPath)
		mergeResultss(categorizedResults, eslintResult)
	}
	err := os.RemoveAll(dirPath)
	if err != nil {
		return nil, fmt.Errorf("failed to remove directory: %v", err)
	}
	newCatResult := utils.ConvertCategorizedResults(categorizedResults)
	finalresult := utils.ReformatScanResults(newCatResult)
	return finalresult, nil
}

func RunAdvancedScans(dirPath string, language string) (map[string]interface{}, error) {
	// Initialize the categorized results
	categorizedResults := utils.InitAdvancedCategorizedResults()

	// Duplicate code detection
	jscpdResults := tools.RunJSCPD(dirPath)
	mergeResultss(categorizedResults, jscpdResults)

	// Infra security
	checkovResults := tools.RunCheckov(dirPath)
	mergeResultss(categorizedResults, checkovResults)

	// Secret detection
	trufflehogResults := tools.RunTrufflehog(dirPath)
	mergeResultss(categorizedResults, trufflehogResults)

	// Infra security and secret detection
	trivyResults := tools.RunTrivy(dirPath)
	mergeResultss(categorizedResults, trivyResults)

	// SCA
	osvscannerResults, err := tools.RunOSVScanner(dirPath)
	mergeResultss(categorizedResults, osvscannerResults)

	// Dead code detection based on language
	if language == "go" {
		godeadcodeResults := tools.RunGoDeadcode(dirPath)
		mergeResultss(categorizedResults, godeadcodeResults)
	} else if language == "py" {
		vulnResults, _ := tools.RunVulture(dirPath)
		mergeResultss(categorizedResults, vulnResults)
	} else if language == "js" {
		eslintResults := tools.RunESLintAdvanced(dirPath)
		mergeResultss(categorizedResults, eslintResults)
	}
	err = os.RemoveAll(dirPath)
	if err != nil {
		return nil, fmt.Errorf("failed to remove directory: %v", err)
	}
	newCatResult := utils.ConvertCategorizedResults(categorizedResults)
	finalresult := utils.ReformatAdvancedScanResults(newCatResult)
	return finalresult, nil
}
func mergeResultss(categorizedResults map[string][]interface{}, newResults map[string]interface{}) {
	for key, newValue := range newResults {

		if newValue == nil {
			if _, exists := categorizedResults[key]; !exists {
				categorizedResults[key] = []interface{}{}
			}
			continue
		}

		newSlice, ok := newValue.([]interface{})
		if !ok {
			continue
		}

		if len(newSlice) == 0 {
			if _, exists := categorizedResults[key]; !exists {
				categorizedResults[key] = []interface{}{}
			}
			continue
		}

		if existingSlice, exists := categorizedResults[key]; exists {
			categorizedResults[key] = append(existingSlice, newSlice...)
		} else {
			categorizedResults[key] = newSlice
		}
	}
}

func ScanFileTask(filePath string) (map[string]interface{}, error) {
	dirPath := filepath.Dir(filePath)

	defer func() {
		if err := os.RemoveAll(dirPath); err != nil {
			log.Printf("Error cleaning up directory: %v", err)
		} else {
			log.Printf("Successfully cleaned up directory: %s", dirPath)
		}
	}()

	language := utils.DetectFileLanguage(filePath)
	if language == "" {
		err := errors.New("unable to detect file language")
		log.Printf("Error: %v", err)
		return map[string]interface{}{"status": "failed", "error": err.Error()}, err
	}

	categorizedResults, err := RunSimpleScan(filePath, language)
	if err != nil {
		log.Printf("Error while running scans: %v", err)
		return map[string]interface{}{"status": "failed", "error": err.Error()}, err
	}

	return categorizedResults, nil
}
