package main

import (
	"fmt"
	"log"

	tea "github.com/charmbracelet/bubbletea"
)

// state connect, we request "enter" to send the request to the server
func (m model) UpdateConnect(msg tea.Msg) (tea.Model, tea.Cmd) {

	switch msg := msg.(type) {

	// Is it a key press?
	case tea.KeyMsg:

		// Cool, what was the actual key pressed?
		switch msg.String() {

		// These keys should exit the program.
		case "ctrl+c", "q":
			return m, tea.Quit

		// The "enter" key and the spacebar (a literal space) toggle
		// the selected state for the item that the cursor is pointing at.
		case "enter", " ":
			m.clientEnv.logJournal = append(m.clientEnv.logJournal, fmt.Sprintf("[%s] -------------- SendSocketTree --------------", m.clientEnv.clientID.Name))
			m.clientChoice.socketChoicesList, m.clientEnv = socket_get_family_list(m.clientEnv)
			if m.clientEnv.err.err != nil {
				return m, nil
			}
			m.state = stateGetFamily
		}
	}

	// Return the updated model to the Bubble Tea runtime for processing.
	// Note that we're not returning a command.
	return m, nil
}

func (m model) UpdateGetFamily(msg tea.Msg) (tea.Model, tea.Cmd) {

	switch msg := msg.(type) {

	// Is it a key press?
	case tea.KeyMsg:

		// Cool, what was the actual key pressed?
		switch msg.String() {

		// These keys should exit the program.
		case "ctrl+c", "q":
			return m, tea.Quit

		// The "up" and "k" keys move the cursor up
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}

		// The "down" and "j" keys move the cursor down
		case "down", "j":
			if m.cursor < len(m.clientChoice.socketChoicesList)-1 {
				m.cursor++
			}

		// the spacebar (a literal space) toggle
		// the selected state for the item that the cursor is pointing at.
		case " ":
			_, ok := m.selected[m.cursor]
			if ok {
				delete(m.selected, m.cursor)
			} else {
				if len(m.selected) == 0 {
					m.selected[m.cursor] = struct{}{}

				} else {
					m.clientEnv.logJournal = append(m.clientEnv.logJournal, fmt.Sprintf("[%s] W: User tried to add more than one selection.", m.clientEnv.clientID))
				}
			}
		// The "enter" key to validate
		case "enter":
			log.Printf("selected: %d", m.selected[m.cursor])

		}
	}

	// Return the updated model to the Bubble Tea runtime for processing.
	// Note that we're not returning a command.
	return m, nil
}

func (m model) UpdateGetType(msg tea.Msg) (tea.Model, tea.Cmd) {

	return m, nil
}

func (m model) UpdateGetProtocol(msg tea.Msg) (tea.Model, tea.Cmd) {

	return m, nil
}

func (m model) UpdateDone(msg tea.Msg) (tea.Model, tea.Cmd) {

	return m, nil
}
