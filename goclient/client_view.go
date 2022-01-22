package main

import (
	"fmt"
)

// The header
func (m model) printHeader() string {

	var s string
	if m.clientID == nil {
		s += fmt.Sprintf("------- Client ID undefined %p -------\n", m.client)
	} else {
		s += fmt.Sprintf("------- %s -------\n", m.clientID)
	}
	switch m.state {
	case stateConnect:
		s += "Request process starting..\n\n"
	case stateGetFamily:
		s += "Please select your socket family..\n"
	case stateGetType:
		s += fmt.Sprintf("Requesting socket type list for family %s \n\n", m.selectedFamily.Name)
	case stateGetProtocol:
		s += fmt.Sprintf("Requesting socket protocol list for family %s and type %s:\n\n", m.selectedFamily.Name, m.selectedType.Name)
	case stateDone:
		s += "Request process done, your final choice will appear below:\n\n"
	default:
		s += "Unknown state, exiting..\n\n"
	}

	return s
}

func (m model) printChoices(i int, choiceValue int32, choiceName string) string {

	// Is the cursor pointing at this choice?
	cursor := " " // no cursor
	if m.cursor == i {
		cursor = ">" // cursor!
	}

	// Is this choice selected?
	checked := " " // not selected
	if _, ok := m.selected[i]; ok {
		checked = "x" // selected!
	}

	// Render the row
	return fmt.Sprintf("%s [%s] %d - %s \n", cursor, checked, choiceValue, choiceName)

}

func (m model) printFooter() string {

	// The footer
	var s string
	s += "\nPress <space> to add a value\nPress <enter> to validate\n"
	s += "Press q to quit.\n"
	s += "Only one selection at a time possible.\n"

	// print log journal
	s += "------ logs ------\n"
	for _, line := range m.logJournal {
		s += fmt.Sprintf("%s\n", line)
	}

	return s
}

func (m model) ViewConnect() string {
	var s string
	s += m.printHeader()
	// TODO add loading bar
	s += fmt.Sprintf("Press enter to start %s..\n", m.connInfo.serverAddr)
	s += m.printFooter()
	return s
}

func (m model) ViewGetFamily() string {
	var s string
	s += m.printHeader()

	// Iterate over our choices - we need to instanciate them first
	for i, choice := range m.socketFamilyChoices {
		// print choices
		s += m.printChoices(i, choice.Value, choice.Name)
	}

	s += m.printFooter()
	return s
}

func (m model) ViewGetType() string {
	var s string
	s += m.printHeader()
	// Iterate over our choices
	for i, choice := range *m.socketTypeChoices {
		// print choices
		s += m.printChoices(i, choice.Value, choice.Name)
	}

	s += m.printFooter()
	return "getType\n"
}

func (m model) ViewGetProtocol() string {
	var s string
	s += m.printHeader()
	// Iterate over our choices
	for i, choice := range *m.socketProtocolChoices {
		// print choices
		s += m.printChoices(i, choice.Value, choice.Name)
	}
	s += m.printFooter()
	return "getProtocol\n"
}

func (m model) ViewDone() string {
	var s string
	s += m.printHeader()

	s += "You choose the following parameters for your socket:\n"
	s += fmt.Sprintf("\t-> Family: %d - %s\n", m.selectedFamily.Value, m.selectedFamily.Name)
	s += fmt.Sprintf("\t--> Type: %d - %s\n", m.selectedType.Value, m.selectedType.Name)
	s += fmt.Sprintf("\t---> Protocol: %d - %s\n", m.selectedProtocol.Value, m.selectedProtocol.Name)

	s += m.printFooter()
	return "done\n"
}