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

type conn_info struct {
	opts       []grpc.DialOption
	serverAddr string
}

type socket_choice struct {
	Name  string
	Value int32
}

type model struct {
	choices  []string         // items on the to-do list
	cursor   int              // which to-do list item our cursor is pointing at
	selected map[int]struct{} // which to-do items are selected

	socketFamilyChoices   []socket_choice      // socket family choices list
	socketTypeChoices     *[]socket_choice     // socket type choices list
	socketProtocolChoices *[]socket_choice     // socket protocol choices list
	selectedFamily        *socket_choice       // selected family
	selectedType          *socket_choice       // selected type
	selectedProtocol      *socket_choice       // selected protocol
	clientID              *pb.SocketTree       // client id for the transaction
	connInfo              conn_info            // conn info
	client                pb.SocketGuideClient // client

	logJournal []string // log journal
	err        errMsg   // possible error message
	state      int      // Current state (connect, getXX, done..) - see const for values
}

type errMsg struct{ err error }

func (e errMsg) Error() string { return e.err.Error() }

/*************
Bubbletea part
*************/

func initialModel(connInfo conn_info) model {

	var initModel model
	initModel.state = stateConnect
	/* Connect to server message (loading bar) */
	initModel.logJournal = append(initModel.logJournal, "Client up, proceeding..")

	initModel.logJournal = append(initModel.logJournal, "Connecting to: "+connInfo.serverAddr)

	conn, err := grpc.Dial(connInfo.serverAddr, connInfo.opts...)
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	defer conn.Close()
	initModel.logJournal = append(initModel.logJournal, "Connected to server, proceeding..")
	initModel.connInfo = connInfo

	client := pb.NewSocketGuideClient(conn)
	initModel.client = client

	client_id := &pb.SocketTree{Name: define_client_id()}

	initModel.logJournal = append(initModel.logJournal, fmt.Sprintf("Created client %p with id %s", &client, client_id.Name))

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

// TODO Factorize get part
func socket_get_family_list(client_id *pb.SocketTree, client pb.SocketGuideClient, logJournal []string) ([]socket_choice, []string) {

	logJournal = append(logJournal, fmt.Sprintf("[%s][GetFamilyList] Entering.", client_id.Name))

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	socketFamilyStream, req_err := client.GetSocketFamilyList(ctx, client_id)
	check(req_err)

	var socketFamilyList []socket_choice
	for {
		family, stream_err := socketFamilyStream.Recv()
		if stream_err == io.EOF {
			break
		}
		check(stream_err)
		logJournal = append(logJournal, fmt.Sprintf("[%s][GetFamilyList] Received family: %s\n", client_id.Name, family))
		socketFamilyList = append(socketFamilyList, socket_choice{
			Name:  family.Name,
			Value: family.Value,
		})
	}
	logJournal = append(logJournal, fmt.Sprintf("[%s][GetFamilyList] len=%d cap=%d\n", client_id.Name, len(socketFamilyList), cap(socketFamilyList)))

	return socketFamilyList, logJournal
}

func socket_get_type_list(client_id *pb.SocketTree, socketFamilyChoice *pb.SocketFamily, client pb.SocketGuideClient) *pb.SocketType {

	log.Printf("[%s][GetTypeList] Entering with family: %d --> %s\n", client_id.Name, socketFamilyChoice.Value, socketFamilyChoice.Name)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	socketTypeStream, req_err := client.GetSocketTypeList(ctx, socketFamilyChoice)
	check(req_err)

	var socketTypeList []pb.SocketType
	for {
		socketType, stream_err := socketTypeStream.Recv()
		if stream_err == io.EOF {
			break
		}
		check(stream_err)
		log.Printf("[%s][GetTypeList] Received family: %s\n", client_id.Name, socketType)
		socketTypeList = append(socketTypeList, pb.SocketType{
			Name:     socketType.Name,
			Value:    socketType.Value,
			ClientId: client_id})
	}
	log.Printf("[%s][GetTypeList] len=%d cap=%d\n", client_id.Name, len(socketTypeList), cap(socketTypeList))

	// TUI: Ask for type choice
	return &socketTypeList[0]
}

func socket_get_protocol_list(socketTypeAndFamilyChoice *pb.SocketTypeAndFamily, client pb.SocketGuideClient) *pb.SocketProtocol {

	client_id := socketTypeAndFamilyChoice.ClientId.Name
	log.Printf("[%s][GetProtocolList] Entering with family: [%d] %s -- [%d] %s\n",
		client_id,
		socketTypeAndFamilyChoice.Family.Value, socketTypeAndFamilyChoice.Family.Name,
		socketTypeAndFamilyChoice.Type.Value, socketTypeAndFamilyChoice.Type.Name)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	socketProtocolStream, req_err := client.GetSocketProtocolList(ctx, socketTypeAndFamilyChoice)
	check(req_err)

	var socketProtocolList []pb.SocketProtocol
	for {
		socketProtocol, stream_err := socketProtocolStream.Recv()
		if stream_err == io.EOF {
			break
		}
		check(stream_err)
		log.Printf("[%s][GetProtocolList] Received protocol: %s\n", client_id, socketProtocol)
		socketProtocolList = append(socketProtocolList, pb.SocketProtocol{
			Name:     socketProtocol.Name,
			Value:    socketProtocol.Value,
			ClientId: socketProtocol.ClientId})
	}
	log.Printf("[%s][GetProtocolList] len=%d cap=%d\n", client_id, len(socketProtocolList), cap(socketProtocolList))

	//TUI: Ask for protocol choice
	return &socketProtocolList[0]

}

func main() {

	// Retrieve connections options (TLS, etc..)
	flag.Parse()
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))

	/* Start TUI */
	if err := tea.NewProgram(initialModel(conn_info{serverAddr: *serverAddr, opts: opts})).Start(); err != nil {
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
