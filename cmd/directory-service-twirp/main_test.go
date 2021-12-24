package main_test

import (
	"context"
	"crypto/tls"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/tabboud/directory-service/rpc/authservice"
)

func Test_gRPC(t *testing.T) {
	const (
		token     = "test-token"
		expiresIn = 60
	)
	expectedResp := &authservice.LoginResponseV1{
		AccessToken: token,
		ExpiresIn:   expiresIn,
	}

	authService := &mockAuthServiceV1{
		login: func() (*authservice.LoginResponseV1, error) {
			return &authservice.LoginResponseV1{
				AccessToken: token,
				ExpiresIn:   expiresIn,
			}, nil
		},
	}
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

func assertNoError(t *testing.T, err error) {
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func assertResponse(t *testing.T, expected, got *authservice.LoginResponseV1) {
	if expected.AccessToken != got.AccessToken {
		t.Fatalf("access tokens are not equal. got: %s, expected: %s",
			got.AccessToken, expected.AccessToken)
	}
	if expected.ExpiresIn != got.ExpiresIn {
		t.Fatalf("expiresIn are not equal. got: %d, expected: %d",
			got.ExpiresIn, expected.ExpiresIn)
	}
}

type mockAuthServiceV1 struct {
	login func() (*authservice.LoginResponseV1, error)
}

func (m *mockAuthServiceV1) Login(ctx context.Context, req *authservice.LoginRequestV1) (*authservice.LoginResponseV1, error) {
	return m.login()
}
