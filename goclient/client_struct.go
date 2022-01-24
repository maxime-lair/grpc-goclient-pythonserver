package main

import (
	pb "main/pb_server"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/spinner"
)

// Info for family/type/protocol, always a string and value
// We have to create it like this since GRPC struct are mutex protected and we can not loop over them
type socketChoice struct {
	Name  string
	Value int32
}

type errMsg struct{ err error }

// Struct for error/log handling when requesting to client
type clientEnv struct {
	clientID   *pb.SocketTree        // client ID for the transaction
	client     *pb.SocketGuideClient // client
	logJournal []string              // log journal
	err        errMsg                // possible error message
}

// All client choices made
type clientChoice struct {
	socketChoicesList []socketChoice // socket family choices list
	selectedFamily    *socketChoice  // selected family
	selectedType      *socketChoice  // selected type
	selectedProtocol  *socketChoice  // selected protocol
}

// TUI model used to print and show informations
type model struct {
	state   int // Current state (connect, getXX, done..) - see const for values
	cursor  int // which to-do list item our cursor is pointing at
	help    help.Model
	keys    keyMap
	spinner spinner.Model

	clientChoice clientChoice // all client choices for the socket
	clientEnv    clientEnv    // all client informations
}

// Bubbles keymap
type keyMap struct {
	Up    key.Binding
	Down  key.Binding
	Enter key.Binding
	Space key.Binding
	Quit  key.Binding
}
