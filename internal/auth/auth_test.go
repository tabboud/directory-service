package auth_test

import (
	"context"
	"testing"

	"github.com/tabboud/directory-service/internal/auth"
	"github.com/tabboud/directory-service/rpc/authservice"
)

func Test_Login(t *testing.T) {
	var (
		tokenProvider = staticTokenProvider{token: "token"}
		ttl           = 60
	)
	for _, tc := range []struct {
		name    string
		req     *authservice.LoginRequestV1
		want    *authservice.LoginResponseV1
		wantErr bool
	}{
		{
			name: "valid request",
			req: &authservice.LoginRequestV1{
				Username: "john",
				Password: "password",
			},
			want: &authservice.LoginResponseV1{
				AccessToken: "token",
				ExpiresIn:   int64(ttl),
			},
			wantErr: false,
		},
		{
			name:    "empty request",
			req:     nil,
			want:    nil,
			wantErr: true,
		},
		{
			name: "no username",
			req: &authservice.LoginRequestV1{
				Password: "password",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "no password",
			req: &authservice.LoginRequestV1{
				Username: "john",
			},
			want:    nil,
			wantErr: true,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			s := auth.NewService(tokenProvider, ttl)
			got, err := s.Login(context.Background(), tc.req)
			if (err != nil) != tc.wantErr {
				t.Fatalf("Login() error = %v, wantErr %v", err, tc.wantErr)
			}
			if tc.wantErr {
				// skip response verification
				return
			}

			// verify response
			if tc.want.AccessToken != got.AccessToken {
				t.Fatalf("access tokens are not equal. got: %s, expected: %s",
					got.AccessToken, tc.want.AccessToken)
			}
			if tc.want.ExpiresIn != got.ExpiresIn {
				t.Fatalf("expiresIn are not equal. got: %d, expected: %d",
					got.ExpiresIn, tc.want.ExpiresIn)
			}
		})
	}
}

type staticTokenProvider struct {
	token string
}

func (s staticTokenProvider) GetToken(ctx context.Context) string {
	return s.token
}
