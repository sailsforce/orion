package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	goHooks       []string
	integrations  []string
	cursor        int
	selectedHooks map[int]struct{}
}

func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}
	return false
}

func remove(slice []string, s int) []string {
	return append(slice[:s], slice[s+1:]...)
}

func initialModel() model {
	return model{
		goHooks:       []string{"Defaults", "All", "go mod tidy (default)", "go fmt (default)", "go vet", "go critic", "golangci-lint"},
		integrations:  []string{"Defaults", "All", "SonarQube (default)", "CodeFactor", "CodeCov"},
		selectedHooks: make(map[int]struct{}),
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
			if m.cursor < len(m.goHooks)-1 {
				m.cursor++
			}

		case "enter", " ":
			_, ok := m.selectedHooks[m.cursor]
			if ok {
				delete(m.selectedHooks, m.cursor)
			} else {
				m.selectedHooks[m.cursor] = struct{}{}
			}
		}
	}

	//fmt.Printf("selected Hooks: %v", m.selectedHooks)
	return m, nil
}

func (m model) View() string {
	s := "What Go Hooks do you want to run?\n\n"

	// iterate over GoHook choices
	for i, hook := range m.goHooks {
		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}

		checked := " "
		if _, ok := m.selectedHooks[i]; ok {
			checked = "x"
		}

		s += fmt.Sprintf("%s [%s] %s\n", cursor, checked, hook)
	}

	s += "\nPress q to quit.\n"

	return s
}

func main() {
	p := tea.NewProgram(initialModel())
	if err := p.Start(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
