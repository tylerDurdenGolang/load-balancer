package main

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"net/http"

	"log"
	"net"
	"template-api/internal/repository"
	"template-api/internal/service"
	http_handler "template-api/internal/transport/http"

	grpc_handler "template-api/internal/transport/grpc"
	proto "template-api/pkg/proto"

	"google.golang.org/grpc"
)

func main() {
	dbPool, err := pgxpool.New(context.Background(), "your_connection_string")
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer dbPool.Close()

	repo := repository.NewProductRepository(dbPool)
	s := service.NewProductService(repo)
	handlers := grpc_handler.NewHandler(s)

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

	httpHandlers := http_handler.NewHandler(s)

	httpServer := &http.Server{
		Addr:    ":8080",
		Handler: httpHandlers.InitRoutes(),
	}

	log.Println("HTTP server is listening on :8080")
	if err := httpServer.ListenAndServe(); err != nil {
		log.Fatalf("failed to serve HTTP: %v", err)
	}
}
