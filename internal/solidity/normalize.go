package solidity

import "strings"

// NormalizeSeverity maps tool-specific severities to a common scale.
func NormalizeSeverity(tool, s string) string {
	level := strings.ToUpper(strings.TrimSpace(s))
	switch tool {
	case "slither":
		// Slither uses: CRITICAL/HIGH/MEDIUM/LOW/INFO (impact)
		return mapToCommon(level)
	case "mythril":
		// Mythril often uses High/Medium/Low
		return mapToCommon(level)
	case "oyente", "securify", "smartcheck", "semgrep":
		return mapToCommon(level)
	default:
		return mapToCommon(level)
	}
}

func mapToCommon(level string) string {
	switch level {
	case "CRITICAL":
		return "CRITICAL"
	case "HIGH":
		return "HIGH"
	case "MEDIUM":
		return "MEDIUM"
	case "LOW":
		return "LOW"
	default:
		return "INFO"
	}
}

// MapRuleToSWC returns an SWC code for a given tool+rule identifier when known.
func MapRuleToSWC(tool, rule string) string {
	key := strings.ToLower(strings.TrimSpace(tool)) + ":" + strings.ToLower(strings.TrimSpace(rule))
	switch key {
	case "slither:reentrancy-eth", "slither:reentrancy-no-eth", "slither:reentrancy-unlimited-gas":
		return "SWC-107"
	case "slither:tx-origin":
		return "SWC-115"
	case "slither:delegatecall":
		return "SWC-112"
	case "slither:selfdestruct", "slither:suicidal":
		return "SWC-106"
	case "slither:uninitialized-state", "slither:uninitialized-storage":
		return "SWC-109"
	case "slither:controlled-delegatecall":
		return "SWC-112"
	case "slither:unchecked-transfer":
		return "SWC-104"
	case "slither:arbitrary-send-eth":
		return "SWC-105"
	default:
		return ""
	}
}

// MapRuleToCWE returns a CWE ID for a given tool+rule identifier when known.
func MapRuleToCWE(tool, rule string) string {
	key := strings.ToLower(strings.TrimSpace(tool)) + ":" + strings.ToLower(strings.TrimSpace(rule))
	switch key {
	case "slither:reentrancy-eth", "slither:reentrancy-no-eth", "slither:reentrancy-unlimited-gas":
		return "CWE-841"
	case "slither:tx-origin":
		return "CWE-285"
	case "slither:delegatecall", "slither:controlled-delegatecall":
		return "CWE-829"
	case "slither:selfdestruct", "slither:suicidal":
		return "CWE-284"
	case "slither:uninitialized-state", "slither:uninitialized-storage":
		return "CWE-665"
	default:
		return ""
	}
}
