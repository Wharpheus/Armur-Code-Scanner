package cmd

import (
	"armur-cli/internal/config"
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage CLI configuration",
	Long:  `Manage the configuration settings for the Armur CLI.`,
}

var configSetCmd = &cobra.Command{
	Use:   "set [key] [value]",
	Short: "Set a configuration value",
	Long:  `Set a specific configuration key to a new value.`,
	Args:  cobra.ExactArgs(2), // Expect exactly two arguments (key, value)
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.LoadConfig()
		if err != nil {
			color.Red("Error loading configuration: %v", err)
			os.Exit(1)
		}

		key := args[0]
		value := args[1]

		switch key {
		case "api_url":
			cfg.API.URL = value
		case "redis_url":
			cfg.Redis.URL = value
		default:
			color.Red("Invalid configuration key.")
			os.Exit(1)
		}

		if err := config.SaveConfig(cfg); err != nil {
			color.Red("Error saving configuration: %v", err)
			os.Exit(1)
		}

		fmt.Println(color.GreenString("Configuration updated successfully."))
	},
}

var configShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Show current configuration",
	Long:  `Display the current configuration settings.`,
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.LoadConfig()
		if err != nil {
			color.Red("Error loading configuration: %v", err)
			os.Exit(1)
		}

		fmt.Println(color.CyanString("Current Configuration:"))
		fmt.Printf("  API URL: %s\n", cfg.API.URL)
		fmt.Printf("  Redis URL: %s\n", cfg.Redis.URL)
	},
}

func init() {
	rootCmd.AddCommand(configCmd)
	configCmd.AddCommand(configSetCmd)
	configCmd.AddCommand(configShowCmd)
}