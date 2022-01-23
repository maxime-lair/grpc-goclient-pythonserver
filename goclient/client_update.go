package main

import (
	"fmt"

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
			selected := m.clientChoice.socketChoicesList[m.cursor]
			m.clientEnv.logJournal = append(m.clientEnv.logJournal, fmt.Sprintf("[%s] User tries to select family %s - currently selected %v", m.clientEnv.clientID.Name, selected.Name, m.clientChoice.selectedFamily))
			if m.clientChoice.selectedFamily != nil && *m.clientChoice.selectedFamily == selected {
				m.clientChoice.selectedFamily = nil
			} else {
				if m.clientChoice.selectedFamily == nil {
					m.clientChoice.selectedFamily = &selected
					m.clientEnv.logJournal = append(m.clientEnv.logJournal, fmt.Sprintf("[%s] User selected %s - selectedFamily: %v - selected addr %v",
						m.clientEnv.clientID.Name, selected.Name,
						m.clientChoice.selectedFamily, &selected))
				} else {
					m.clientEnv.logJournal = append(m.clientEnv.logJournal, fmt.Sprintf("[%s] W: User tried to add more than one selection.", m.clientEnv.clientID.Name))
				}
			}
		// The "enter" key to validate
		case "enter":
			m.clientEnv.logJournal = append(m.clientEnv.logJournal, fmt.Sprintf("[%s] selected entry %d - which is choice: %d %s",
				m.clientEnv.clientID.Name, m.cursor,
				m.clientChoice.selectedFamily.Value, m.clientChoice.selectedFamily.Name))
			if m.clientChoice.selectedFamily != nil {
				m.clientEnv.logJournal = append(m.clientEnv.logJournal, fmt.Sprintf("[%s] -------------- GetSocketTypeList --------------", m.clientEnv.clientID.Name))
				m.clientChoice.socketChoicesList = nil // reset choiceList
				m.clientChoice.socketChoicesList, m.clientEnv = socket_get_type_list(m.clientEnv, m.clientChoice)
				if m.clientEnv.err.err != nil {
					return m, nil
				}
				m.state = stateGetType
			}
		}
	}

	// Return the updated model to the Bubble Tea runtime for processing.
	// Note that we're not returning a command.
	return m, nil
}

func (m model) UpdateGetType(msg tea.Msg) (tea.Model, tea.Cmd) {

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
			selected := m.clientChoice.socketChoicesList[m.cursor]
			m.clientEnv.logJournal = append(m.clientEnv.logJournal, fmt.Sprintf("[%s] User tries to select type %s - currently selected %v", m.clientEnv.clientID.Name, selected.Name, m.clientChoice.selectedFamily))
			if m.clientChoice.selectedType != nil && *m.clientChoice.selectedType == selected {
				m.clientChoice.selectedType = nil
			} else {
				if m.clientChoice.selectedType == nil {
					m.clientChoice.selectedType = &selected
					m.clientEnv.logJournal = append(m.clientEnv.logJournal, fmt.Sprintf("[%s] User selected %s - selectedType: %v - selected addr %v",
						m.clientEnv.clientID.Name, selected.Name,
						m.clientChoice.selectedType, &selected))
				} else {
					m.clientEnv.logJournal = append(m.clientEnv.logJournal, fmt.Sprintf("[%s] W: User tried to add more than one selection.", m.clientEnv.clientID.Name))
				}
			}
		// The "enter" key to validate
		case "enter":
			m.clientEnv.logJournal = append(m.clientEnv.logJournal, fmt.Sprintf("[%s] selected entry %d - which is choice: %d %s",
				m.clientEnv.clientID.Name, m.cursor,
				m.clientChoice.selectedType.Value, m.clientChoice.selectedType.Name))
			if m.clientChoice.selectedType != nil {
				m.clientEnv.logJournal = append(m.clientEnv.logJournal, fmt.Sprintf("[%s] -------------- GetSocketProtocolList --------------", m.clientEnv.clientID.Name))
				m.clientChoice.socketChoicesList = nil // reset choiceList
				m.clientChoice.socketChoicesList, m.clientEnv = socket_get_protocol_list(m.clientEnv, m.clientChoice)
				if m.clientEnv.err.err != nil {
					return m, nil
				}
				m.state = stateGetProtocol
			}
		}
	}

	return m, nil
}

func (m model) UpdateGetProtocol(msg tea.Msg) (tea.Model, tea.Cmd) {

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
			selected := m.clientChoice.socketChoicesList[m.cursor]
			m.clientEnv.logJournal = append(m.clientEnv.logJournal, fmt.Sprintf("[%s] User tries to select protocol %s - currently selected %v", m.clientEnv.clientID.Name, selected.Name, m.clientChoice.selectedFamily))
			if m.clientChoice.selectedProtocol != nil && *m.clientChoice.selectedProtocol == selected {
				m.clientChoice.selectedProtocol = nil
			} else {
				if m.clientChoice.selectedProtocol == nil {
					m.clientChoice.selectedProtocol = &selected
					m.clientEnv.logJournal = append(m.clientEnv.logJournal, fmt.Sprintf("[%s] User selected %s - selectedProtocol: %v - selected addr %v",
						m.clientEnv.clientID.Name, selected.Name,
						m.clientChoice.selectedProtocol, &selected))
				} else {
					m.clientEnv.logJournal = append(m.clientEnv.logJournal, fmt.Sprintf("[%s] W: User tried to add more than one selection.", m.clientEnv.clientID.Name))
				}
			}
		// The "enter" key to validate
		case "enter":
			m.clientEnv.logJournal = append(m.clientEnv.logJournal, fmt.Sprintf("[%s] selected entry %d - which is choice: %d %s",
				m.clientEnv.clientID.Name, m.cursor,
				m.clientChoice.selectedProtocol.Value, m.clientChoice.selectedProtocol.Name))
			if m.clientChoice.selectedProtocol != nil {
				m.clientEnv.logJournal = append(m.clientEnv.logJournal, fmt.Sprintf("[%s] -------------- Done --------------", m.clientEnv.clientID.Name))
				m.clientChoice.socketChoicesList = nil // reset choiceList

				m.state = stateDone
			}
		}
	}

	return m, nil
}

func (m model) UpdateDone(msg tea.Msg) (tea.Model, tea.Cmd) {

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
		}
	}

	return m, nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	switch m.state {
	case stateConnect:
		return m.UpdateConnect(msg)
	case stateGetFamily:
		return m.UpdateGetFamily(msg)
	case stateGetType:
		return m.UpdateGetType(msg)
	case stateGetProtocol:
		return m.UpdateGetProtocol(msg)
	case stateDone:
		return m.UpdateDone(msg)
	default:
		return m, nil
	}
}
