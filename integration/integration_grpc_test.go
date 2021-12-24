package integration

import (
	"context"
	"log"
	"net"
	"testing"
	"time"

	"github.com/tabboud/directory-service/internal/auth"
	"github.com/tabboud/directory-service/rpc/authservice"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
)

func Test_gRPC(t *testing.T) {
	var (
		token         = "token"
		ttl           = 60
		tokenProvider = hardcodedTokenProvider{token: token}
		authService   = auth.NewService(tokenProvider, ttl)
	)

	// setup the gRPC server
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	srv, listener := startGRPCServer(authService)
	defer srv.Stop()
	conn, err := grpc.DialContext(ctx, "",
		grpc.WithInsecure(),
		grpc.WithContextDialer(getDialer(listener)))
	if err != nil {
		t.Fatalf("Failed to dial: %v", err)
	}
	defer func() {
		_ = conn.Close()
	}()

	// test/assert
	req := &authservice.LoginRequestV1{
		Username: "user",
		Password: "pass",
	}
	expectedResp := &authservice.LoginResponseV1{
		AccessToken: "token",
		ExpiresIn:   60,
	}
	client := authservice.NewAuthServiceV1Client(conn)
	resp, err := client.Login(context.Background(), req)
	if err != nil {
		t.Fatalf("Failed to Login: %v", err)
	}
	assertResponse(t, expectedResp, resp)
}

func startGRPCServer(as *auth.Service) (*grpc.Server, *bufconn.Listener) {
	listener := bufconn.Listen(1024 * 1024)
	srv := grpc.NewServer()
	authservice.RegisterAuthServiceV1Server(srv, as)
	go func() {
		if err := srv.Serve(listener); err != nil {
			log.Fatalf("Failed to serve: %v", err)
		}
	}()
	return srv, listener
}

func getDialer(listener *bufconn.Listener) func(context.Context, string) (net.Conn, error) {
	return func(ctx context.Context, _ string) (net.Conn, error) {
		return listener.DialContext(ctx)
	}
}
