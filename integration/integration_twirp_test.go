package integration

import (
	"context"
	"crypto/tls"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/tabboud/directory-service/internal/auth"
	"github.com/tabboud/directory-service/rpc/authservice"
)

func Test_twirp(t *testing.T) {
	const (
		token     = "test-token"
		expiresIn = 60
	)
	expectedResp := &authservice.LoginResponseV1{
		AccessToken: token,
		ExpiresIn:   expiresIn,
	}

	tokenProvider := hardcodedTokenProvider{token: token}
	authService := auth.NewService(tokenProvider, expiresIn)
	handler := authservice.NewAuthServiceV1Server(authService)
	srv := httptest.NewUnstartedServer(handler)
	srv.TLS = &tls.Config{InsecureSkipVerify: true}
	srv.Start()
	defer srv.Close()

	loginReq := &authservice.LoginRequestV1{
		Username: "test-user",
		Password: "test-password",
	}
	httpClient := http.DefaultClient

	// Protobuf client test
	protoClient := authservice.NewAuthServiceV1ProtobufClient(srv.URL, httpClient)
	resp, err := protoClient.Login(context.Background(), loginReq)
	assertNoError(t, err)
	assertResponse(t, expectedResp, resp)

	// JSON client test
	jsonClient := authservice.NewAuthServiceV1JSONClient(srv.URL, httpClient)
	resp, err = jsonClient.Login(context.Background(), loginReq)
	assertNoError(t, err)
	assertResponse(t, expectedResp, resp)
}
