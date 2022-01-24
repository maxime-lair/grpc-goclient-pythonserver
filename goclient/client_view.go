package main

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/lipgloss"
)

// The header
func (m model) printHeader() string {

	var s string

	{
		// Tabs
		row := lipgloss.JoinHorizontal(
			lipgloss.Top,
			activeTab.Render("GO Client"),
		)
		gap := tabGap.Render(strings.Repeat("", max(0, 96-lipgloss.Width(row)-2)))
		row = lipgloss.JoinHorizontal(lipgloss.Bottom, row, gap)
		s += fmt.Sprintf("%s\n\n", row)

		// Title
		desc := lipgloss.JoinVertical(lipgloss.Center,
			descStyle.Render("Go client to request sockets through grpc"),
			infoStyle.Render("Built with"+divider+url("GRPC")+divider+url("BubbleTea")+divider+url("Bubbles")+divider+url("LipGloss")),
		)

		row = lipgloss.JoinHorizontal(lipgloss.Top, activeTab.Render(desc), gap)
		s += fmt.Sprintf("%s\n\n", row)
	}

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
		s += "Please select your socket family (only one selection possible)\n"
	case stateGetType:
		s += fmt.Sprintf("Requesting socket type list for family %s \n\n", m.clientChoice.selectedFamily.Name)
		s += "Please select your socket type (only one selection possible)\n"
	case stateGetProtocol:
		s += fmt.Sprintf("Requesting socket protocol list for family %s and type %s:\n\n", m.clientChoice.selectedFamily.Name, m.clientChoice.selectedType.Name)
		s += "Please select your socket protocol (only one selection possible)\n"
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

func (m model) printHelp() string {

	helpList := [][]key.Binding{
		{DefaultKeyMap.Up},
		{DefaultKeyMap.Down},
		{DefaultKeyMap.Space},
		{DefaultKeyMap.Enter},
		{DefaultKeyMap.Quit},
	}

	var s string
	s += "\n"
	// print help, idk why fullHelpView does not work here, so had to do it dirty
	for _, group := range helpList {
		s += fmt.Sprintf("%s\n", m.help.ShortHelpView(group))
	}

	return s
}

func (m model) printLogs() string {
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

func (m model) ViewConnect() string {
	var s string
	s += m.printHeader()
	// TODO add loading bar
	s += "Press enter to start ..\n"
	s += m.printHelp()
	s += m.printLogs()
	return s
}

func (m model) ViewGetFamily() string {
	var s string
	s += m.printHeader()

	// Iterate over our choices
	for i, choice := range m.clientChoice.socketChoicesList {
		s += m.printChoices(i, m.clientChoice.selectedFamily, choice)
	}

	s += m.printHelp()
	s += m.printLogs()
	return s
}

func (m model) ViewGetType() string {
	var s string
	s += m.printHeader()
	// Iterate over our choices
	for i, choice := range m.clientChoice.socketChoicesList {
		s += m.printChoices(i, m.clientChoice.selectedType, choice)
	}

	s += m.printHelp()
	s += m.printLogs()
	return s
}

func (m model) ViewGetProtocol() string {
	var s string
	s += m.printHeader()
	// Iterate over our choices
	for i, choice := range m.clientChoice.socketChoicesList {
		s += m.printChoices(i, m.clientChoice.selectedProtocol, choice)
	}
	s += m.printHelp()
	s += m.printLogs()
	return s
}

func (m model) ViewDone() string {
	var s string
	s += m.printHeader()

	s += "You choose the following parameters for your socket:\n"
	s += fmt.Sprintf("\t-> Family: %d - %s\n", m.clientChoice.selectedFamily.Value, m.clientChoice.selectedFamily.Name)
	s += fmt.Sprintf("\t--> Type: %d - %s\n", m.clientChoice.selectedType.Value, m.clientChoice.selectedType.Name)
	s += fmt.Sprintf("\t---> Protocol: %d - %s\n", m.clientChoice.selectedProtocol.Value, m.clientChoice.selectedProtocol.Name)

	s += m.printHelp()
	s += m.printLogs()
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
