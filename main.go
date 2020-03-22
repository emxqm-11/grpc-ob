package main

import (
	"fmt"
	"log"
	"net"

	ag "github.com/emxqm-11/grpc-ob/aggregator"
	"github.com/emxqm-11/grpc-ob/customob"
	"google.golang.org/grpc"
)

// main start a gRPC server and waits for connection
func main() {
	// create a listener on TCP port 7777
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", 7777))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// create a server instance
	s := customob.Server{}

	// create a gRPC server object
	log.Printf("Starting new grpc server...")
	grpcServer := grpc.NewServer()

	// attach the Ping service to the server
	ag.RegisterBankingAggregatorServiceServer(grpcServer, &s)
	// start the server
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %s", err)
	}
	log.Printf("Server started on %s", "port :7777")
}
