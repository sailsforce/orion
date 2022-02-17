package main

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/muesli/reflow/indent"
)

type model struct {
	goHooks              []string
	integrations         []string
	cursor               int
	selectedHooks        map[int]struct{}
	selectedIntegrations map[int]struct{}
	hooksChosen          bool
	integrationsChosen   bool
	quitting             bool
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// Make sure these keys always quit
	if msg, ok := msg.(tea.KeyMsg); ok {
		k := msg.String()
		if k == "q" || k == "ctrl+c" {
			m.quitting = true
			return m, tea.Quit
		}
	}

	// Hand off the message and model to the appropriate update func
	if !m.hooksChosen {
		return updateHookChoices(msg, m)
	} else if !m.integrationsChosen {
		return updateIntegrationChoices(msg, m)
	}
	return updateChosen(msg, m)
}

func (m model) View() string {
	var s string
	if m.quitting {
		return "\n See you later!\n\n"
	}

	if !m.hooksChosen {
		s = hookChoicesView(m)
	} else if !m.integrationsChosen {
		s = integrationChoicesView(m)
	} else {
		s = chosenView(m)
	}

	return indent.String("\n"+s+"\n\n", 2)
}

func updateHookChoices(msg tea.Msg, m model) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {

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

	return m, nil
}

func updateIntegrationChoices(msg tea.Msg, m model) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {

		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}

		case "down", "j":
			if m.cursor < len(m.goHooks)-1 {
				m.cursor++
			}

		case "enter", " ":
			_, ok := m.selectedIntegrations[m.cursor]
			if ok {
				delete(m.selectedIntegrations, m.cursor)
			} else {
				m.selectedIntegrations[m.cursor] = struct{}{}
			}
		}
	}

	return m, nil
}

func updateChosen(msg tea.Msg, m model) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {

		case "enter":
			m.quitting = true

		}
	}

	return m, nil
}

func hookChoicesView(m model) string {
	c := m.cursor

	t := "What Go Hooks do you want to run?\n\n"
	t += "%s\n\n"
	t += "j/k, up/down: select" + "enter, space: choose" + "q, ctrl+c: quit"

	choices := fmt.Sprintf(
		"%s\n%s\n%s\n%s\n%s\n%s\n%s",
		checkbox("Defaults", c == 0),
		checkbox("All", c == 1),
		checkbox("go mod tidy (default)", c == 2),
		checkbox("go fmt (default)", c == 3),
		checkbox("go vet", c == 4),
		checkbox("go critic", c == 5),
		checkbox("golangci-lint", c == 6),
	)

	return choices
}

func integrationChoicesView(m model) string {
	c := m.cursor

	t := "What Go Integrations do you want to run?\n\n"
	t += "%s\n\n"
	t += "j/k, up/down: select" + "enter, space: choose" + "q, ctrl+c: quit"

	choices := fmt.Sprintf(
		"%s\n%s\n%s\n%s\n%s",
		checkbox("Defaults", c == 0),
		checkbox("All", c == 1),
		checkbox("SonarQube (default)", c == 2),
		checkbox("CodeFactor", c == 3),
		checkbox("CodeCov", c == 4),
	)

	return choices
}

func chosenView(m model) string {
	msg := "You chose these Go Hooks:\n"
	for _, v := range m.selectedHooks {
		msg += fmt.Sprintf("%s\n", v)
	}

	msg += "\n\nYou chose these Integrations:\n"
	for _, v := range m.selectedIntegrations {
		msg += fmt.Sprintf("%s\n", v)
	}

	return msg
}

func checkbox(label string, checked bool) string {
	if checked {
		return "[x] " + label
	}
	return fmt.Sprintf("[ ] %s", label)
}

func main() {
	initialModel := model{[]string{}, []string{}, 0, make(map[int]struct{}), make(map[int]struct{}), false, false, false}
	p := tea.NewProgram(initialModel)
	if err := p.Start(); err != nil {
		fmt.Printf("error running the program:\n %v", err)
	}
}
