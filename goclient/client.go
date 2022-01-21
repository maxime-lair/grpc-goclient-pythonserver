package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"math/rand"
	"os"
	"time"

	pb "main/pb_server"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	tea "github.com/charmbracelet/bubbletea"
)

var (
	serverAddr = flag.String("addr", "localhost:50051", "The server address in the format of host:port")
)

type model struct {
	choices  []string         // items on the to-do list
	cursor   int              // which to-do list item our cursor is pointing at
	selected map[int]struct{} // which to-do items are selected
}

/*************
Bubbletea part
*************/

func initialModel() model {
	return model{
		// Our shopping list is a grocery list
		choices: []string{"Buy carrots", "Buy celery", "Buy kohlrabi"},

		// A map which indicates which choices are selected. We're using
		// the  map like a mathematical set. The keys refer to the indexes
		// of the `choices` slice, above.
		selected: make(map[int]struct{}),
	}
}

func (m model) Init() tea.Cmd {
	// Just return `nil`, which means "no I/O right now, please."
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	// Is it a key press?
	case tea.KeyMsg:

		// Cool, what was the actual key pressed?
		switch msg.String() {

		// These keys should exit the program.
		case "ctrl+c", "q":
			return m, tea.Quit

		// The "up" and "k" keys move the cursor up
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}

		// The "down" and "j" keys move the cursor down
		case "down", "j":
			if m.cursor < len(m.choices)-1 {
				m.cursor++
			}

		// The "enter" key and the spacebar (a literal space) toggle
		// the selected state for the item that the cursor is pointing at.
		case "enter", " ":
			_, ok := m.selected[m.cursor]
			if ok {
				delete(m.selected, m.cursor)
			} else {
				m.selected[m.cursor] = struct{}{}
			}
		}
	}

	// Return the updated model to the Bubble Tea runtime for processing.
	// Note that we're not returning a command.
	return m, nil
}

func (m model) View() string {
	// The header
	s := "What should we buy at the market?\n\n"

	// Iterate over our choices
	for i, choice := range m.choices {

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
		s += fmt.Sprintf("%s [%s] %s\n", cursor, checked, choice)
	}

	// The footer
	s += "\nPress q to quit.\n"

	// Send the UI for rendering
	return s
}

/*************
GRPC part
*************/

func check(e error) {
	if e != nil {
		log.Fatal(e)
	}
}

func get_random_line(open_file *os.File) string {

	// Load it all into memory,
	// a better way would be to have the same byte count on each line and just get a multiple
	var lines []string
	scanner := bufio.NewScanner(open_file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	random_seed := rand.New(rand.NewSource(time.Now().UnixNano()))

	line_number := random_seed.Intn(len(lines))

	log.Printf("Picking line %s at line %d among total %d\n", lines[line_number], line_number, len(lines))

	return lines[line_number]
}

func define_client_id() string {
	// Get random line from a color wordlist
	color_file, color_err := os.Open("../wordlist/color.txt")
	check(color_err)
	defer color_file.Close()
	color_picked := get_random_line(color_file)
	// Get random line from an animal wordlist
	animal_file, animal_err := os.Open("../wordlist/animal.txt")
	check(animal_err)
	defer animal_file.Close()
	animal_picked := get_random_line(animal_file)

	// return both string concat
	log.Printf("Created new client id : %s_%s\n", color_picked, animal_picked)
	return color_picked + "_" + animal_picked
}

func socket_get_family_list(client_id *pb.SocketTree, client pb.SocketGuideClient) *pb.SocketFamily {

	log.Printf("[%s][GetFamilyList] Entering.", client_id.Name)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	socketFamilyStream, req_err := client.GetSocketFamilyList(ctx, client_id)
	check(req_err)

	var socketFamilyList []pb.SocketFamily
	for {
		family, stream_err := socketFamilyStream.Recv()
		if stream_err == io.EOF {
			break
		}
		check(stream_err)
		log.Printf("[%s][GetFamilyList] Received family: %s\n", client_id.Name, family)
		socketFamilyList = append(socketFamilyList, pb.SocketFamily{
			Name:     family.Name,
			Value:    family.Value,
			ClientId: client_id})
	}
	log.Printf("[%s][GetFamilyList] len=%d cap=%d\n", client_id.Name, len(socketFamilyList), cap(socketFamilyList))

	// TUI: ask for family choice
	return &socketFamilyList[1]
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

	p := tea.NewProgram(initialModel())
	if err := p.Start(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}

	flag.Parse()
	var opts []grpc.DialOption

	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))

	log.Printf("Connecting to: %s", *serverAddr)

	conn, err := grpc.Dial(*serverAddr, opts...)
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	defer conn.Close()
	log.Printf("Connected to server, proceeding..")

	client := pb.NewSocketGuideClient(conn)

	client_id := &pb.SocketTree{Name: define_client_id()}

	log.Printf("Created client %v with id %s", &client, client_id.Name)

	log.Printf("-------------- SendSocketTree --------------\n")
	socketFamilyChoice := socket_get_family_list(client_id, client)
	log.Printf("-------------- GetSocketTypeList --------------\n")
	socketTypeChoice := socket_get_type_list(client_id, socketFamilyChoice, client)
	log.Printf("-------------- GetSocketProtocolList --------------\n")
	if socketTypeChoice.Name != "" {
		log.Printf("Socket type choice is not empty, choosing protocol\n")
		socketTypeAndFamilyChoice := &pb.SocketTypeAndFamily{
			Family:   socketFamilyChoice,
			Type:     socketTypeChoice,
			ClientId: client_id,
		}
		socketProtocolChoice := socket_get_protocol_list(socketTypeAndFamilyChoice, client)

		log.Printf("[%s] Chosen socket: %s - %s - %s", client_id.Name, socketFamilyChoice.Name, socketTypeChoice.Name, socketProtocolChoice.Name)

	} else {
		log.Printf("Socket type choice is empty, no protocols available\n")
	}
	log.Printf("Finished requesting socket family, type and protocol\n")
}
