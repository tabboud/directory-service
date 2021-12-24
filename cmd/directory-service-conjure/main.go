package main

import (
	"context"
	"flag"
	"log"

	"github.com/palantir/witchcraft-go-server/v2/wrouter"
	"github.com/palantir/witchcraft-go-server/v2/wrouter/whttprouter"
	"github.com/tabboud/directory-service/rpc/conjure/api/auth"
)

func main() {
	addr := flag.String("addr", "localhost:8080", "the address to run the grpc server")
	flag.Parse()
	_ = addr

	router := wrouter.New(whttprouter.New())
	service := newService()
	if err := auth.RegisterRoutesAuthServiceV1(router, service); err != nil {
		log.Fatalf("Failed to register routes: %v", err)
	}
}

type service struct{}

func newService() *service {
	return &service{}
}

func (s *service) Login(ctx context.Context, requestArg auth.LoginRequestV1) (auth.LoginResponseV1, error) {
	return auth.LoginResponseV1{
		AccessToken: "test-token",
		ExpiresIn:   50,
	}, nil
}
