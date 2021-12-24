package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/tabboud/directory-service/internal/auth"
	"github.com/tabboud/directory-service/internal/token"
	"github.com/tabboud/directory-service/rpc/authservice"
)

func main() {
	addr := flag.String("addr", "localhost:8080", "the address to connect to")
	flag.Parse()

	// Main application handler
	tokenProvider := token.NewUUIDProvider()
	authService := auth.NewService(tokenProvider, 60)
	handler := authservice.NewAuthServiceV1Server(authService)

	srv := &http.Server{
		Addr:    *addr,
		Handler: handler,
	}

	log.Printf("Starting application on %v", *addr)
	if err := srv.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
