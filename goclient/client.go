package main

import (
	"flag"
	"fmt"
	"log"
	"os"

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

/*************
Bubbletea part
*************/

// For messages that contain errors it's often handy to also implement the
// error interface on the message.
func (e errMsg) Error() string { return e.err.Error() }

func initialModel(conn grpc.ClientConnInterface) model {

	var initModel model
	initModel.state = stateConnect
	/* Connect to server message (loading bar) */
	initModel.clientEnv.logJournal = append(initModel.clientEnv.logJournal, fmt.Sprintf("Connecting to: %s", *serverAddr))

	// Create the client
	client := pb.NewSocketGuideClient(conn)
	initModel.clientEnv.client = &client
	// Define a client ID we will use (random)
	client_id := &pb.SocketTree{Name: define_client_id()}
	initModel.clientEnv.clientID = client_id

	initModel.clientEnv.logJournal = append(initModel.clientEnv.logJournal, fmt.Sprintf("Created client %p with id %s", &client, client_id.Name))

	return initModel
}

func (m model) Init() tea.Cmd {
	// Just return `nil`, which means "no I/O right now, please."
	return nil
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
	defer conn.Close() // This needs to be done in main() as TUI is threaded

	/* Start TUI */
	if err := tea.NewProgram(initialModel(conn)).Start(); err != nil {
		log.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
	/* End of TUI */
}
