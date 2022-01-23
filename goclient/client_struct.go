package main

import pb "main/pb_server"

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
	state  int // Current state (connect, getXX, done..) - see const for values
	cursor int // which to-do list item our cursor is pointing at

	clientChoice clientChoice // all client choices for the socket
	clientEnv    clientEnv    // all client informations
}
