package main

import (
	"flag"
	"log"
	"net"

	"github.com/tabboud/directory-service/internal/auth"
	"github.com/tabboud/directory-service/internal/token"
	"github.com/tabboud/directory-service/rpc/authservice"
	"google.golang.org/grpc"
)

func main() {
	addr := flag.String("addr", "localhost:8080", "the address to run the grpc server")
	flag.Parse()

	listener, err := net.Listen("tcp", *addr)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	tokenProvider := token.NewUUIDProvider()
	service := auth.NewService(tokenProvider, 60)
	srv := grpc.NewServer()
	authservice.RegisterAuthServiceV1Server(srv, service)
	if err := srv.Serve(listener); err != nil {
		log.Fatalf("Failed to server: %v", err)
	}
}
