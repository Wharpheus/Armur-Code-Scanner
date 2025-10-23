package tools

import (
	"encoding/json"
	"log"
	"os/exec"
	"strings"

	utils "armur-codescanner/pkg"
)

type TrivyResult struct {
	Results []struct {
		Target          string `json:"Target"`
		Class           string `json:"Class"`
		Type            string `json:"Type"`
		Vulnerabilities []struct {
			VulnerabilityID string `json:"VulnerabilityID"`
			PkgID           string `json:"PkgID"`
			PkgName         string `json:"PkgName"`
			PkgIdentifier   struct {
				PURL string `json:"PURL"`
				UID  string `json:"UID"`
			} `json:"PkgIdentifier"`
			InstalledVersion string         `json:"InstalledVersion"`
			FixedVersion     string         `json:"FixedVersion"`
			Status           string         `json:"Status"`
			SeveritySource   string         `json:"SeveritySource"`
			PrimaryURL       string         `json:"PrimaryURL"`
			Title            string         `json:"Title"`
			Description      string         `json:"Description"`
			Severity         string         `json:"Severity"`
			CweIDs           []string       `json:"CweIDs"`
			VendorSeverity   map[string]int `json:"VendorSeverity"`
			CVSS             map[string]struct {
				V3Vector string  `json:"V3Vector"`
				V3Score  float64 `json:"V3Score"`
			} `json:"CVSS"`
			References []string `json:"References"`
		} `json:"Vulnerabilities"`
	} `json:"Results"`
}

func RunTrivy(target string) map[string]interface{} {
	log.Println("Running Trivy")
	cmd := exec.Command("trivy", "fs", "--format", "json", target)
	output, _ := cmd.Output()
	ans := categorizeTrivyResults(string(output))
	newcat := utils.ConvertCategorizedResults(ans)
	return newcat
}

func categorizeTrivyResults(results string) map[string][]interface{} {
	categorizedResults := utils.InitAdvancedCategorizedResults()

	var trivyResults TrivyResult
	err := json.Unmarshal([]byte(results), &trivyResults)
	if err != nil {
		log.Printf("Failed to parse Trivy results: %v", err)
		return categorizedResults
	}

	for _, trivyResult := range trivyResults.Results {
		for _, vulnerability := range trivyResult.Vulnerabilities {
			entry := map[string]interface{}{
				"path":              trivyResult.Target,
				"package":           vulnerability.PkgName,
				"check_id":          vulnerability.VulnerabilityID,
				"installed_version": vulnerability.InstalledVersion,
				"fixed_version":     vulnerability.FixedVersion,
				"severity":          vulnerability.Severity,
				"message":           vulnerability.Title,
				"cwe":               vulnerability.CweIDs,
			}

			if isSeverityHigh(vulnerability.Severity) {
				categorizedResults[utils.INFRA_SECURITY] = append(categorizedResults[utils.INFRA_SECURITY], entry)
			}

			if strings.Contains(strings.ToLower(vulnerability.Title), "secret") {
				categorizedResults[utils.SECRET_DETECTION] = append(categorizedResults[utils.SECRET_DETECTION], entry)
			}
		}
	}

	if len(categorizedResults[utils.INFRA_SECURITY]) == 0 {
		categorizedResults[utils.INFRA_SECURITY] = []interface{}{}
	}
	if len(categorizedResults[utils.SCA]) == 0 {
		categorizedResults[utils.SCA] = []interface{}{}
	}

	return categorizedResults
}

func isSeverityHigh(severity string) bool {
	highSeverities := []string{"CRITICAL", "HIGH", "MEDIUM", "LOW"}
	for _, s := range highSeverities {
		if strings.EqualFold(severity, s) {
			return true
		}
	}
	return false
}
