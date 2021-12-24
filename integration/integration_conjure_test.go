package integration

import (
	"context"
	"crypto/tls"
	"net/http/httptest"
	"testing"

	"github.com/palantir/conjure-go-runtime/v2/conjure-go-client/httpclient"
	"github.com/palantir/witchcraft-go-server/v2/wrouter"
	"github.com/palantir/witchcraft-go-server/v2/wrouter/whttprouter"
	"github.com/tabboud/directory-service/internal/auth"
	authapi "github.com/tabboud/directory-service/rpc/conjure/api/auth"
)

func Test_conjure(t *testing.T) {
	const (
		token     = "test-token"
		expiresIn = 60
	)
	// Server setup/wiring
	router := wrouter.New(whttprouter.New())
	tokenProvider := hardcodedTokenProvider{token: token}
	authService := auth.NewConjureService(tokenProvider, expiresIn)
	if err := authapi.RegisterRoutesAuthServiceV1(router, authService); err != nil {
		t.Fatalf("Failed to register auth service routes: %v", err)
	}
	srv := httptest.NewUnstartedServer(router)
	srv.TLS = &tls.Config{InsecureSkipVerify: true}
	srv.Start()
	defer srv.Close()

	// Client testing
	loginReq := authapi.LoginRequestV1{
		Username: "test-user",
		Password: "test-password",
	}
	expected := &authapi.LoginResponseV1{
		AccessToken: token,
		ExpiresIn:   expiresIn,
	}
	httpClient, err := httpclient.NewClient(
		httpclient.WithBaseURLs([]string{srv.URL}),
		httpclient.WithTLSInsecureSkipVerify(),
	)
	if err != nil {
		t.Fatalf("Failed to create http client: %v", err)
	}

	client := authapi.NewAuthServiceV1Client(httpClient)
	resp, err := client.Login(context.Background(), loginReq)
	assertNoError(t, err)
	if expected.AccessToken != resp.AccessToken {
		t.Fatalf("access tokens are not equal. got: %s, expected: %s",
			resp.AccessToken, expected.AccessToken)
	}
	if expected.ExpiresIn != resp.ExpiresIn {
		t.Fatalf("expiresIn are not equal. got: %d, expected: %d",
			resp.ExpiresIn, expected.ExpiresIn)
	}
}
