package utils

import (
	pkg "armur-codescanner/pkg/common"
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/jung-kurt/gofpdf"
)

type Issue struct {
	Path    string `json:"path"`
	Line    int    `json:"line"`
	Column  int    `json:"column,omitempty"`
	Message string `json:"message"`
	RuleID  string `json:"ruleId,omitempty"`
}

type SimpleScanResults struct {
	ComplexFunctions []Issue `json:"complex_functions"`
	DocstringAbsent  []Issue `json:"docstring_absent"`
	AntipatternsBugs []Issue `json:"antipatterns_bugs"`
	SecurityIssues   []Issue `json:"security_issues"`
}

type CWEData struct {
	CWE                 string            `json:"cwe"`
	GoodPracticeExample map[string]string `json:"good_practice_example"`
	BadPracticeExample  map[string]string `json:"bad_practice_example"`
}

const (
	SimpleScan   = "simple_scan"
	AdvancedScan = "advanced_scan"
	FileScan     = "file_scan"
	LocalScan    = "local_scan"
)

// Constants
const (
	DEAD_CODE          = "dead_code"
	DUPLICATE_CODE     = "duplicate_code"
	SECRET_DETECTION   = "secret_detection"
	INFRA_SECURITY     = "infra_security"
	SCA                = "sca"
	COMPLEX_FUNCTIONS  = "complex_functions"
	DOCKSTRING_ABSENT  = "dockstring_absent"
	ANTIPATTERNS_BUGS  = "antipatterns_bugs"
	SECURITY_ISSUES    = "security_issues"
	UNKNOWN            = "unknown"
	GAS_ISSUES         = "gas_issues"
	GAS_OPTIMIZATIONS  = "gas_optimizations"
	LP_PAIRING_CHECKS  = "lp_pairing_checks"
	DEFI_OPTIMIZATIONS = "defi_optimizations"
	BUILD_ISSUES       = "build_issues"
)

// ReformatScanResults reformats simple scan results
func ReformatScanResults(results map[string]interface{}) *SimpleScanResults {
	reformattedResults := &SimpleScanResults{
		ComplexFunctions: []Issue{},
		DocstringAbsent:  []Issue{},
		AntipatternsBugs: []Issue{},
		SecurityIssues:   []Issue{},
	}

	if data, ok := results[COMPLEX_FUNCTIONS]; ok {
		if issues, ok := data.([]interface{}); ok {
			for _, item := range issues {
				if issueMap, ok := item.(map[string]interface{}); ok {
					issue := Issue{
						Path:    getString(issueMap, "path"),
						Line:    getInt(issueMap, "line"),
						Column:  getInt(issueMap, "column"),
						Message: getString(issueMap, "message"),
						RuleID:  getString(issueMap, "ruleId"),
					}
					reformattedResults.ComplexFunctions = append(reformattedResults.ComplexFunctions, issue)
				}
			}
		}
	}
	if data, ok := results[DOCKSTRING_ABSENT]; ok {
		if issues, ok := data.([]interface{}); ok {
			for _, item := range issues {
				if issueMap, ok := item.(map[string]interface{}); ok {
					issue := Issue{
						Path:    getString(issueMap, "path"),
						Line:    getInt(issueMap, "line"),
						Column:  getInt(issueMap, "column"),
						Message: getString(issueMap, "message"),
						RuleID:  getString(issueMap, "ruleId"),
					}
					reformattedResults.DocstringAbsent = append(reformattedResults.DocstringAbsent, issue)
				}
			}
		}
	}
	if data, ok := results[ANTIPATTERNS_BUGS]; ok {
		if issues, ok := data.([]interface{}); ok {
			for _, item := range issues {
				if issueMap, ok := item.(map[string]interface{}); ok {
					issue := Issue{
						Path:    getString(issueMap, "path"),
						Line:    getInt(issueMap, "line"),
						Column:  getInt(issueMap, "column"),
						Message: getString(issueMap, "message"),
						RuleID:  getString(issueMap, "ruleId"),
					}
					reformattedResults.AntipatternsBugs = append(reformattedResults.AntipatternsBugs, issue)
				}
			}
		}
	}
	if data, ok := results[SECURITY_ISSUES]; ok {
		if issues, ok := data.([]interface{}); ok {
			for _, item := range issues {
				if issueMap, ok := item.(map[string]interface{}); ok {
					issue := Issue{
						Path:    getString(issueMap, "path"),
						Line:    getInt(issueMap, "line"),
						Column:  getInt(issueMap, "column"),
						Message: getString(issueMap, "message"),
						RuleID:  getString(issueMap, "ruleId"),
					}
					reformattedResults.SecurityIssues = append(reformattedResults.SecurityIssues, issue)
				}
			}
		}
	}

	return reformattedResults
}

func ReformatAdvancedScanResults(results map[string]interface{}) map[string]interface{} {
	reformattedResults := map[string]interface{}{
		DEAD_CODE:        []interface{}{},
		DUPLICATE_CODE:   []interface{}{},
		SECRET_DETECTION: []interface{}{},
		INFRA_SECURITY:   []interface{}{},
		SCA:              []interface{}{},
	}

	if data := ReformatDeadCode(results); data != nil {
		reformattedResults[DEAD_CODE] = data
	}
	if data := ReformatDuplicateCode(results); data != nil {
		reformattedResults[DUPLICATE_CODE] = data
	}
	if data := ReformatSecretDetection(results); data != nil {
		reformattedResults[SECRET_DETECTION] = data
	}
	if data := ReformatInfraSecurity(results); data != nil {
		reformattedResults[INFRA_SECURITY] = data
	}
	if data := ReformatSCAIssues(results); data != nil {
		reformattedResults[SCA] = data
	}

	return reformattedResults
}

func ReformatDeadCode(results map[string]interface{}) []map[string]interface{} {
	deadCodeGroupedIssues := make(map[string][]map[string]interface{})
	deadCodeResults, ok := results[DEAD_CODE]
	if !ok {
		return []map[string]interface{}{}
	}
	deadCodeSlice, ok := deadCodeResults.([]interface{})
	if !ok || len(deadCodeSlice) == 0 {
		return []map[string]interface{}{}
	}
	for _, issue := range deadCodeSlice {
		issueMap, ok := issue.(map[string]interface{})
		if !ok {
			continue
		}

		checkID, ok := issueMap["check_id"].(string)
		if !ok {
			continue
		}

		deadCodeGroupedIssues[checkID] = append(deadCodeGroupedIssues[checkID], issueMap)
	}

	// Prepare the result array
	var result []map[string]interface{}
	for checkID, issues := range deadCodeGroupedIssues {
		result = append(result, map[string]interface{}{
			"check_id": checkID,
			"issues":   issues,
		})
	}

	return result
}

func ReformatVultureOutput(vultureResults string, results map[string]interface{}) {
	// If the DEAD_CODE category doesn't exist, initialize it with an empty slice
	if _, ok := results[DEAD_CODE]; !ok {
		results[DEAD_CODE] = []interface{}{}
	}

	// Split the Vulture results by new lines
	lines := strings.Split(vultureResults, "\n")

	// Loop through each line to extract information
	for _, line := range lines {
		// Skip empty lines
		if len(line) == 0 {
			continue
		}

		// Parse the line into file, line, and message
		parts := strings.Split(line, " line:")
		if len(parts) < 2 {
			continue
		}

		fileAndMessage := parts[0]
		message := parts[1]

		// Split the file and message part
		fileParts := strings.Split(fileAndMessage, " message:")
		if len(fileParts) < 2 {
			continue
		}

		file := strings.TrimSpace(fileParts[0])
		msg := strings.TrimSpace(fileParts[1])

		// Create the issue map
		issue := map[string]interface{}{
			"file":    file,
			"line":    strings.TrimSpace(message),
			"message": msg,
		}

		// Append to the DEAD_CODE category
		results[DEAD_CODE] = append(results[DEAD_CODE].([]interface{}), issue)
	}
}

// ReformatDuplicateCode reformats duplicate code results
func ReformatDuplicateCode(results map[string]interface{}) interface{} {
	return results[DUPLICATE_CODE]
}

// ReformatSecretDetection reformats secret detection results
func ReformatSecretDetection(results map[string]interface{}) interface{} {
	return results[SECRET_DETECTION]
}

func ReformatInfraSecurity(results map[string]interface{}) []map[string]interface{} {
	if results[INFRA_SECURITY] == nil {
		return []map[string]interface{}{}
	}

	infraSecurityGroupedIssues := make(map[string][]interface{})

	for _, issue := range results[INFRA_SECURITY].([]interface{}) {
		issueMap := issue.(map[string]interface{})
		message := issueMap["message"].(string)
		infraSecurityGroupedIssues[message] = append(infraSecurityGroupedIssues[message], issue)
	}

	var result []map[string]interface{}
	for message, issues := range infraSecurityGroupedIssues {
		result = append(result, map[string]interface{}{
			"message": message,
			"issues":  issues,
		})
	}

	return result
}

func LoadCWEData(filePath string) ([]CWEData, error) {
	// Read the file
	file, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("error reading file: %v", err)
	}

	// Unmarshal the JSON into a slice of CWEData
	var data []CWEData
	err = json.Unmarshal(file, &data)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling CWE data: %v", err)
	}

	return data, nil
}

func GetPracticesFromJSON(data []CWEData, language string, cwe string) map[string]string {
	for _, item := range data {
		if item.CWE == cwe {
			goodPractice, goodFound := item.GoodPracticeExample[language]
			badPractice, badFound := item.BadPracticeExample[language]

			if goodFound && badFound {
				return map[string]string{
					"good_practice": goodPractice,
					"bad_practice":  badPractice,
				}
			}
		}
	}
	return map[string]string{}
}

func ReformatComplexFunctions(results map[string]interface{}) []map[string]interface{} {
	complexGroupedIssues := make(map[string][]interface{})

	issues, ok := results[COMPLEX_FUNCTIONS]
	if !ok || issues == nil {
		log.Printf("No COMPLEX_FUNCTIONS key or it is nil")
		return []map[string]interface{}{}
	}

	issueList, ok := issues.([]interface{})
	if !ok {
		log.Printf("COMPLEX_FUNCTIONS is not a []interface{}")
		return []map[string]interface{}{}
	}

	for _, issue := range issueList {
		issueMap := issue.(map[string]interface{})
		path := issueMap["path"].(string)
		complexGroupedIssues[path] = append(complexGroupedIssues[path], issue)
	}

	var result []map[string]interface{}
	for path, issues := range complexGroupedIssues {
		result = append(result, map[string]interface{}{
			"path":   path,
			"issues": issues,
		})
	}
	return result
}

// ReformatDocstringIssues reformats docstring issues results
func ReformatDocstringIssues(results map[string]interface{}) []map[string]interface{} {
	docstringGroupedIssues := make(map[string][]interface{})

	docstringIssues, ok := results[DOCKSTRING_ABSENT]
	if !ok || docstringIssues == nil {
		return []map[string]interface{}{}
	}

	issueList, ok := docstringIssues.([]interface{})
	if !ok {
		return []map[string]interface{}{}
	}

	for _, issue := range issueList {
		issueMap, ok := issue.(map[string]interface{})
		if !ok {
			continue
		}
		path, ok := issueMap["path"].(string)
		if !ok {
			continue
		}
		docstringGroupedIssues[path] = append(docstringGroupedIssues[path], issue)
	}

	var result []map[string]interface{}
	for path, issues := range docstringGroupedIssues {
		result = append(result, map[string]interface{}{
			"path":   path,
			"issues": issues,
		})
	}
	return result
}

// ReformatSecurityIssues reformats security issues results
func ReformatSecurityIssues(results map[string]interface{}) []map[string]interface{} {
	securityGroupedIssues := make(map[string]map[string][]interface{})

	if securityIssues, ok := results[SECURITY_ISSUES].([]interface{}); ok {
		for _, issue := range securityIssues {
			issueMap, ok := issue.(map[string]interface{})
			if !ok {
				fmt.Println("Skipping invalid issue format")
				continue
			}

			cweKey := UNKNOWN

			if cweArray, ok := issueMap["cwe"].([]string); ok && len(cweArray) > 0 {
				cweKey = cweArray[0]
			} else if cweInterfaceArray, ok := issueMap["cwe"].([]interface{}); ok && len(cweInterfaceArray) > 0 {
				if firstCWE, ok := cweInterfaceArray[0].(string); ok {
					cweKey = firstCWE
				} else {
					fmt.Println("Invalid CWE format in array:", cweInterfaceArray)
				}
			} else {
				fmt.Println("CWE not found or invalid type:", issueMap["cwe"])
			}

			path, _ := issueMap["path"].(string)

			if _, exists := securityGroupedIssues[cweKey]; !exists {
				securityGroupedIssues[cweKey] = make(map[string][]interface{})
			}

			securityGroupedIssues[cweKey][path] = append(securityGroupedIssues[cweKey][path], issue)
		}
	} else {
		fmt.Println("No 'security_issues' found or invalid format")
	}

	// Convert the grouped map into the desired structure
	var resultCwePathGrouping []map[string]interface{}
	for cwe, files := range securityGroupedIssues {
		var fileIssues []map[string]interface{}
		for path, issues := range files {
			fileIssues = append(fileIssues, map[string]interface{}{
				"path":   path,
				"issues": issues,
			})
		}

		resultCwePathGrouping = append(resultCwePathGrouping, map[string]interface{}{
			"cwe":   cwe,
			"files": fileIssues,
		})
	}

	return resultCwePathGrouping
}

// ReformatAntipatternsBugs reformats antipatterns bugs results
func ReformatAntipatternsBugs(results map[string]interface{}) []map[string]interface{} {
	antipatternGroupedIssues := make(map[string]map[string][]interface{})

	antipatternIssues, ok := results[ANTIPATTERNS_BUGS]
	if !ok || antipatternIssues == nil {
		return []map[string]interface{}{}
	}

	issueList, ok := antipatternIssues.([]interface{})
	if !ok {
		return []map[string]interface{}{}
	}

	for _, issue := range issueList {
		issueMap, ok := issue.(map[string]interface{})
		if !ok {
			continue
		}
		messageKey := UNKNOWN
		if message, ok := issueMap["message"].(string); ok {
			messageKey = message
		}

		path, ok := issueMap["path"].(string)
		if !ok {
			continue
		}
		if _, exists := antipatternGroupedIssues[messageKey]; !exists {
			antipatternGroupedIssues[messageKey] = make(map[string][]interface{})
		}
		antipatternGroupedIssues[messageKey][path] = append(antipatternGroupedIssues[messageKey][path], issue)
	}

	var resultMessagePathGrouping []map[string]interface{}
	for message, files := range antipatternGroupedIssues {
		var fileList []map[string]interface{}
		for path, issues := range files {
			fileList = append(fileList, map[string]interface{}{
				"path":   path,
				"issues": issues,
			})
		}
		resultMessagePathGrouping = append(resultMessagePathGrouping, map[string]interface{}{
			"message":       message,
			"seen_in_files": len(files),
			"files":         fileList,
		})
	}
	return resultMessagePathGrouping
}

// ReformatSCAIssues reformats SCA issues results
func ReformatSCAIssues(results map[string]interface{}) []map[string]interface{} {
	if results[SCA] == nil {
		return []map[string]interface{}{}
	}

	scaGroupedIssues := make(map[string][]interface{})

	for _, issue := range results[SCA].([]interface{}) {
		issueMap := issue.(map[string]interface{})
		path := issueMap["path"].(string)
		scaGroupedIssues[path] = append(scaGroupedIssues[path], issue)
	}

	var result []map[string]interface{}
	for path, issues := range scaGroupedIssues {
		result = append(result, map[string]interface{}{
			"path":   path,
			"issues": issues,
		})
	}
	return result
}

// CloneRepo clones a repository to a temporary directory
func CloneRepo(repositoryURL string) (string, error) {
	baseDir := "/tmp/armur/repos"
	if err := os.MkdirAll(baseDir, 0755); err != nil {
		return "", fmt.Errorf("error creating base directory: %w", err)
	}

	tempDir, err := os.MkdirTemp(baseDir, "repo")
	if err != nil {
		return "", fmt.Errorf("error creating temp directory: %w", err)
	}

	_, err = git.PlainClone(tempDir, false, &git.CloneOptions{
		URL:   repositoryURL,
		Depth: 1,
		// Disable progress to avoid noisy logs in server mode
		Tags: git.NoTags,
	})
	if err != nil {
		return "", fmt.Errorf("error cloning repository: %w", err)
	}
	return tempDir, nil
}

// DetectRepoLanguage detects the language of a repository
func DetectRepoLanguage(directory string) string {
	languages := map[string]int{"go": 0, "py": 0, "js": 0, "solidity": 0}

	skipDirs := map[string]struct{}{
		".git": {}, "node_modules": {}, "vendor": {}, "venv": {}, "__pycache__": {},
		"dist": {}, "build": {}, ".next": {}, ".cache": {},
	}

	_ = filepath.Walk(directory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			if _, skip := skipDirs[info.Name()]; skip {
				return filepath.SkipDir
			}
			return nil
		}
		name := strings.ToLower(info.Name())
		switch {
		case strings.HasSuffix(name, ".go"):
			languages["go"]++
		case strings.HasSuffix(name, ".py"):
			languages["py"]++
		case strings.HasSuffix(name, ".js"):
			languages["js"]++
		case strings.HasSuffix(name, ".sol"):
			languages["solidity"]++
		}
		return nil
	})

	maxLang := ""
	maxCount := 0
	for lang, count := range languages {
		if count > maxCount {
			maxLang = lang
			maxCount = count
		}
	}

	return strings.ToLower(maxLang)
}

// DetectFileLanguage detects the language of a file
func DetectFileLanguage(file string) string {
	name := strings.ToLower(file)
	switch {
	case strings.HasSuffix(name, ".go"):
		return "go"
	case strings.HasSuffix(name, ".py"):
		return "py"
	case strings.HasSuffix(name, ".js"):
		return "js"
	case strings.HasSuffix(name, ".sol"):
		return "solidity"
	default:
		return ""
	}
}

func RemoveNonRelevantFiles(dirPath string, language string) error {
	// Get extensions for the specified language
	extensions, ok := pkg.LanguageFileExtensions[strings.ToLower(language)]
	if !ok {
		extensions = []string{}
	}

	skipDirs := map[string]struct{}{
		".git": {}, "node_modules": {}, "vendor": {}, "venv": {}, "__pycache__": {},
		"dist": {}, "build": {}, ".next": {}, ".cache": {},
	}

	return filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			if _, skip := skipDirs[info.Name()]; skip {
				return filepath.SkipDir
			}
			return nil
		}

		name := strings.ToLower(info.Name())
		shouldKeep := false
		for _, ext := range extensions {
			if strings.HasSuffix(name, ext) {
				shouldKeep = true
				break
			}
		}
		if !shouldKeep {
			if err := os.Remove(path); err != nil {
				return err
			}
		}
		return nil
	})
}

// InitCategorizedResults initializes categorized results
func InitCategorizedResults() map[string][]interface{} {
	return map[string][]interface{}{
		DOCKSTRING_ABSENT: {},
		SECURITY_ISSUES:   {},
		COMPLEX_FUNCTIONS: {},
		ANTIPATTERNS_BUGS: {},
		GAS_ISSUES:        {},
	}
}

func ConvertCategorizedResults(input map[string][]interface{}) map[string]interface{} {
	converted := make(map[string]interface{})
	for key, value := range input {
		if value == nil {
			converted[key] = []interface{}{}
		} else {
			converted[key] = value
		}
	}
	return converted
}

// InitAdvancedCategorizedResults initializes advanced categorized results
func InitAdvancedCategorizedResults() map[string][]interface{} {
	return map[string][]interface{}{
		DEAD_CODE:        {},
		DUPLICATE_CODE:   {},
		SECRET_DETECTION: {},
		INFRA_SECURITY:   {},
		SCA:              {},
	}
}

type ReportItem struct {
	Column               int    `json:"column"`
	Confidence           string `json:"confidence"`
	File                 string `json:"file"`
	Line                 int    `json:"line"`
	Message              string `json:"message"`
	Owasp                string `json:"owasp"`
	Severity             string `json:"severity"`
	SuggestedRemediation string `json:"suggested_remediation"`
}

func GenerateOwaspReport(taskResult interface{}) ([]ReportItem, error) {
	var owaspReport []ReportItem

	// Ensure taskResult is a map
	taskResultMap, ok := taskResult.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("taskResult is not a valid map")
	}

	// Get SECURITY_ISSUES from the map
	securityIssues, ok := taskResultMap[SECURITY_ISSUES].([]interface{})
	if !ok {
		return nil, fmt.Errorf("SECURITY_ISSUES key not found or invalid")
	}

	// Iterate over SECURITY_ISSUES
	for _, rawIssue := range securityIssues {
		issueMap, ok := rawIssue.(map[string]interface{})
		if !ok {
			continue
		}

		// Iterate over "files" in each issue
		files, ok := issueMap["files"].([]interface{})
		if !ok {
			continue
		}
		for _, rawFile := range files {
			fileInfo, ok := rawFile.(map[string]interface{})
			if !ok {
				continue
			}

			// Iterate over "issues" in each file
			fileIssues, ok := fileInfo["issues"].([]interface{})
			if !ok {
				continue
			}
			for _, rawFileIssue := range fileIssues {
				fileIssue, ok := rawFileIssue.(map[string]interface{})
				if !ok {
					continue
				}

				// Extract OWASP entries
				owaspItems, ok := fileIssue["owasp"].([]interface{})
				if !ok {
					continue
				}
				for _, rawOwasp := range owaspItems {
					owaspItem, ok := rawOwasp.(string)
					if !ok {
						continue
					}

					// Construct the report item
					reportItem := ReportItem{
						Owasp:      owaspItem,
						File:       getString(fileIssue, "path"),
						Line:       getInt(fileIssue, "line"),
						Column:     getInt(fileIssue, "column"),
						Message:    getString(fileIssue, "message"),
						Confidence: getString(fileIssue, "confidence"),
						Severity:   getString(fileIssue, "severity"),
						SuggestedRemediation: fmt.Sprintf(
							"Bad Practice: %s\nSuggested Fix: %s",
							getString(fileIssue, "bad_practice"),
							getString(fileIssue, "good_practice"),
						),
					}

					// Append to the OWASP report
					owaspReport = append(owaspReport, reportItem)
				}
			}
		}
	}

	return owaspReport, nil
}

// Helper function to safely get a string value from a map
func getString(data map[string]interface{}, key string) string {
	if val, ok := data[key]; ok {
		if strVal, ok := val.(string); ok {
			return strVal
		}
	}
	return ""
}

// Helper function to safely get an int value from a map
func getInt(data map[string]interface{}, key string) int {
	if val, ok := data[key]; ok {
		if floatVal, ok := val.(float64); ok {
			return int(floatVal)
		}
	}
	return 0
}

type SANSReportItem struct {
	CWE                  string `json:"cwe"`
	File                 string `json:"file"`
	Line                 int    `json:"line"`
	Column               int    `json:"column"`
	Message              string `json:"message"`
	Confidence           string `json:"confidence"`
	Severity             string `json:"severity"`
	SuggestedRemediation string `json:"suggested_remediation"`
}

func GenerateSANSReports(taskResult interface{}) ([]SANSReportItem, error) {
	var sansReport []SANSReportItem

	// Ensure taskResult is a map
	taskResultMap, ok := taskResult.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("taskResult is not a valid map")
	}

	// Get SECURITY_ISSUES from the map
	securityIssues, ok := taskResultMap[SECURITY_ISSUES].([]interface{})
	if !ok {
		return nil, fmt.Errorf("SECURITY_ISSUES key not found or invalid")
	}

	// Iterate over SECURITY_ISSUES
	for _, rawIssue := range securityIssues {
		issueMap, ok := rawIssue.(map[string]interface{})
		if !ok {
			continue
		}

		// Get CWE for this issue
		cwe, _ := issueMap["cwe"].(string)

		// Iterate over "files" in each issue
		files, ok := issueMap["files"].([]interface{})
		if !ok {
			continue
		}
		for _, rawFile := range files {
			fileInfo, ok := rawFile.(map[string]interface{})
			if !ok {
				continue
			}

			// Iterate over "issues" in each file
			fileIssues, ok := fileInfo["issues"].([]interface{})
			if !ok {
				continue
			}
			for _, rawFileIssue := range fileIssues {
				fileIssue, ok := rawFileIssue.(map[string]interface{})
				if !ok {
					continue
				}

				// Construct the report item
				reportItem := SANSReportItem{
					CWE:        cwe,
					File:       getString(fileIssue, "path"),
					Line:       getInt(fileIssue, "line"),
					Column:     getInt(fileIssue, "column"),
					Message:    getString(fileIssue, "message"),
					Confidence: getString(fileIssue, "confidence"),
					Severity:   getString(fileIssue, "severity"),
					SuggestedRemediation: fmt.Sprintf(
						"Bad Practice: %s\nSuggested Fix: %s",
						getString(fileIssue, "bad_practice"),
						getString(fileIssue, "good_practice"),
					),
				}

				// Append to the SANS report
				sansReport = append(sansReport, reportItem)
			}
		}
	}

	return sansReport, nil
}

func GeneratePdfReport(taskResult interface{}) ([]byte, error) {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 16)

	// --- Title Page ---
	pdf.SetFillColor(40, 40, 40)
	pdf.Rect(0, 0, 210, 297, "F") // Full page background
	pdf.SetTextColor(255, 255, 255)

	pdf.SetY(100)
	pdf.SetFont("Arial", "B", 32)
	pdf.Cell(0, 10, "Security Assessment Report")
	pdf.Ln(20)

	pdf.SetFont("Arial", "", 16)
	// TODO: Get project name from task data
	pdf.Cell(0, 10, "Project: DNC GEMINI (Placeholder)")
	pdf.Ln(10)
	pdf.Cell(0, 10, fmt.Sprintf("Date: %s", time.Now().Format("02 Jan 2006")))

	// --- Findings Summary Page ---
	pdf.AddPage()
	pdf.SetTextColor(0, 0, 0)
	pdf.SetFont("Arial", "B", 24)
	pdf.Cell(0, 20, "Findings Summary")
	pdf.Ln(25)

	taskResultMap, ok := taskResult.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("taskResult is not a valid map")
	}

	// --- Summary Table ---
	pdf.SetFont("Arial", "B", 12)
	pdf.Cell(40, 10, "Category")
	pdf.Cell(40, 10, "Count")
	pdf.Ln(10)
	pdf.SetFont("Arial", "", 12)

	totalIssues := 0
	categories := []string{SECURITY_ISSUES, GAS_ISSUES, ANTIPATTERNS_BUGS, COMPLEX_FUNCTIONS, DOCKSTRING_ABSENT}

	for _, category := range categories {
		if issues, ok := taskResultMap[category].([]interface{}); ok {
			count := len(issues)
			totalIssues += count
			pdf.Cell(40, 10, category)
			pdf.Cell(40, 10, fmt.Sprintf("%d", count))
			pdf.Ln(10)
		}
	}

	pdf.Ln(10)
	pdf.SetFont("Arial", "B", 14)
	pdf.Cell(40, 10, "Total Issues Found:")
	pdf.Cell(40, 10, fmt.Sprintf("%d", totalIssues))
	pdf.Ln(20)

	// --- Disclaimer ---
	pdf.SetY(250)
	pdf.SetFont("Arial", "I", 8)
	pdf.MultiCell(0, 5, "This is an automated security assessment. The findings in this report are generated by static analysis tools and should be manually verified. This report does not guarantee the absence of all vulnerabilities.", "", "C", false)

	var buf bytes.Buffer
	err := pdf.Output(&buf)
	return buf.Bytes(), err
}

func ReplaceTarget(args []string, old, new string) []string {
	result := make([]string, len(args))
	copy(result, args)
	for i, arg := range result {
		if arg == old {
			result[i] = new
		}
	}
	return result
}
