package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/palantir/conjure-go-runtime/v2/conjure-go-contract/errors"
	"github.com/palantir/conjure-go-runtime/v2/conjure-go-server/httpserver"
	"github.com/palantir/pkg/bearertoken"
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
	router.AddRouteHandlerMiddleware(getAuthnMiddleware(router.RegisteredRoutes()))

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

func getAuthnMiddleware(registeredRoutes []wrouter.RouteSpec) wrouter.RouteHandlerMiddleware {
	// pre-compute route selection strategy
	for _, route := range registeredRoutes {
		if route.Method == http.MethodGet {
			// do stuff
		}
	}
	return func(rw http.ResponseWriter, r *http.Request, reqVals wrouter.RequestVals, next wrouter.RouteRequestHandler) {
		// TODO(tabboud): Add matchFunc for route selection based on registeredRoutes
		authHeader, err := httpserver.ParseBearerTokenHeader(r)
		if err != nil {
			errors.WriteErrorResponse(rw, errors.WrapWithPermissionDenied(err))
			return
		}
		token := bearertoken.Token(authHeader)
		_ = token

		// Now verify the authHeader
		next(rw, r, reqVals)
	}
}
