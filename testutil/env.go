package testutil

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// LoadDotenv loads environment variables from a .env file if present.
// It is safe to call multiple times; subsequent calls are no-ops.
func LoadDotenv() {
	// Only attempt to load if a .env file exists in repository root
	if _, err := os.Stat(".env"); err == nil {
		if err := godotenv.Load(".env"); err != nil {
			log.Printf("testutil: failed to load .env: %v", err)
		}
	}
}
