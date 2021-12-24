package integration

import (
	"context"
	"testing"

	"github.com/tabboud/directory-service/rpc/authservice"
)

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

type hardcodedTokenProvider struct {
	token string
}

func (h hardcodedTokenProvider) GetToken(context.Context) string {
	return h.token
}
