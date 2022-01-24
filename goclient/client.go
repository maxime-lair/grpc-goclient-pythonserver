package main

import (
	"flag"
	"log"
	"os"
	"time"

	pb "main/pb_server"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/spinner"
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

func initialModel(client *pb.SocketGuideClient) model {
	s := spinner.New()
	s.Spinner = spinner.Spinner{
		Frames: []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏", "🌑", "🌒", "🌓", "🌔", "🌕", "🌖", "🌗", "🌘", "⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"},
		FPS:    time.Second / 24, //nolint:gomnd
	}
	s.Style = spinnerStyle
	return model{
		state:    stateConnect,
		help:     help.New(),
		keys:     DefaultKeyMap,
		spinner:  s,
		progress: progress.New(progress.WithDefaultGradient()),
		clientEnv: clientEnv{
			client:   client,
			clientID: &pb.SocketTree{Name: define_client_id()},
		},
	}
}

func (m model) Init() tea.Cmd {
	// Just return `nil`, which means "no I/O right now, please."
	return nil
}

func main() {

	/* Create connection */
	// Retrieve connections options (TLS, etc..)
	flag.Parse()
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))

	conn, err := grpc.Dial(*serverAddr, opts...)
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	defer conn.Close() // This needs to be done in main() as TUI is threaded
	client := pb.NewSocketGuideClient(conn)
	/* Connection created, let's start TUI */

	/* Start TUI */
	if err := tea.NewProgram(initialModel(&client)).Start(); err != nil {
		log.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
	/* End of TUI */
}
