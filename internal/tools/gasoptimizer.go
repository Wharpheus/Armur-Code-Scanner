package tools

import (
	utils "armur-codescanner/pkg"
	"bufio"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

const (
	GAS_OPTIMIZATIONS  = "gas_optimizations"
	LP_PAIRING_CHECKS  = "lp_pairing_checks"
	DEFI_OPTIMIZATIONS = "defi_optimizations"
)

type GasOptimization struct {
	Path           string `json:"path"`
	Line           int    `json:"line"`
	Message        string `json:"message"`
	Severity       string `json:"severity"`
	Savings        string `json:"estimated_savings"`
	Category       string `json:"category"`
	Recommendation string `json:"recommendation"`
}

type LPPairingCheck struct {
	Path          string `json:"path"`
	Line          int    `json:"line"`
	Message       string `json:"message"`
	Severity      string `json:"severity"`
	Protocol      string `json:"protocol"`
	Compatibility string `json:"compatibility"`
}

func RunGasOptimizer(directory string) map[string]interface{} {
	log.Println("Running Gas Optimizer...")
	results := analyzeSolidityFiles(directory)
	categorizedResults := categorizeGasOptimizations(results)
	return utils.ConvertCategorizedResults(categorizedResults)
}

func RunLPPairingChecks(directory string) map[string]interface{} {
	log.Println("Running LP Pairing Checks...")
	results := analyzeLPPairing(directory)
	categorizedResults := categorizeLPPairingResults(results)
	return utils.ConvertCategorizedResults(categorizedResults)
}

func RunDeFiOptimizations(directory string) map[string]interface{} {
	log.Println("Running DeFi Optimizations...")
	results := analyzeDeFiOptimizations(directory)
	categorizedResults := categorizeDeFiResults(results)
	return utils.ConvertCategorizedResults(categorizedResults)
}

func analyzeSolidityFiles(directory string) []GasOptimization {
	var optimizations []GasOptimization

	err := filepath.Walk(directory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !strings.HasSuffix(strings.ToLower(info.Name()), ".sol") {
			return nil
		}

		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()

		relPath, _ := filepath.Rel(directory, path)
		optimizations = append(optimizations, analyzeFile(file, relPath)...)

		return nil
	})

	if err != nil {
		log.Printf("Error walking directory: %v", err)
	}

	return optimizations
}

func analyzeFile(file *os.File, relPath string) []GasOptimization {
	var optimizations []GasOptimization
	scanner := bufio.NewScanner(file)
	lineNum := 0
	var lines []string

	// Read all lines for context
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
		lineNum++
	}

	// Analyze each line
	for i, line := range lines {
		lineNum := i + 1
		opts := analyzeLine(line, lineNum, lines, relPath)
		optimizations = append(optimizations, opts...)
	}

	return optimizations
}

func analyzeLine(line string, lineNum int, allLines []string, relPath string) []GasOptimization {
	var optimizations []GasOptimization

	// Storage optimization patterns
	if strings.Contains(line, "uint256") && !strings.Contains(line, "constant") && !strings.Contains(line, "immutable") {
		// Check for packing opportunities
		if canPackVariables(allLines, lineNum) {
			optimizations = append(optimizations, GasOptimization{
				Path:           relPath,
				Line:           lineNum,
				Message:        "Consider packing uint256 variables with smaller types for gas savings",
				Severity:       "MEDIUM",
				Savings:        "200-500 gas per slot",
				Category:       "storage_packing",
				Recommendation: "Pack multiple small variables into single storage slots",
			})
		}
	}

	// External calls in loops
	if strings.Contains(line, "for") || strings.Contains(line, "while") {
		for j := lineNum; j < len(allLines) && j < lineNum+10; j++ {
			if strings.Contains(allLines[j], ".call(") || strings.Contains(allLines[j], ".transfer(") || strings.Contains(allLines[j], ".send(") {
				optimizations = append(optimizations, GasOptimization{
					Path:           relPath,
					Line:           lineNum,
					Message:        "External calls inside loops are expensive",
					Severity:       "HIGH",
					Savings:        "21000+ gas per iteration",
					Category:       "external_calls",
					Recommendation: "Cache external call results or use batching",
				})
				break
			}
		}
	}

	// Memory vs Storage usage
	if strings.Contains(line, "memory") && (strings.Contains(line, "[]") || strings.Contains(line, " mapping")) {
		optimizations = append(optimizations, GasOptimization{
			Path:           relPath,
			Line:           lineNum,
			Message:        "Large data structures in memory can be expensive",
			Severity:       "LOW",
			Savings:        "Variable",
			Category:       "memory_usage",
			Recommendation: "Consider storage for frequently accessed data",
		})
	}

	// Function visibility optimization
	if strings.Contains(line, "function ") && strings.Contains(line, " public ") {
		optimizations = append(optimizations, GasOptimization{
			Path:           relPath,
			Line:           lineNum,
			Message:        "Public functions are more expensive than external",
			Severity:       "LOW",
			Savings:        "50-100 gas",
			Category:       "function_visibility",
			Recommendation: "Use external instead of public when possible",
		})
	}

	// State variable reads in loops
	if strings.Contains(line, "for") || strings.Contains(line, "while") {
		for j := lineNum; j < len(allLines) && j < lineNum+10; j++ {
			if regexp.MustCompile(`\b\w+\s*=`).MatchString(allLines[j]) {
				optimizations = append(optimizations, GasOptimization{
					Path:           relPath,
					Line:           lineNum,
					Message:        "State variable assignments inside loops",
					Severity:       "MEDIUM",
					Savings:        "5000+ gas per assignment",
					Category:       "state_variables",
					Recommendation: "Cache state variables in memory variables",
				})
				break
			}
		}
	}

	return optimizations
}

func canPackVariables(lines []string, startLine int) bool {
	// Simple heuristic: look for multiple uint declarations that could be packed
	uintCount := 0
	for i := startLine - 5; i < startLine+5 && i < len(lines); i++ {
		if i >= 0 && strings.Contains(lines[i], "uint") {
			uintCount++
		}
	}
	return uintCount > 1
}

func analyzeLPPairing(directory string) []LPPairingCheck {
	var checks []LPPairingCheck

	err := filepath.Walk(directory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !strings.HasSuffix(strings.ToLower(info.Name()), ".sol") {
			return nil
		}

		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()

		relPath, _ := filepath.Rel(directory, path)
		checks = append(checks, analyzeLPCompatibility(file, relPath)...)

		return nil
	})

	if err != nil {
		log.Printf("Error walking directory for LP checks: %v", err)
	}

	return checks
}

func analyzeLPCompatibility(file *os.File, relPath string) []LPPairingCheck {
	var checks []LPPairingCheck
	scanner := bufio.NewScanner(file)
	lineNum := 0
	content := ""

	for scanner.Scan() {
		line := scanner.Text()
		content += line + "\n"
		lineNum++

		// Uniswap V2 compatibility checks
		if strings.Contains(line, "IUniswapV2Pair") || strings.Contains(line, "UniswapV2") {
			checks = append(checks, LPPairingCheck{
				Path:          relPath,
				Line:          lineNum,
				Message:       "Uniswap V2 interface detected",
				Severity:      "INFO",
				Protocol:      "Uniswap V2",
				Compatibility: "Compatible with Polygon and Amoy testnet",
			})
		}

		// Uniswap V3 compatibility checks
		if strings.Contains(line, "IUniswapV3Pool") || strings.Contains(line, "UniswapV3") {
			checks = append(checks, LPPairingCheck{
				Path:          relPath,
				Line:          lineNum,
				Message:       "Uniswap V3 interface detected",
				Severity:      "INFO",
				Protocol:      "Uniswap V3",
				Compatibility: "Compatible with Polygon and Amoy testnet",
			})
		}

		// ERC20 checks for LP tokens
		if strings.Contains(line, "IERC20") || strings.Contains(line, "ERC20") {
			checks = append(checks, LPPairingCheck{
				Path:          relPath,
				Line:          lineNum,
				Message:       "ERC20 interface for LP token compatibility",
				Severity:      "INFO",
				Protocol:      "ERC20",
				Compatibility: "Required for liquidity pool tokens",
			})
		}
	}

	// Check for flash loan protection
	if !strings.Contains(content, "reentrancy") && !strings.Contains(content, "nonReentrant") {
		checks = append(checks, LPPairingCheck{
			Path:          relPath,
			Line:          1,
			Message:       "Consider adding reentrancy protection for LP operations",
			Severity:      "MEDIUM",
			Protocol:      "General",
			Compatibility: "Important for DeFi security",
		})
	}

	return checks
}

func analyzeDeFiOptimizations(directory string) []GasOptimization {
	var optimizations []GasOptimization

	err := filepath.Walk(directory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !strings.HasSuffix(strings.ToLower(info.Name()), ".sol") {
			return nil
		}

		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()

		relPath, _ := filepath.Rel(directory, path)
		opts := analyzeDeFiFile(file, relPath)
		optimizations = append(optimizations, opts...)

		return nil
	})

	if err != nil {
		log.Printf("Error walking directory for DeFi analysis: %v", err)
	}

	return optimizations
}

func analyzeDeFiFile(file *os.File, relPath string) []GasOptimization {
	var optimizations []GasOptimization
	scanner := bufio.NewScanner(file)
	lineNum := 0
	content := ""

	for scanner.Scan() {
		line := scanner.Text()
		content += line + "\n"
		lineNum++

		// Flash loan protection
		if strings.Contains(line, "flashLoan") || strings.Contains(line, "flashloan") {
			optimizations = append(optimizations, GasOptimization{
				Path:           relPath,
				Line:           lineNum,
				Message:        "Flash loan operation detected - ensure proper validation",
				Severity:       "HIGH",
				Savings:        "Prevent potential losses",
				Category:       "flash_loan_protection",
				Recommendation: "Add slippage and amount validation",
			})
		}

		// Oracle price feeds
		if strings.Contains(line, "price") && (strings.Contains(line, "oracle") || strings.Contains(line, "feed")) {
			optimizations = append(optimizations, GasOptimization{
				Path:           relPath,
				Line:           lineNum,
				Message:        "Price oracle usage - consider staleness checks",
				Severity:       "MEDIUM",
				Savings:        "Prevent stale price attacks",
				Category:       "oracle_safety",
				Recommendation: "Add timestamp validation for price feeds",
			})
		}

		// LP token operations
		if strings.Contains(line, "addLiquidity") || strings.Contains(line, "removeLiquidity") {
			optimizations = append(optimizations, GasOptimization{
				Path:           relPath,
				Line:           lineNum,
				Message:        "Liquidity operation - optimize for gas efficiency",
				Severity:       "MEDIUM",
				Savings:        "1000-5000 gas",
				Category:       "liquidity_operations",
				Recommendation: "Batch operations and minimize external calls",
			})
		}
	}

	// Check for missing DeFi best practices
	if strings.Contains(content, "transferFrom") && !strings.Contains(content, "safeTransferFrom") {
		optimizations = append(optimizations, GasOptimization{
			Path:           relPath,
			Line:           1,
			Message:        "Consider using SafeERC20 for secure token transfers",
			Severity:       "MEDIUM",
			Savings:        "Prevent transfer failures",
			Category:       "token_transfers",
			Recommendation: "Use OpenZeppelin's SafeERC20 library",
		})
	}

	return optimizations
}

func categorizeGasOptimizations(optimizations []GasOptimization) map[string][]interface{} {
	categorized := utils.InitCategorizedResults()

	for _, opt := range optimizations {
		issue := map[string]interface{}{
			"path":              opt.Path,
			"line":              opt.Line,
			"message":           opt.Message,
			"severity":          opt.Severity,
			"estimated_savings": opt.Savings,
			"category":          opt.Category,
			"recommendation":    opt.Recommendation,
			"tool":              "gas_optimizer",
		}
		categorized[GAS_OPTIMIZATIONS] = append(categorized[GAS_OPTIMIZATIONS], issue)
	}

	return categorized
}

func categorizeLPPairingResults(checks []LPPairingCheck) map[string][]interface{} {
	categorized := utils.InitCategorizedResults()

	for _, check := range checks {
		issue := map[string]interface{}{
			"path":          check.Path,
			"line":          check.Line,
			"message":       check.Message,
			"severity":      check.Severity,
			"protocol":      check.Protocol,
			"compatibility": check.Compatibility,
			"tool":          "lp_pairing_checker",
		}
		categorized[LP_PAIRING_CHECKS] = append(categorized[LP_PAIRING_CHECKS], issue)
	}

	return categorized
}

func categorizeDeFiResults(optimizations []GasOptimization) map[string][]interface{} {
	categorized := utils.InitCategorizedResults()

	for _, opt := range optimizations {
		issue := map[string]interface{}{
			"path":              opt.Path,
			"line":              opt.Line,
			"message":           opt.Message,
			"severity":          opt.Severity,
			"estimated_savings": opt.Savings,
			"category":          opt.Category,
			"recommendation":    opt.Recommendation,
			"tool":              "defi_optimizer",
		}
		categorized[DEFI_OPTIMIZATIONS] = append(categorized[DEFI_OPTIMIZATIONS], issue)
	}

	return categorized
}
