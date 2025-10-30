package tasks

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"armur-codescanner/internal/solidity"
	"armur-codescanner/internal/tools"
	utils "armur-codescanner/pkg"
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

	// For local scans, detect language if not specified but NEVER delete files
	if language == "" {
		language = utils.DetectRepoLanguage(dirPath)
		log.Println("Language detected:", language)
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

		// Run advanced Solidity tools in parallel for better performance
		type toolResult struct {
			results map[string]interface{}
		}

		resultsChan := make(chan toolResult, 4)

		go func() { resultsChan <- toolResult{results: tools.RunSlither(dirPath)} }()
		go func() { resultsChan <- toolResult{results: tools.RunMythril(dirPath)} }()
		go func() { resultsChan <- toolResult{results: tools.RunOyente(dirPath)} }()
		go func() { resultsChan <- toolResult{results: tools.RunSecurify(dirPath)} }()

		// Collect results
		for i := 0; i < 4; i++ {
			res := <-resultsChan
			mergeResults(categorizedResults, res.results)
		}

		// Run SmartCheck separately as it might have different requirements
		smartcheckResults := tools.RunSmartCheck(dirPath)
		mergeResults(categorizedResults, smartcheckResults)

		// Run custom Solidity Semgrep rules
		semgrepSolidityResults := tools.RunSemgrepSolidity(dirPath)
		mergeResults(categorizedResults, semgrepSolidityResults)

		// Check for known vulnerable dependencies
		slitherDepResults := tools.RunSlitherDependencies(dirPath)
		mergeResults(categorizedResults, slitherDepResults)

		// Run gas optimizer
		gasOptResults := tools.RunGasOptimizer(dirPath)
		mergeResults(categorizedResults, gasOptResults)

		// Run LP pairing checks
		lpCheckResults := tools.RunLPPairingChecks(dirPath)
		mergeResults(categorizedResults, lpCheckResults)

		// Run DeFi optimizations
		defiOptResults := tools.RunDeFiOptimizations(dirPath)
		mergeResults(categorizedResults, defiOptResults)
	}
	if cleanup {
		err := os.RemoveAll(dirPath)
		if err != nil {
			return nil, fmt.Errorf("failed to remove directory: %v", err)
		}
	}

	newCatResult := utils.ConvertCategorizedResults(categorizedResults)
	finalresult := utils.ReformatScanResults(newCatResult)
	return map[string]interface{}{
		"complex_functions": finalresult.ComplexFunctions,
		"docstring_absent":  finalresult.DocstringAbsent,
		"antipatterns_bugs": finalresult.AntipatternsBugs,
		"security_issues":   finalresult.SecurityIssues,
	}, nil
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
	case "solidity":
		// Parallel execution for advanced Solidity scanning
		resultsChan := make(chan map[string]interface{}, 5)

		go func() { resultsChan <- tools.RunSlither(dirPath) }()
		go func() { resultsChan <- tools.RunMythril(dirPath) }()
		go func() { resultsChan <- tools.RunOyente(dirPath) }()
		go func() { resultsChan <- tools.RunSecurify(dirPath) }()
		go func() { resultsChan <- tools.RunSmartCheck(dirPath) }()

		for i := 0; i < 5; i++ {
			result := <-resultsChan
			mergeResults(categorizedResults, result)
		}

		// Additional checks
		solcResults := tools.RunSolcCheck(dirPath, "", nil) // AST analysis
		mergeResults(categorizedResults, solcResults)

		semgrepResults := tools.RunSemgrepSolidity(dirPath)
		mergeResults(categorizedResults, semgrepResults)

		slitherDepsResults := tools.RunSlitherDependencies(dirPath)
		mergeResults(categorizedResults, slitherDepsResults)
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

func RunBatchScanTask(taskData map[string]interface{}) map[string]interface{} {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Error while running batch scan: %v", r)
		}
	}()

	// Extract batch scan parameters
	contractPathsInterface, ok := taskData["contract_paths"]
	if !ok {
		return map[string]interface{}{
			"status": "failed",
			"error":  "contract_paths not provided",
		}
	}

	// Type assertion for contract paths
	var contractPaths []string
	if paths, ok := contractPathsInterface.([]string); ok {
		contractPaths = paths
	} else if pathsInterface, ok := contractPathsInterface.([]interface{}); ok {
		for _, p := range pathsInterface {
			if pathStr, ok := p.(string); ok {
				contractPaths = append(contractPaths, pathStr)
			}
		}
	} else {
		return map[string]interface{}{
			"status": "failed",
			"error":  "invalid contract_paths format",
		}
	}

	language, _ := taskData["language"].(string)
	network, _ := taskData["network"].(string)

	log.Printf("Starting batch scan for %d contracts on %s network", len(contractPaths), network)

	// Initialize results aggregation
	batchResults := map[string]interface{}{
		"batch_summary": map[string]interface{}{
			"total_contracts": len(contractPaths),
			"network":         network,
			"language":        language,
			"scan_timestamp":  time.Now().Format(time.RFC3339),
		},
		"contract_results":   []interface{}{},
		"gas_optimizations":  []interface{}{},
		"lp_pairing_checks":  []interface{}{},
		"defi_optimizations": []interface{}{},
		"security_issues":    []interface{}{},
	}

	// Process each contract
	for i, contractPath := range contractPaths {
		log.Printf("Processing contract %d/%d: %s", i+1, len(contractPaths), contractPath)

		// Create temporary directory for this contract
		tempDir, err := os.MkdirTemp("", "batch_contract_*")
		if err != nil {
			log.Printf("Error creating temp dir for contract %s: %v", contractPath, err)
			continue
		}

		// Copy contract to temp directory
		contractFile := filepath.Base(contractPath)
		tempContractPath := filepath.Join(tempDir, contractFile)

		if err := copyFile(contractPath, tempContractPath); err != nil {
			log.Printf("Error copying contract %s: %v", contractPath, err)
			os.RemoveAll(tempDir)
			continue
		}

		// Run gas optimization scan on this contract
		categorizedResults, err := RunSimpleScan(tempDir, language)
		if err != nil {
			log.Printf("Error scanning contract %s: %v", contractPath, err)
			os.RemoveAll(tempDir)
			continue
		}

		// Extract results for this contract
		contractResult := map[string]interface{}{
			"contract_path":    contractPath,
			"contract_name":    contractFile,
			"scan_results":     categorizedResults,
			"gas_savings":      estimateGasSavings(categorizedResults),
			"lp_compatibility": checkLPCompatibility(categorizedResults),
		}

		// Aggregate results
		if gasOpts, ok := categorizedResults["gas_optimizations"]; ok {
			if opts, ok := gasOpts.([]interface{}); ok {
				batchResults["gas_optimizations"] = append(batchResults["gas_optimizations"].([]interface{}), opts...)
			}
		}
		if lpChecks, ok := categorizedResults["lp_pairing_checks"]; ok {
			if checks, ok := lpChecks.([]interface{}); ok {
				batchResults["lp_pairing_checks"] = append(batchResults["lp_pairing_checks"].([]interface{}), checks...)
			}
		}
		if defiOpts, ok := categorizedResults["defi_optimizations"]; ok {
			if opts, ok := defiOpts.([]interface{}); ok {
				batchResults["defi_optimizations"] = append(batchResults["defi_optimizations"].([]interface{}), opts...)
			}
		}
		if secIssues, ok := categorizedResults["security_issues"]; ok {
			if issues, ok := secIssues.([]interface{}); ok {
				batchResults["security_issues"] = append(batchResults["security_issues"].([]interface{}), issues...)
			}
		}

		batchResults["contract_results"] = append(batchResults["contract_results"].([]interface{}), contractResult)

		// Clean up
		os.RemoveAll(tempDir)
	}

	// Calculate batch summary
	totalGasSavings := calculateTotalGasSavings(batchResults)
	if summary, ok := batchResults["batch_summary"].(map[string]interface{}); ok {
		summary["total_estimated_gas_savings"] = totalGasSavings
		summary["contracts_processed"] = len(batchResults["contract_results"].([]interface{}))
	}

	log.Printf("Batch scan completed. Processed %d contracts with estimated gas savings: %s",
		len(batchResults["contract_results"].([]interface{})), totalGasSavings)

	return batchResults
}

func copyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = destFile.ReadFrom(sourceFile)
	return err
}

func estimateGasSavings(results map[string]interface{}) string {
	totalSavings := 0

	if gasOpts, ok := results["gas_optimizations"]; ok {
		if opts, ok := gasOpts.([]interface{}); ok {
			for _, opt := range opts {
				if optMap, ok := opt.(map[string]interface{}); ok {
					if savings, ok := optMap["estimated_savings"]; ok {
						if savingsStr, ok := savings.(string); ok {
							// Simple parsing of savings estimates
							if strings.Contains(savingsStr, "200-500 gas") {
								totalSavings += 350
							} else if strings.Contains(savingsStr, "5000+ gas") {
								totalSavings += 5000
							} else if strings.Contains(savingsStr, "21000+ gas") {
								totalSavings += 21000
							}
						}
					}
				}
			}
		}
	}

	if totalSavings > 0 {
		return fmt.Sprintf("%d gas per transaction", totalSavings)
	}
	return "Variable"
}

func checkLPCompatibility(results map[string]interface{}) string {
	if lpChecks, ok := results["lp_pairing_checks"]; ok {
		if checks, ok := lpChecks.([]interface{}); ok && len(checks) > 0 {
			return "Compatible with Uniswap V2/V3 on Polygon/Amoy"
		}
	}
	return "Not analyzed"
}

func calculateTotalGasSavings(batchResults map[string]interface{}) string {
	totalSavings := 0
	contractResults, ok := batchResults["contract_results"].([]interface{})
	if !ok {
		return "Unknown"
	}

	for _, result := range contractResults {
		if resultMap, ok := result.(map[string]interface{}); ok {
			if scanResults, ok := resultMap["scan_results"].(map[string]interface{}); ok {
				if gasOpts, ok := scanResults["gas_optimizations"]; ok {
					if opts, ok := gasOpts.([]interface{}); ok {
						totalSavings += len(opts) * 1000 // Rough estimate per optimization
					}
				}
			}
		}
	}

	if totalSavings > 0 {
		return fmt.Sprintf("%d+ gas per optimized contract", totalSavings)
	}
	return "Variable"
}
