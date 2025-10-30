package tasks

import (
	"armur-codescanner/internal/solidity"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"sort"
	"strings"
)

// GenericIssue represents a minimal common shape many tool adapters already emit.
// This adapter helps convert existing per-tool issue maps into Finding for deduplication.
type GenericIssue struct {
	Path     string
	Line     int
	Message  string
	Severity string
	RuleID   string
	SWC      string
	CWE      string
	Tool     string
	Category string
}

// ToFinding converts a GenericIssue to a Finding, filling StartLine/EndLine and normalization.
func (gi GenericIssue) ToFinding() Finding {
	sev := solidity.NormalizeSeverity(gi.Tool, gi.Severity)
	return Finding{
		Path:      gi.Path,
		StartLine: gi.Line,
		EndLine:   gi.Line,
		Message:   gi.Message,
		Severity:  sev,
		RuleID:    gi.RuleID,
		SWC:       gi.SWC,
		CWE:       gi.CWE,
		Category:  gi.Category,
		Tool:      gi.Tool,
	}
}

// NormalizeFromMap attempts to read common fields from a raw map[string]any emitted by tools.
// Unknown/missing fields are set to defaults that won't break dedup.
func NormalizeFromMap(m map[string]any) Finding {
	getS := func(k string) string {
		if v, ok := m[k]; ok && v != nil { return fmt.Sprintf("%v", v) }
		return ""
	}
	getI := func(k string) int {
		if v, ok := m[k]; ok && v != nil {
			s := fmt.Sprintf("%v", v)
			if s == "-" || s == "" { return 0 }
			var n int
			fmt.Sscanf(s, "%d", &n)
			return n
		}
		return 0
	}

	tool := strings.ToLower(getS("tool"))
	sev := solidity.NormalizeSeverity(tool, getS("severity"))
	rule := getS("rule")
	if rule == "" {
		// synthesize a stable rule id from message when missing
		h := sha1.Sum([]byte(getS("message")))
		rule = tool + ":" + hex.EncodeToString(h[:4])
	}
	return Finding{
		Path:      getS("path"),
		StartLine: getI("line"),
		EndLine:   getI("line"),
		Message:   getS("message"),
		Severity:  sev,
		RuleID:    rule,
		SWC:       getS("swc"),
		CWE:       getS("cwe"),
		Category:  getS("category"),
		Tool:      tool,
	}
}

// SortFindings sorts by severity desc, then path, then start line.
func SortFindings(findings []Finding) {
	rank := map[string]int{"CRITICAL":4, "HIGH":3, "MEDIUM":2, "LOW":1, "INFO":0}
	sort.SliceStable(findings, func(i, j int) bool {
		li, lj := strings.ToUpper(findings[i].Severity), strings.ToUpper(findings[j].Severity)
		if rank[li] != rank[lj] { return rank[li] > rank[lj] }
		if findings[i].Path != findings[j].Path { return findings[i].Path < findings[j].Path }
		return findings[i].StartLine < findings[j].StartLine
	})
}
