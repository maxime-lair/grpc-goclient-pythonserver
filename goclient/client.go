package main

import (
	"bufio"
	"context"
	"flag"
	"io"
	"log"
	"math/rand"
	"os"
	"time"

	pb "main/pb_server"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	serverAddr = flag.String("addr", "localhost:50051", "The server address in the format of host:port")
)

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
	log.Printf("Created new client id : %s - %s\n", color_picked, animal_picked)
	return color_picked + "_" + animal_picked
}

func socket_get_family_list(client_id *pb.SocketTree, client pb.SocketGuideClient) *pb.SocketFamily {

	log.Printf("[%s][GetFamilyList] Entering.", client_id.Name)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	socketFamilyList, req_err := client.GetSocketFamilyList(ctx, client_id)
	check(req_err)

	for {
		family, stream_err := socketFamilyList.Recv()
		if stream_err == io.EOF {
			break
		}
		check(stream_err)
		log.Printf("Received family: %s", family)
	}

	var socketFamilyChoice pb.SocketFamily
	return &socketFamilyChoice
}

func socket_get_type_list(client_id *pb.SocketTree, socketFamilyChoice *pb.SocketFamily, client pb.SocketGuideClient) *pb.SocketType {

	var socketTypeChoice pb.SocketType

	return &socketTypeChoice
}

func socket_get_protocol_list(socketTypeAndFamilyChoice *pb.SocketTypeAndFamily, client pb.SocketGuideClient) *pb.SocketProtocol {

	var socketProtocolChoice pb.SocketProtocol

	return &socketProtocolChoice

}

func main() {
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
