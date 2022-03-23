package main

import (
	"log"
	"net"

	"grpc-demo/file-transfer/fileStreaming"

	"google.golang.org/grpc"
)

func main() {
	// Setup listener
	lis, err := net.Listen("tcp", ":9000")
	if err != nil {
		log.Fatalf("Listen: %v", err)
	}

	// Register gRPC service
	grpcServer := grpc.NewServer()
	s := fileStreaming.Server{}
	fileStreaming.RegisterFileUploadServiceServer(grpcServer, &s)

	// Listen for requests
	log.Println("Listening on :9000...")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Serve: %v", err)
	}
}
