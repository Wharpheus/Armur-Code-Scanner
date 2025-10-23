package solidity

import (
	"bufio"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// Config holds detected Solidity project settings
type Config struct {
	Version     string
	Remappings  []string
	Framework   string // foundry|hardhat|truffle|unknown
	SourcesDir  string // contracts/ or src/
}

// DetectSolidityConfig inspects common framework config files to extract solc version and remappings.
func DetectSolidityConfig(root string) Config {
	cfg := Config{Framework: "unknown", SourcesDir: ""}

	// Foundry (foundry.toml)
	foundry := filepath.Join(root, "foundry.toml")
	if fileExists(foundry) {
		cfg.Framework = "foundry"
		cfg.SourcesDir = firstExisting(root, []string{"src", "contracts"})
		cfg.Version = parseFoundryVersion(foundry)
		cfg.Remappings = parseFoundryRemappings(root)
		return cfg
	}

	// Hardhat (hardhat.config.js/ts)
	for _, name := range []string{"hardhat.config.ts", "hardhat.config.js"} {
		p := filepath.Join(root, name)
		if fileExists(p) {
			cfg.Framework = "hardhat"
			cfg.SourcesDir = firstExisting(root, []string{"contracts", "src"})
			cfg.Version = parseHardhatVersion(p)
			return cfg
		}
	}

	// Truffle (truffle-config.js)
	truffle := filepath.Join(root, "truffle-config.js")
	if fileExists(truffle) {
		cfg.Framework = "truffle"
		cfg.SourcesDir = firstExisting(root, []string{"contracts", "src"})
		cfg.Version = parseTruffleVersion(truffle)
		return cfg
	}

	// Fallbacks
	cfg.SourcesDir = firstExisting(root, []string{"contracts", "src"})
	return cfg
}

func firstExisting(root string, dirs []string) string {
	for _, d := range dirs {
		if dirExists(filepath.Join(root, d)) {
			return d
		}
	}
	return ""
}

func fileExists(p string) bool {
	st, err := os.Stat(p)
	return err == nil && !st.IsDir()
}

func dirExists(p string) bool {
	st, err := os.Stat(p)
	return err == nil && st.IsDir()
}

func parseFoundryVersion(path string) string {
	// solc_version = "0.8.23" or via profile
	b, _ := os.ReadFile(path)
	re := regexp.MustCompile(`(?mi)^\s*solc_version\s*=\s*"([^"]+)"`)
	m := re.FindStringSubmatch(string(b))
	if len(m) > 1 {
		return m[1]
	}
	return ""
}

func parseFoundryRemappings(root string) []string {
	remaps := []string{}
	for _, file := range []string{"remappings.txt", filepath.Join(".gitmodules") } {
		p := filepath.Join(root, file)
		if !fileExists(p) { continue }
		f, err := os.Open(p)
		if err != nil { continue }
		s := bufio.NewScanner(f)
		for s.Scan() {
			line := strings.TrimSpace(s.Text())
			if line == "" || strings.HasPrefix(line, "#") { continue }
			remaps = append(remaps, line)
		}
		_ = f.Close()
	}
	return remaps
}

func parseHardhatVersion(path string) string {
	b, _ := os.ReadFile(path)
	// solidity: "0.8.23" or solidity: { version: "0.8.23" }
	re := regexp.MustCompile(`solidity\s*:\s*(?:\{[^}]*version\s*:\s*"([^"]+)"|"([^"]+)")`)
	m := re.FindStringSubmatch(string(b))
	if len(m) > 2 {
		if m[1] != "" { return m[1] }
		return m[2]
	}
	return ""
}

func parseTruffleVersion(path string) string {
	b, _ := os.ReadFile(path)
	re := regexp.MustCompile(`compilers\s*:\s*\{\s*solc\s*:\s*\{[^}]*version\s*:\s*"([^"]+)"`)
	m := re.FindStringSubmatch(string(b))
	if len(m) > 1 { return m[1] }
	return ""
}
