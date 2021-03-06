package main

import (
	"fmt"

	key "github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
)

var DefaultKeyMap = keyMap{
	Up: key.NewBinding(
		key.WithKeys("k", "up"),        // actual keybindings
		key.WithHelp("↑/k", "move up"), // corresponding help text
	),
	Down: key.NewBinding(
		key.WithKeys("j", "down"),
		key.WithHelp("↓/j", "move down"),
	),
	Enter: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("<enter>", "Validate entry"),
	),
	Space: key.NewBinding(
		key.WithKeys(" "),
		key.WithHelp("<spacebar>", "Select entry"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q", "ctrl+c"),
		key.WithHelp("⌘+c/q", "Quit client"),
	),
}

// For messages that contain errors it's often handy to also implement the
// error interface on the message.
func (e errMsg) Error() string { return e.err.Error() }

// state connect, we request "enter" to send the request to the server
func (m model) UpdateConnect(msg tea.Msg) (tea.Model, tea.Cmd) {

	switch msg := msg.(type) {

	case tickMsg:
		if m.progress.Percent() >= 0.2 {
			return m, nil
		}

		// Note that you can also use progress.Model.SetPercent to set the
		// percentage value explicitly, too.
		cmd := m.progress.IncrPercent(0.05)
		return m, tea.Batch(tickCmd(), cmd)

	// FrameMsg is sent when the progress bar wants to animate itself
	case progress.FrameMsg:
		progressModel, cmd := m.progress.Update(msg)
		m.progress = progressModel.(progress.Model)
		return m, cmd

	// Is it a key press?
	case tea.KeyMsg:

		// Cool, what was the actual key pressed?
		switch {

		// These keys should exit the program.
		case key.Matches(msg, DefaultKeyMap.Quit):
			return m, tea.Quit

		case key.Matches(msg, DefaultKeyMap.Space):
			m.clientEnv.logJournal = append(m.clientEnv.logJournal, fmt.Sprintf("Add percent %v", m.progress))
			return m, m.progress.IncrPercent(0.25)
		// The "enter" key and the spacebar (a literal space) toggle
		// the selected state for the item that the cursor is pointing at.
		case key.Matches(msg, DefaultKeyMap.Enter):
			m.clientEnv.logJournal = append(m.clientEnv.logJournal, fmt.Sprintf("Created client %p with id %s", &m.clientEnv.client, m.clientEnv.clientID.Name))
			m.clientEnv.logJournal = append(m.clientEnv.logJournal, fmt.Sprintf("[%s] -------------- SendSocketTree --------------", m.clientEnv.clientID.Name))
			m.clientChoice.socketChoicesList, m.clientEnv = socket_get_family_list(m.clientEnv)
			if m.clientEnv.err.err != nil {
				return m, nil
			}
			m.state = stateGetFamily
			return m, tickCmd()
		}
	default:
		return m, m.spinner.Tick
	}

	// Return the updated model to the Bubble Tea runtime for processing.
	return m, nil
}

func (m model) UpdateGetFamily(msg tea.Msg) (tea.Model, tea.Cmd) {

	switch msg := msg.(type) {

	case tickMsg:
		if m.progress.Percent() >= 0.4 {
			return m, nil
		}

		// Note that you can also use progress.Model.SetPercent to set the
		// percentage value explicitly, too.
		cmd := m.progress.IncrPercent(0.05)
		return m, tea.Batch(tickCmd(), cmd)

	// FrameMsg is sent when the progress bar wants to animate itself
	case progress.FrameMsg:
		progressModel, cmd := m.progress.Update(msg)
		m.progress = progressModel.(progress.Model)
		return m, cmd

	// Is it a key press?
	case tea.KeyMsg:

		// Cool, what was the actual key pressed?
		switch {

		// These keys should exit the program.
		case key.Matches(msg, DefaultKeyMap.Quit):
			return m, tea.Quit

		// The "up" and "k" keys move the cursor up
		case key.Matches(msg, DefaultKeyMap.Up):
			if m.cursor > 0 {
				m.cursor--
			}

		// The "down" and "j" keys move the cursor down
		case key.Matches(msg, DefaultKeyMap.Down):
			if m.cursor < len(m.clientChoice.socketChoicesList)-1 {
				m.cursor++
			}

		// the spacebar (a literal space) toggle
		// the selected state for the item that the cursor is pointing at.
		case key.Matches(msg, DefaultKeyMap.Space):
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
		case key.Matches(msg, DefaultKeyMap.Enter):
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
				m.cursor = 0 // reset cursor
				return m, tickCmd()
			}
		}
	default:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	}

	// Return the updated model to the Bubble Tea runtime for processing.
	return m, nil
}

func (m model) UpdateGetType(msg tea.Msg) (tea.Model, tea.Cmd) {

	switch msg := msg.(type) {

	case tickMsg:
		if m.progress.Percent() >= 0.6 {
			return m, nil
		}

		// Note that you can also use progress.Model.SetPercent to set the
		// percentage value explicitly, too.
		cmd := m.progress.IncrPercent(0.05)
		return m, tea.Batch(tickCmd(), cmd)

	// FrameMsg is sent when the progress bar wants to animate itself
	case progress.FrameMsg:
		progressModel, cmd := m.progress.Update(msg)
		m.progress = progressModel.(progress.Model)
		return m, cmd

	// Is it a key press?
	case tea.KeyMsg:

		// Cool, what was the actual key pressed?
		switch {

		// These keys should exit the program.
		case key.Matches(msg, DefaultKeyMap.Quit):
			return m, tea.Quit

		// The "up" and "k" keys move the cursor up
		case key.Matches(msg, DefaultKeyMap.Up):
			if m.cursor > 0 {
				m.cursor--
			}

		// The "down" and "j" keys move the cursor down
		case key.Matches(msg, DefaultKeyMap.Down):
			if m.cursor < len(m.clientChoice.socketChoicesList)-1 {
				m.cursor++
			}

		// the spacebar (a literal space) toggle
		// the selected state for the item that the cursor is pointing at.
		case key.Matches(msg, DefaultKeyMap.Space):
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
		case key.Matches(msg, DefaultKeyMap.Enter):
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
				m.cursor = 0 // reset cursor
				return m, tickCmd()
			}
		}
	default:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	}

	return m, nil
}

func (m model) UpdateGetProtocol(msg tea.Msg) (tea.Model, tea.Cmd) {

	switch msg := msg.(type) {

	case tickMsg:
		if m.progress.Percent() >= 0.8 {
			return m, nil
		}

		// Note that you can also use progress.Model.SetPercent to set the
		// percentage value explicitly, too.
		cmd := m.progress.IncrPercent(0.05)
		return m, tea.Batch(tickCmd(), cmd)

	// FrameMsg is sent when the progress bar wants to animate itself
	case progress.FrameMsg:
		progressModel, cmd := m.progress.Update(msg)
		m.progress = progressModel.(progress.Model)
		return m, cmd

	// Is it a key press?
	case tea.KeyMsg:

		// Cool, what was the actual key pressed?
		switch {

		// These keys should exit the program.
		case key.Matches(msg, DefaultKeyMap.Quit):
			return m, tea.Quit

		// The "up" and "k" keys move the cursor up
		case key.Matches(msg, DefaultKeyMap.Up):
			if m.cursor > 0 {
				m.cursor--
			}

		// The "down" and "j" keys move the cursor down
		case key.Matches(msg, DefaultKeyMap.Down):
			if m.cursor < len(m.clientChoice.socketChoicesList)-1 {
				m.cursor++
			}

		// the spacebar (a literal space) toggle
		// the selected state for the item that the cursor is pointing at.
		case key.Matches(msg, DefaultKeyMap.Space):
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
		case key.Matches(msg, DefaultKeyMap.Enter):
			m.clientEnv.logJournal = append(m.clientEnv.logJournal, fmt.Sprintf("[%s] selected entry %d - which is choice: %d %s",
				m.clientEnv.clientID.Name, m.cursor,
				m.clientChoice.selectedProtocol.Value, m.clientChoice.selectedProtocol.Name))
			if m.clientChoice.selectedProtocol != nil {
				m.clientEnv.logJournal = append(m.clientEnv.logJournal, fmt.Sprintf("[%s] -------------- Done --------------", m.clientEnv.clientID.Name))
				m.clientChoice.socketChoicesList = nil // reset choiceList

				m.state = stateDone
				m.cursor = 0 // reset cursor
				return m, tickCmd()
			}
		}
	default:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	}

	return m, nil
}

func (m model) UpdateDone(msg tea.Msg) (tea.Model, tea.Cmd) {

	switch msg := msg.(type) {

	case tickMsg:
		if m.progress.Percent() >= 1.0 {
			return m, nil
		}

		// Note that you can also use progress.Model.SetPercent to set the
		// percentage value explicitly, too.
		cmd := m.progress.IncrPercent(0.05)
		return m, tea.Batch(tickCmd(), cmd)

	// FrameMsg is sent when the progress bar wants to animate itself
	case progress.FrameMsg:
		progressModel, cmd := m.progress.Update(msg)
		m.progress = progressModel.(progress.Model)
		return m, cmd

	// Is it a key press?
	case tea.KeyMsg:

		// Cool, what was the actual key pressed?
		switch {

		// These keys should exit the program.
		case key.Matches(msg, DefaultKeyMap.Quit):
			return m, tea.Quit
		}
	default:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
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
