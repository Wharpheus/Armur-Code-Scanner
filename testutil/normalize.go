package testutil

import (
	"regexp"
)

var (
	absPathRe  = regexp.MustCompile(`(?m)/[A-Za-z0-9_./-]+`)
	timestampRe = regexp.MustCompile(`\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}Z`)
)

// NormalizeOutput strips unstable elements like absolute paths and timestamps to make output comparable.
func NormalizeOutput(s string) string {
	s = absPathRe.ReplaceAllString(s, "/path")
	s = timestampRe.ReplaceAllString(s, "0000-00-00T00:00:00Z")
	return s
}
