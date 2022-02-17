package main

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
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

	return s
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

		case "n":
			m.hooksChosen = true
			return m, nil
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

		case "n":
			m.integrationsChosen = true
			return m, nil
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
	t := "What Go Hooks do you want to run?\n\n"
	t += "%s\n\n"
	t += "j/k, up/down: select | " + "enter, space: choose | " + "n: next" + "q, ctrl+c: quit"

	choices := fmt.Sprintf(
		"%s\n%s\n%s\n%s\n%s\n%s\n%s",
		checkbox("Defaults", m.goHooks),
		checkbox("All", m.goHooks),
		checkbox("go mod tidy (default)", m.goHooks),
		checkbox("go fmt (default)", m.goHooks),
		checkbox("go vet", m.goHooks),
		checkbox("go critic", m.goHooks),
		checkbox("golangci-lint", m.goHooks),
	)

	return fmt.Sprintf(t, choices)
}

func integrationChoicesView(m model) string {
	t := "What Go Integrations do you want to run?\n\n"
	t += "%s\n\n"
	t += "j/k, up/down: select | " + "enter, space: choose | " + "n: next" + "q, ctrl+c: quit"

	choices := fmt.Sprintf(
		"%s\n%s\n%s\n%s\n%s",
		checkbox("Defaults", m.integrations),
		checkbox("All", m.integrations),
		checkbox("SonarQube (default)", m.integrations),
		checkbox("CodeFactor", m.integrations),
		checkbox("CodeCov", m.integrations),
	)

	return fmt.Sprintf(t, choices)
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

func checkbox(label string, list []string) string {
	checked := false
	for _, v := range list {
		if v == label {
			checked = true
		}
	}

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
