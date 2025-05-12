package cmd

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/fatih/color"
)

type menuItem struct {
	title       string
	description string
	action      func()
}

type model struct {
	cursor int
	items  []menuItem
}

func initialModel() model {
	return model{
		items: []menuItem{
			{"Scan", "Initiate a new scan", scanAction},
			{"Status", "Check the status of a scan task", statusAction},
			{"Config", "Manage CLI configuration", configAction},
			{"API Audit", "Send an API audit request", apiAction},
			{"Exit", "Exit the application", func() {
				fmt.Println("Goodbye!")
				os.Exit(0)
			}},
		},
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch msg.String() {

		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}

		case "down", "j":
			if m.cursor < len(m.items)-1 {
				m.cursor++
			}

		case "enter":
			// Execute the selected action
			m.items[m.cursor].action()
		}
	}

	return m, nil
}

func (m model) View() string {
	s := "\n"
	s += color.CyanString(`
    █████╗ ██████╗ ███╗   ███╗██╗   ██╗██████╗ 
   ██╔══██╗██╔══██╗████╗ ████║██║   ██║██╔══██╗
   ███████║██████╔╝██╔████╔██║██║   ██║██████╔╝
   ██╔══██║██╔═══╝ ██║╚██╔╝██║██║   ██║██╔═══╝ 
   ██║  ██║██║     ██║ ╚═╝ ██║╚██████╔╝██║     
   ╚═╝  ╚═╝╚═╝     ╚═╝     ╚═╝ ╚═════╝ ╚═╝     
	`)
	s += "\n\nUse ↑/↓ to navigate and Enter to select:\n\n"

	for i, item := range m.items {
		cursor := " " // no cursor
		if m.cursor == i {
			cursor = ">" // cursor
		}
		s += fmt.Sprintf("%s %s - %s\n", cursor, item.title, item.description)
	}

	s += "\nPress Ctrl+C to quit.\n"

	return s
}
