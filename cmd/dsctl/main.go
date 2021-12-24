package main

import (
	"context"
	"flag"
	"log"
	"time"

	"github.com/tabboud/directory-service/rpc/authservice"
	"google.golang.org/grpc"
)

func main() {
	addr := flag.String("addr", "localhost:50051", "the address to connect to")
	user := flag.String("username", "john", "Username")
	pass := flag.String("password", "doe", "password")
	flag.Parse()

	// Client for the grpc server
	conn, err := grpc.Dial(*addr, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer func() {
		_ = conn.Close()
	}()

	// contact the server and print the response
	client := authservice.NewAuthServiceV1Client(conn)
	req := &authservice.LoginRequestV1{
		Username: *user,
		Password: *pass,
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	resp, err := client.Login(ctx, req)
	if err != nil {
		log.Fatalf("Failed to login: %v", err)
	}
	log.Printf("AccessToken: %s, ExpiresIn: %d", resp.GetAccessToken(), resp.GetExpiresIn())
}
