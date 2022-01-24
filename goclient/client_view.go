package main

import (
	"fmt"

	"github.com/charmbracelet/bubbles/key"
)

// ShortHelp returns keybindings to be shown in the mini help view. It's part
// of the key.Map interface.
func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Quit}
}

// FullHelp returns keybindings for the expanded help view. It's part of the
// key.Map interface.
func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down},             // first column
		{k.Space, k.Enter, k.Quit}, // second column
	}
}

// The header
func (m model) printHeader() string {

	var s string
	if m.clientEnv.clientID == nil {
		s += fmt.Sprintf("------- Client ID undefined %p -------\n", m.clientEnv.client)
	} else {
		s += fmt.Sprintf("------- %s -------\n", m.clientEnv.clientID.Name)
	}
	switch m.state {
	case stateConnect:
		s += "Request process starting..\n\n"
	case stateGetFamily:
		if m.clientChoice.selectedFamily != nil {
			s += fmt.Sprintf("\nCurrently selected value : [%d] %s\n", m.clientChoice.selectedFamily.Value, m.clientChoice.selectedFamily.Name)
		}
		s += "Please select your socket family..\n"
	case stateGetType:
		s += fmt.Sprintf("Requesting socket type list for family %s \n\n", m.clientChoice.selectedFamily.Name)
	case stateGetProtocol:
		s += fmt.Sprintf("Requesting socket protocol list for family %s and type %s:\n\n", m.clientChoice.selectedFamily.Name, m.clientChoice.selectedType.Name)
	case stateDone:
		s += "Request process done, your final choice will appear below:\n\n"
	default:
		s += "Unknown state, exiting..\n\n"
	}

	return s
}

func (m model) printChoices(i int, selectedValue *socketChoice, possibleChoice socketChoice) string {
	// Is the cursor pointing at this choice?
	cursor := " " // no cursor
	if m.cursor == i {
		cursor = ">" // cursor!
	}

	// Is this choice selected?
	checked := " " // not selected
	if selectedValue != nil && *selectedValue == possibleChoice {
		checked = "x" // selected!
	}
	return fmt.Sprintf("%s [%s] %d - %s \n", cursor, checked, possibleChoice.Value, possibleChoice.Name)
}

func (m model) printFooter() string {

	// The footer
	m.help.ShowAll = true // show all help
	var s string
	s += "\n"
	s += m.help.View(m.keys)
	s += "Only one selection at a time possible.\n"

	return s
}

func (m model) ViewConnect() string {
	var s string
	s += m.printHeader()
	// TODO add loading bar
	s += "Press enter to start ..\n"
	s += m.printFooter()
	s += m.ViewLogs()
	return s
}

func (m model) ViewGetFamily() string {
	var s string
	s += m.printHeader()

	// Iterate over our choices
	for i, choice := range m.clientChoice.socketChoicesList {
		s += m.printChoices(i, m.clientChoice.selectedFamily, choice)
	}

	s += m.printFooter()
	s += m.ViewLogs()
	return s
}

func (m model) ViewGetType() string {
	var s string
	s += m.printHeader()
	// Iterate over our choices
	for i, choice := range m.clientChoice.socketChoicesList {
		s += m.printChoices(i, m.clientChoice.selectedType, choice)
	}

	s += m.printFooter()
	s += m.ViewLogs()
	return s
}

func (m model) ViewGetProtocol() string {
	var s string
	s += m.printHeader()
	// Iterate over our choices
	for i, choice := range m.clientChoice.socketChoicesList {
		s += m.printChoices(i, m.clientChoice.selectedProtocol, choice)
	}
	s += m.printFooter()
	s += m.ViewLogs()
	return s
}

func (m model) ViewDone() string {
	var s string
	s += m.printHeader()

	s += "You choose the following parameters for your socket:\n"
	s += fmt.Sprintf("\t-> Family: %d - %s\n", m.clientChoice.selectedFamily.Value, m.clientChoice.selectedFamily.Name)
	s += fmt.Sprintf("\t--> Type: %d - %s\n", m.clientChoice.selectedType.Value, m.clientChoice.selectedType.Name)
	s += fmt.Sprintf("\t---> Protocol: %d - %s\n", m.clientChoice.selectedProtocol.Value, m.clientChoice.selectedProtocol.Name)

	s += m.printFooter()
	s += m.ViewLogs()
	return s
}

func (m model) ViewLogs() string {
	// print log journal
	var s string
	var recentLogs []string
	if len(m.clientEnv.logJournal) > 5 {
		recentLogs = m.clientEnv.logJournal[len(m.clientEnv.logJournal)-5:]
	} else {
		recentLogs = m.clientEnv.logJournal
	}
	s += fmt.Sprintf("\n------ last %d logs (total %d)------\n", len(m.clientEnv.logJournal), len(recentLogs))

	for _, line := range recentLogs {
		s += fmt.Sprintf("%s\n", line)
	}
	return s
}

func (m model) View() string {

	switch m.state {
	case stateConnect:
		return m.ViewConnect()
	case stateGetFamily:
		return m.ViewGetFamily()
	case stateGetType:
		return m.ViewGetType()
	case stateGetProtocol:
		return m.ViewGetProtocol()
	case stateDone:
		return m.ViewDone()
	default:
		return "Unknown state\n"
	}
}
