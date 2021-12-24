package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/palantir/witchcraft-go-server/v2/wrouter"
	"github.com/palantir/witchcraft-go-server/v2/wrouter/whttprouter"
	"github.com/tabboud/directory-service/internal/auth"
	"github.com/tabboud/directory-service/internal/token"
	authapi "github.com/tabboud/directory-service/rpc/conjure/api/auth"
)

func main() {
	addr := flag.String("addr", "localhost:8080", "the address to run the grpc server")
	flag.Parse()

	router := wrouter.New(whttprouter.New())
	tokenProvider := token.NewUUIDProvider()
	service := auth.NewConjureService(tokenProvider, 60)
	if err := authapi.RegisterRoutesAuthServiceV1(router, service); err != nil {
		log.Fatalf("Failed to register routes: %v", err)
	}

	// Startup the server with the router as the main handler
	srv := &http.Server{
		Addr:    *addr,
		Handler: router,
	}
	log.Printf("Server listening on: %v", *addr)
	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
}
