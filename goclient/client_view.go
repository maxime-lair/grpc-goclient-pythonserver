package main

import (
	"fmt"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/lipgloss"
)

// The header
func (m model) printHeader() string {

	var s string

	{
		// Title
		desc := lipgloss.JoinVertical(lipgloss.Center,
			descStyle.Render("Go client to request sockets through grpc"),
			infoStyle.Render("Built with"+divider+url("GRPC")+divider+url("BubbleTea")+divider+url("Bubbles")+divider+url("LipGloss")),
		)

		row := lipgloss.JoinHorizontal(lipgloss.Top, activeTab.Render(desc))
		s += fmt.Sprintf("%s\n\n", row)
	}

	// Status bar
	{
		status := ""

		switch m.state {
		case stateConnect:
			status = "Request process starting.."
		case stateGetFamily:
			if m.clientChoice.selectedFamily != nil {
				status = fmt.Sprintf("Selected family value : [%d] %s", m.clientChoice.selectedFamily.Value, m.clientChoice.selectedFamily.Name)
			} else {
				status = "No selected family value"
			}
		case stateGetType:
			if m.clientChoice.selectedType != nil {
				status = fmt.Sprintf("Family selected %s, selected type [%d] %s", m.clientChoice.selectedFamily.Name, m.clientChoice.selectedType.Value, m.clientChoice.selectedType.Name)
			} else {
				status = fmt.Sprintf("Family selected %s, no selected type", m.clientChoice.selectedFamily.Name)
			}
		case stateGetProtocol:
			if m.clientChoice.selectedProtocol != nil {
				status = fmt.Sprintf("Family selected %s, type %s, selected protocol [%d] %s", m.clientChoice.selectedFamily.Name, m.clientChoice.selectedType.Name,
					m.clientChoice.selectedProtocol.Value, m.clientChoice.selectedProtocol.Name)
			} else {
				status = fmt.Sprintf("Family selected %s, type %s, no selected protocol", m.clientChoice.selectedFamily.Name, m.clientChoice.selectedType.Name)
			}
		case stateDone:
			status = "Request process done, your final choice will appear below:"
		default:
			status = "Unknown state, exiting.."
		}

		w := lipgloss.Width

		statusKey := statusStyle.Render("Status")
		clientIDKey := clientIDKeyStyle.Render("ClientID")
		clientIDName := clientIDStyle.Render(m.clientEnv.clientID.Name)
		statusVal := statusText.Copy().
			Width(90 - w(statusKey) - w(clientIDKey) - w(clientIDName)).
			Render(status)

		bar := lipgloss.JoinHorizontal(lipgloss.Top,
			statusKey,
			statusVal,
			clientIDKey,
			clientIDName,
		)

		s += statusBarStyle.Width(90).Render(bar) + "\n\n"
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
