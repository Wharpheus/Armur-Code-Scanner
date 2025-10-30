package tasks

import (
	"context"
	"crypto/sha1"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
)

var (
	resultStoreRedisClient     *redis.Client
	resultStoreRedisClientOnce sync.Once
)

// Finding represents a normalized finding across tools.
// This allows deduplication and consistent reporting.
// Kept minimal to avoid breaking existing API contracts.
 type Finding struct {
	Path       string   `json:"path"`
	StartLine  int      `json:"start_line"`
	EndLine    int      `json:"end_line"`
	Message    string   `json:"message"`
	Severity   string   `json:"severity"`
	RuleID     string   `json:"rule_id"`
	SWC        string   `json:"swc,omitempty"`
	CWE        string   `json:"cwe,omitempty"`
	Category   string   `json:"category,omitempty"`
	Tool       string   `json:"tool"`
	Sources    []string `json:"sources,omitempty"` // tools that reported it
}

func initResultStoreRedisClient() *redis.Client {
	resultStoreRedisClientOnce.Do(func() {
		addr := os.Getenv("REDIS_ADDR")
		if addr == "" {
			addr = "localhost:6379"
		}

		password := os.Getenv("REDIS_PASSWORD")
		db := getEnvAsIntResultStore("REDIS_DB", 0)

		resultStoreRedisClient = redis.NewClient(&redis.Options{
			Addr:     addr,
			Password: password,
			DB:       db,
		})
	})
	return resultStoreRedisClient
}

func getEnvAsIntResultStore(key string, defaultValue int) int {
	value, exists := os.LookupEnv(key)
	if !exists {
		return defaultValue
	}
	intValue, err := strconv.Atoi(value)
	if err != nil {
		return defaultValue
	}
	return intValue
}

// DeduplicateFindings merges similar findings keeping highest severity and aggregated sources.
func DeduplicateFindings(findings []Finding) []Finding {
	if len(findings) == 0 {
		return findings
	}
	byKey := make(map[string]Finding)
	sevRank := map[string]int{"CRITICAL":4, "HIGH":3, "MEDIUM":2, "LOW":1, "INFO":0}
	for _, f := range findings {
		key := findingKey(f)
		if existing, ok := byKey[key]; ok {
			// keep highest severity
			if sevRank[strings.ToUpper(f.Severity)] > sevRank[strings.ToUpper(existing.Severity)] {
				existing.Severity = f.Severity
				existing.Message = f.Message
				existing.RuleID = f.RuleID
			}
			// merge sources
			existing.Sources = mergeStrings(existing.Sources, append([]string{existing.Tool}, f.Tool))
			byKey[key] = existing
		} else {
			f.Sources = []string{f.Tool}
			byKey[key] = f
		}
	}
	// stable order
	keys := make([]string, 0, len(byKey))
	for k := range byKey { keys = append(keys, k) }
	sort.Strings(keys)
	out := make([]Finding, 0, len(byKey))
	for _, k := range keys { out = append(out, byKey[k]) }
	return out
}

func findingKey(f Finding) string {
	base := fmt.Sprintf("%s|%d|%d|%s|%s|%s", f.Path, f.StartLine, f.EndLine, strings.ToLower(f.RuleID), strings.ToLower(f.SWC), strings.ToLower(f.CWE))
	h := sha1.Sum([]byte(base))
	return fmt.Sprintf("%x", h)
}

func mergeStrings(a []string, b []string) []string {
	m := make(map[string]struct{}, len(a)+len(b))
	for _, s := range a { m[s] = struct{}{} }
	for _, s := range b { m[s] = struct{}{} }
	res := make([]string, 0, len(m))
	for s := range m { res = append(res, s) }
	sort.Strings(res)
	return res
}

func SaveTaskResult(taskID string, result map[string]any) error {
	ctx := context.Background()

	resultData, err := json.Marshal(result)
	if err != nil {
		return err
	}

	client := initResultStoreRedisClient()
	return client.Set(ctx, taskID, resultData, 24*time.Hour).Err()
}

func GetTaskResult(taskID string) (any, error) {
	ctx := context.Background()

	client := initResultStoreRedisClient()
	resultData, err := client.Get(ctx, taskID).Result()
	if err != nil {
		if err.Error() == "redis: nil" {
			return nil, errors.New("task result not found")
		}
		return nil, err
	}

	var result interface{}
	if err := json.Unmarshal([]byte(resultData), &result); err != nil {
		return nil, err
	}
	return result, nil
}
