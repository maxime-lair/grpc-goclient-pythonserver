package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	pb "main/pb_server"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	tea "github.com/charmbracelet/bubbletea"
)

const (
	stateConnect     = 0
	stateGetFamily   = 1
	stateGetType     = 2
	stateGetProtocol = 3
	stateDone        = 4
)

var (
	serverAddr = flag.String("addr", "localhost:50051", "The server address in the format of host:port")
)

// Info used to create the client connection
type connInfo struct {
	opts       []grpc.DialOption
	serverAddr string
}

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
	state    int              // Current state (connect, getXX, done..) - see const for values
	cursor   int              // which to-do list item our cursor is pointing at
	selected map[int]struct{} // which to-do items are selected

	clientChoice clientChoice // all client choices for the socket
	clientEnv    clientEnv    // all client informations
}

func (e errMsg) Error() string { return e.err.Error() }

/*************
Bubbletea part
*************/

func initialModel(conn grpc.ClientConnInterface) model {

	var initModel model
	initModel.state = stateConnect
	/* Connect to server message (loading bar) */
	initModel.clientEnv.logJournal = append(initModel.clientEnv.logJournal, "Client up, proceeding..")

	initModel.clientEnv.logJournal = append(initModel.clientEnv.logJournal, fmt.Sprintf("Connecting to: %s", *serverAddr))

	initModel.clientEnv.logJournal = append(initModel.clientEnv.logJournal, "Connected to server, proceeding..")

	client := pb.NewSocketGuideClient(conn)
	initModel.clientEnv.client = &client
	client_id := &pb.SocketTree{Name: define_client_id()}
	initModel.clientEnv.clientID = client_id

	initModel.clientEnv.logJournal = append(initModel.clientEnv.logJournal, fmt.Sprintf("Created client %p with id %s", &client, client_id.Name))

	return initModel
}

func (m model) Init() tea.Cmd {
	// Just return `nil`, which means "no I/O right now, please."
	return nil
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

/*************
GRPC part
*************/

func socket_get_family_list(clientEnv clientEnv) ([]socketChoice, clientEnv) {

	clientEnv.logJournal = append(clientEnv.logJournal, fmt.Sprintf("[%s][GetFamilyList] Entering.", clientEnv.clientID.Name))

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	socketFamilyStream, req_err := (*clientEnv.client).GetSocketFamilyList(ctx, clientEnv.clientID)
	if req_err != nil {
		clientEnv.err = errMsg{req_err}
		clientEnv.logJournal = append(clientEnv.logJournal, clientEnv.err.Error())
		return nil, clientEnv
	}
	var socketFamilyList []socketChoice
	for {
		family, stream_err := socketFamilyStream.Recv()
		if stream_err == io.EOF {
			break
		}
		if stream_err != nil {
			clientEnv.err = errMsg{req_err}
			clientEnv.logJournal = append(clientEnv.logJournal, clientEnv.err.Error())
			return nil, clientEnv
		}
		//clientEnv.logJournal = append(clientEnv.logJournal, fmt.Sprintf("[%s][GetFamilyList] Received family: %s", clientEnv.clientID.Name, family))
		socketFamilyList = append(socketFamilyList, socketChoice{
			Name:  family.Name,
			Value: family.Value,
		})
	}
	clientEnv.logJournal = append(clientEnv.logJournal, fmt.Sprintf("[%s][GetFamilyList] len=%d cap=%d", clientEnv.clientID.Name, len(socketFamilyList), cap(socketFamilyList)))

	return socketFamilyList, clientEnv
}

func socket_get_type_list(clientEnv clientEnv, clientChoice clientChoice) ([]socketChoice, clientEnv) {

	clientEnv.logJournal = append(clientEnv.logJournal, fmt.Sprintf("[%s][GetTypeList] Entering with family: %d --> %s\n", clientEnv.clientID.Name, clientChoice.selectedFamily.Value, clientChoice.selectedFamily.Name))

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	socketTypeStream, req_err := (*clientEnv.client).GetSocketTypeList(ctx, &pb.SocketFamily{
		Name:     clientChoice.selectedFamily.Name,
		Value:    clientChoice.selectedFamily.Value,
		ClientId: clientEnv.clientID,
	})
	if req_err != nil {
		clientEnv.err = errMsg{req_err}
		clientEnv.logJournal = append(clientEnv.logJournal, clientEnv.err.Error())
		return nil, clientEnv
	}

	var socketTypeList []socketChoice
	for {
		socketType, stream_err := socketTypeStream.Recv()
		if stream_err == io.EOF {
			break
		}
		if stream_err != nil {
			clientEnv.err = errMsg{req_err}
			clientEnv.logJournal = append(clientEnv.logJournal, clientEnv.err.Error())
			return nil, clientEnv
		}
		clientEnv.logJournal = append(clientEnv.logJournal, fmt.Sprintf("[%s][GetTypeList] Received family: %s\n", clientEnv.clientID.Name, socketType))
		socketTypeList = append(socketTypeList, socketChoice{
			Name:  socketType.Name,
			Value: socketType.Value,
		})
	}
	clientEnv.logJournal = append(clientEnv.logJournal, fmt.Sprintf("[%s][GetTypeList] len=%d cap=%d\n", clientEnv.clientID.Name, len(socketTypeList), cap(socketTypeList)))

	return socketTypeList, clientEnv
}

func socket_get_protocol_list(clientEnv clientEnv, clientChoice clientChoice) ([]socketChoice, clientEnv) {

	clientEnv.logJournal = append(clientEnv.logJournal, fmt.Sprintf("[%s][GetProtocolList] Entering with family: [%d] %s -- [%d] %s\n",
		clientEnv.clientID.Name,
		clientChoice.selectedFamily.Value, clientChoice.selectedFamily.Name,
		clientChoice.selectedType.Value, clientChoice.selectedType.Name))

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	socketProtocolStream, req_err := (*clientEnv.client).GetSocketProtocolList(ctx, &pb.SocketTypeAndFamily{
		Family: &pb.SocketFamily{
			Name:  clientChoice.selectedFamily.Name,
			Value: clientChoice.selectedFamily.Value,
		},
		Type: &pb.SocketType{
			Name:  clientChoice.selectedType.Name,
			Value: clientChoice.selectedType.Value,
		},
		ClientId: clientEnv.clientID,
	})
	if req_err != nil {
		clientEnv.err = errMsg{req_err}
		clientEnv.logJournal = append(clientEnv.logJournal, clientEnv.err.Error())
		return nil, clientEnv
	}

	var socketProtocolList []socketChoice
	for {
		socketProtocol, stream_err := socketProtocolStream.Recv()
		if stream_err == io.EOF {
			break
		}
		if stream_err != nil {
			clientEnv.err = errMsg{req_err}
			clientEnv.logJournal = append(clientEnv.logJournal, clientEnv.err.Error())
			return nil, clientEnv
		}
		clientEnv.logJournal = append(clientEnv.logJournal, fmt.Sprintf("[%s][GetProtocolList] Received protocol: %s\n", clientEnv.clientID.Name, socketProtocol))
		socketProtocolList = append(socketProtocolList, socketChoice{
			Name:  socketProtocol.Name,
			Value: socketProtocol.Value,
		})
	}
	clientEnv.logJournal = append(clientEnv.logJournal, fmt.Sprintf("[%s][GetProtocolList] len=%d cap=%d\n", clientEnv.clientID.Name, len(socketProtocolList), cap(socketProtocolList)))

	return socketProtocolList, clientEnv

}

func main() {

	// Retrieve connections options (TLS, etc..)
	flag.Parse()
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))

	conn, err := grpc.Dial(*serverAddr, opts...)
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	defer conn.Close()

	/* Start TUI */
	if err := tea.NewProgram(initialModel(conn)).Start(); err != nil {
		log.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
	/* End of TUI */

	/* Update with family list */

	/* Show family choice, update with type list */
	/*
		log.Printf("[%s] -------------- GetSocketTypeList --------------\n", client_id)
		socketTypeChoice := socket_get_type_list(client_id, socketFamilyChoice, client)
		// Show family/type choice, update with protocol list
		log.Printf("[%s] -------------- GetSocketProtocolList --------------\n", client_id)
		if socketTypeChoice.Name != "" {
			log.Printf("[%s] Socket type choice is not empty, choosing protocol\n", client_id)
			socketTypeAndFamilyChoice := &pb.SocketTypeAndFamily{
				Family:   socketFamilyChoice,
				Type:     socketTypeChoice,
				ClientId: client_id,
			}
			socketProtocolChoice := socket_get_protocol_list(socketTypeAndFamilyChoice, client)

			log.Printf("[%s] Chosen socket: %s - %s - %s", client_id.Name, socketFamilyChoice.Name, socketTypeChoice.Name, socketProtocolChoice.Name)

		} else {
			log.Printf("[%s] Socket type choice is empty, no protocols available\n", client_id)
		}
		// Show family/type/protocol choice
		log.Printf("[%s] Finished requesting socket family, type and protocol\n", client_id)
	*/
}
