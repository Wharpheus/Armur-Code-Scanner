package cmd

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "armur",
	Short: "Armur Code Scanner CLI",
	Long:  `A fast and powerful CLI for interacting with the Armur Code Scanner service.`,
}

func Execute() {
	if len(os.Args) == 1 {
		// No arguments provided; launch interactive menu
		p := tea.NewProgram(initialModel())
		if _, err := p.Run(); err != nil {
			fmt.Println("Error running program:", err)
			os.Exit(1)
		}
		return
	}

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
