package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/tabboud/directory-service/internal/auth"
	"github.com/tabboud/directory-service/internal/server"
)

func main() {
	var (
		addr = "localhost:8080"
	)

	// Main application handler
	authService := auth.NewService()
	handler := server.New(authService)

	srv := &http.Server{
		Addr:    addr,
		Handler: handler,
	}

	// Ensure all open connections are killed before shutting down
	idleConnsClosed := make(chan struct{})
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt)
		<-sigint

		// We received an interrupt signal, shut down.
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
		defer cancel()
		if err := srv.Shutdown(ctx); err != nil {
			log.Fatal(err)
		}
		close(idleConnsClosed)
	}()

	log.Printf("Starting application at port %v", addr)
	if err := srv.ListenAndServe(); err != nil {
		log.Fatal(err)
	}

	<-idleConnsClosed
}
