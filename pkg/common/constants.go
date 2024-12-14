// internal/common/constants.go
package pkg

// LanguageFileExtensions maps programming languages to their file extensions
var LanguageFileExtensions = map[string][]string{
	"py": {
		".py",    // Python files
		".pyc",   // Compiled Python files
		".pyo",   // Optimized Python files
		".pyd",   // Python extension modules
		".ipynb", // Jupyter Notebook files
	},
	"js": {
		".js",   // JavaScript files
		".jsx",  // React JavaScript files
		".mjs",  // ES module JavaScript files
		".json", // JSON files
		".ts",   // TypeScript files
		".tsx",  // React TypeScript files
		".html", // HTML templates
		".css",  // CSS files
		".scss", // SCSS files
		".less", // LESS files
		".sass", // SASS files
	},
	"go": {
		".go",    // Go source files
		".mod",   // Go module files
		".sum",   // Go module sum files
		".cgo",   // Go Cgo files
		".proto", // Protocol Buffer files
	},
}

const (
	// General constants
	Unknown = "UNKNOWN"

	// Issue categories
	DocstringAbsent  = "docstring_absent"
	SecurityIssues   = "security_issues"
	ComplexFunctions = "complex_functions"
	AntipatternsBugs = "antipatterns_bugs"
	SCA              = "sca"

	// Advanced categories
	DeadCode        = "dead_code"
	DuplicateCode   = "duplicate_code"
	SecretDetection = "secret_detection"
	InfraSecurity   = "infra_security"

	// Thresholds
	DuplicateCodeLineThreshold = 25
)
