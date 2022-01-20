package main

import (
	"bufio"
	"flag"
	"log"
	"math/rand"
	"os"

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

	random_line := ""

	var lines []string
	scanner := bufio.NewScanner(open_file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	line_number := rand.Intn(len(lines))

	log.Printf("Picking line random %d among total %d\n", line_number, len(lines))

	return random_line
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

	log.Printf("Created client %v with id %s", client, client_id)

}
