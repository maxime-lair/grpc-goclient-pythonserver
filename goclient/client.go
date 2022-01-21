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

var (
	serverAddr = flag.String("addr", "localhost:50051", "The server address in the format of host:port")
)

/*************
GRPC part
*************/

// TODO Factorize get part
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

	log.Printf("[%s] -------------- SendSocketTree --------------\n", client_id)
	socketFamilyChoice := socket_get_family_list(client_id, client)
	log.Printf("[%s] -------------- GetSocketTypeList --------------\n", client_id)
	socketTypeChoice := socket_get_type_list(client_id, socketFamilyChoice, client)
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
	log.Printf("[%s] Finished requesting socket family, type and protocol\n", client_id)
}
