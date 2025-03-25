package main

import (
	"log"
	"net"

	grpc_handler "template-api/internal/transport/grpc"
	proto "template-api/pkg/proto"

	"google.golang.org/grpc"
)

func main() {
	handlers := grpc_handler.NewHandler()
	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	proto.RegisterItemServiceServer(grpcServer, handlers)

	log.Println("gRPC server is listening on :50051")
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
