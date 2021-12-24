package server

import (
	"log"

	"github.com/tabboud/directory-service/rpc/authservice"
)

func New(as authservice.AuthServiceV1) authservice.TwirpServer {
	server := authservice.NewAuthServiceV1Server(as)
	log.Printf("Server pathPrefix: %s", server.PathPrefix())
	return server
}
