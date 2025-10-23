package tasks

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"armur-codescanner/internal/tools"
	utils "armur-codescanner/pkg"
	"armur-codescanner/pkg/solidity"
)

func RunScanTask(repositoryURL, language string) map[string]interface{} {
	defer func() {
		if r := recover(); r != nil { // To-do: Create struct for advanced results
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

	language, err = prepareScanDirectory(dirPath, language)
	if err != nil {
		return map[string]interface{}{"status": "failed", "error": err.Error()}
	}
	categorizedResults, err := RunSimpleScan(dirPath, language)
	if err != nil {
		return map[string]interface{}{
			"status": "failed",
			"error":  err.Error(),
		}
	}

	return categorizedResults
}

func RunScanTaskLocal(repoUrl, language string) map[string]interface{} {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Error while running scan: %v", r)
		}
	}()

	dirPath := repoUrl

	language, err := prepareScanDirectory(dirPath, language)
	if err != nil {
		return map[string]interface{}{"status": "failed", "error": err.Error()}
	}
	categorizedResults, err := RunSimpleScanLocal(dirPath, language)
	if err != nil {
		return map[string]interface{}{
			"status": "failed",
			"error":  err.Error(),
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

	language, err = prepareScanDirectory(dirPath, language)
	if err != nil {
		return map[string]interface{}{"status": "failed", "error": err.Error()}
	}
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

func prepareScanDirectory(dirPath, language string) (string, error) {
	if language == "" {
		language = utils.DetectRepoLanguage(dirPath)
		log.Println("Language detected:", language)
	} else {
		// Remove non-relevant files based on the language
		if err := utils.RemoveNonRelevantFiles(dirPath, language); err != nil {
			log.Println("Error removing non-relevant files:", err)
			return "", err
		}
	}
	return language, nil
}

func RunSimpleScan(dirPath string, language string) (map[string]interface{}, error) {
	return runSimpleScanInternal(dirPath, language, true)
}

func runSimpleScanInternal(dirPath string, language string, cleanup bool) (map[string]interface{}, error) {
	categorizedResults := utils.InitCategorizedResults()

	semgrepResult := tools.RunSemgrep(dirPath, "--config=auto")
	mergeResults(categorizedResults, semgrepResult)

	// If the language is Go, we run Go-specific tools
	switch language {
	case "go":
		gosecResults := tools.RunGosec(dirPath)
		mergeResults(categorizedResults, gosecResults)

		// Run Golint
		golintResults := tools.RunGolint(dirPath)
		mergeResults(categorizedResults, golintResults)

		// Run Govet
		govetResults := tools.RunGovet(dirPath)
		mergeResults(categorizedResults, govetResults)

		// Run Staticcheck
		staticcheckResults := tools.RunStaticCheck(dirPath)
		mergeResults(categorizedResults, staticcheckResults)

		// Run Gocyclo
		gocyloResults := tools.RunGocyclo(dirPath)
		mergeResults(categorizedResults, gocyloResults)
	case "py":
		// Run Bandit for Python
		banditResults := tools.RunBandit(dirPath)
		mergeResults(categorizedResults, banditResults)

		// Run Pydocstyle
		pydocstyleResults := tools.RunPydocstyle(dirPath)
		mergeResults(categorizedResults, pydocstyleResults)

		// Run Radon
		radonResults := tools.RunRadon(dirPath)
		mergeResults(categorizedResults, radonResults)

		// Run Pylint
		pylintResults := tools.RunPylint(dirPath)
		mergeResults(categorizedResults, pylintResults)
	case "js":
		eslintResult := tools.RunESLintOnRepo(dirPath)
		mergeResults(categorizedResults, eslintResult)
	case "solidity":
		// Detect project config and run pre-compilation to surface build issues
		conf := solidity.DetectSolidityConfig(dirPath)
		solcDiagnostics := tools.RunSolcCheck(dirPath, conf.Version, conf.Remappings)
		mergeResults(categorizedResults, solcDiagnostics)

		// Run Slither and Mythril with same directory target
		slitherResults := tools.RunSlither(dirPath)
		mergeResults(categorizedResults, slitherResults)

		mythrilResults := tools.RunMythril(dirPath)
		mergeResults(categorizedResults, mythrilResults)
	}
	if cleanup {
		err := os.RemoveAll(dirPath)
		if err != nil {
			return nil, fmt.Errorf("failed to remove directory: %v", err)
		}
	}

	newCatResult := utils.ConvertCategorizedResults(categorizedResults)
	finalresult := utils.ReformatScanResults(newCatResult)
	return finalresult, nil
}

func RunSimpleScanLocal(dirPath string, language string) (map[string]interface{}, error) {
	return runSimpleScanInternal(dirPath, language, false)
}

func RunAdvancedScans(dirPath string, language string) (map[string]interface{}, error) { // To-do: Create struct for advanced results
	// Initialize the categorized results
	categorizedResults := utils.InitAdvancedCategorizedResults()

	// Duplicate code detection
	jscpdResults := tools.RunJSCPD(dirPath)
	mergeResults(categorizedResults, jscpdResults)

	// Infra security
	checkovResults := tools.RunCheckov(dirPath)
	mergeResults(categorizedResults, checkovResults)

	// Secret detection
	trufflehogResults := tools.RunTrufflehog(dirPath)
	mergeResults(categorizedResults, trufflehogResults)

	// Infra security and secret detection
	trivyResults := tools.RunTrivy(dirPath)
	mergeResults(categorizedResults, trivyResults)

	// SCA
	osvscannerResults, err := tools.RunOSVScanner(dirPath)
	if err != nil {
		log.Println("error running OSV Scanner: ", err)
	}
	mergeResults(categorizedResults, osvscannerResults)

	// Dead code detection based on language
	switch language {
	case "go":
		godeadcodeResults := tools.RunGoDeadcode(dirPath)
		mergeResults(categorizedResults, godeadcodeResults)
	case "py":
		vulnResults, _ := tools.RunVulture(dirPath)
		mergeResults(categorizedResults, vulnResults)
	case "js":
		eslintResults := tools.RunESLintAdvanced(dirPath)
		mergeResults(categorizedResults, eslintResults)
	}
	err = os.RemoveAll(dirPath)
	if err != nil {
		return nil, fmt.Errorf("failed to remove directory: %v", err)
	}
	newCatResult := utils.ConvertCategorizedResults(categorizedResults)
	finalresult := utils.ReformatAdvancedScanResults(newCatResult)
	return finalresult, nil
}
func mergeResults(categorizedResults map[string][]interface{}, newResults map[string]interface{}) {
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

	categorizedResults, err := RunSimpleScan(dirPath, language)
	if err != nil {
		log.Printf("Error while running scans: %v", err)
		return map[string]interface{}{"status": "failed", "error": err.Error()}, err
	}

	return categorizedResults, nil
}
