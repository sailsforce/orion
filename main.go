package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	hookChoices        []string
	integrationChoices []string

	hookSelected        map[int]struct{}
	integrationSelected map[int]struct{}

	cursor int
}

func initialModel() model {
	return model{
		hookChoices:        []string{"Defaults", "All", "go mod tidy (default)", "go fmt (default)", "go vet", "go critic", "golangci-lint"},
		integrationChoices: []string{"Defaults", "All", "SonarQube (default)", "Code Factor", "Code Cov"},
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:

		switch msg.String() {

		case "ctrl+c", "q":
			return m, tea.Quit

		case "up", "k":
			if m.cursor > 0 {
				m.cursor--

			}

		case "down", "j":
			if m.cursor < len(m.hookChoices)-1 {
				m.cursor++

			}

		case "enter", " ":
			_, ok := m.hookSelected[m.cursor]
			if ok {
				delete(m.hookSelected, m.cursor)

			} else {
				m.hookSelected[m.cursor] = struct{}{}
			}
		}
	}

	return m, nil
}

func (m model) View() string {
	// The header
	s := "Select which Go Hooks you want to run...\n\n"

	// Iterate over our choices
	for i, hook := range m.hookChoices {

		// Is the cursor pointing at this choice?
		cursor := " " // no cursor
		if m.cursor == i {
			cursor = ">" // cursor!

		}

		// Is this choice selected?
		checked := " " // not selected
		if _, ok := m.hookSelected[i]; ok {
			checked = "x" // selected!

		}

		// Render the row
		s += fmt.Sprintf("%s [%s] %s\n", cursor, checked, hook)
	}

	// The footer
	s += "\nPress q to quit.\n"

	// Send the UI for rendering
	return s

}

func main() {
	p := tea.NewProgram(initialModel())
	if err := p.Start(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
