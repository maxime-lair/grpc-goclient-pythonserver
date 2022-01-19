package main

import (
	"flag"
	"log"

	pb "github.com/maxime-lair/grpc-goclient-pythonserver/server"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	serverAddr = flag.String("addr", "localhost:50051", "The server address in the format of host:port")
)

func main() {
	flag.Parse()
	var opts []grpc.DialOption

	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))

	log.Printf("Connecting to: %s", *serverAddr)

	conn, err := grpc.Dial(*serverAddr, opts...)
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	log.Printf("Connected to server, proceeding..")

	client := pb.NewSocketGuideClient(conn)

	defer conn.Close()

}
