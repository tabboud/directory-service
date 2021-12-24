package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/palantir/conjure-go-runtime/v2/conjure-go-client/httpclient"
	"github.com/tabboud/directory-service/rpc/authservice"
	authapi "github.com/tabboud/directory-service/rpc/conjure/api/auth"
	"google.golang.org/grpc"
)

func main() {
	addr := flag.String("addr", "localhost:8080", "the address to connect to")
	user := flag.String("username", "john", "Username")
	pass := flag.String("password", "doe", "password")
	serverType := flag.String("type", "twirp", "Server type to connect to ('grpc', 'twirp', 'conjure')")
	flag.Parse()

	switch *serverType {
	case "grpc":
		runGRPCClient(*addr, *user, *pass)
	case "twirp":
		runTwirpClient(getHttpAddr(*addr), *user, *pass)
	case "conjure":
		runConjureClient(getHttpAddr(*addr), *user, *pass)
	default:
		log.Fatalf("Unknown server type: %s, use one of ('grpc', 'twirp', 'conjure')", *serverType)
	}
}

func getHttpAddr(addr string) string {
	if !strings.HasPrefix(addr, "http://") || !strings.HasPrefix(addr, "https://") {
		return "http://" + addr
	}
	return addr
}

func runGRPCClient(addr, user, pass string) {
	// Client for the grpc server
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("[grpc] Failed to connect: %v", err)
	}
	defer func() {
		_ = conn.Close()
	}()

	// contact the server and print the response
	client := authservice.NewAuthServiceV1Client(conn)
	req := &authservice.LoginRequestV1{
		Username: user,
		Password: pass,
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	resp, err := client.Login(ctx, req)
	if err != nil {
		log.Fatalf("[grpc] Failed to login: %v", err)
	}
	logResponse(resp.GetAccessToken(), resp.GetExpiresIn())
}

func runTwirpClient(addr, user, pass string) {
	httpClient := http.DefaultClient
	client := authservice.NewAuthServiceV1JSONClient(addr, httpClient)
	resp, err := client.Login(context.Background(), &authservice.LoginRequestV1{
		Username: user,
		Password: pass,
	})
	if err != nil {
		log.Fatalf("[twirp] Failed to login: %v", err)
	}
	logResponse(resp.GetAccessToken(), resp.GetExpiresIn())
}
func runConjureClient(addr, user, pass string) {
	httpClient, err := httpclient.NewClient(
		httpclient.WithBaseURLs([]string{addr}),
	)
	if err != nil {
		log.Fatalf("[conjure] Failed to create http client")
	}
	client := authapi.NewAuthServiceV1Client(httpClient)
	resp, err := client.Login(context.Background(), authapi.LoginRequestV1{
		Username: user,
		Password: pass,
	})
	if err != nil {
		log.Fatalf("[conjure] Failed to login: %v", err)
	}
	logResponse(resp.AccessToken, int64(resp.ExpiresIn))
}

func logResponse(accessToken string, expiresIn int64) {
	log.Printf("AccessToken: %s, ExpiresIn: %d", accessToken, expiresIn)
}
