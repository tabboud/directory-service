package auth_test

import (
	"context"
	"testing"

	"github.com/tabboud/directory-service/internal/auth"
	"github.com/tabboud/directory-service/rpc/authservice"
)

func TestServicess(t *testing.T) {
	type args struct {
		tp  auth.TokenProvider
		ttl int
	}
	for _, tc := range []struct {
		name    string
		args    args
		req     *authservice.LoginRequestV1
		want    *authservice.LoginResponseV1
		wantErr bool
	}{
		{
			name: "valid request",
			args: args{
				tp: tokenProviderFunc(func() string {
					return "token"
				}),
				ttl: 60,
			},
			req: &authservice.LoginRequestV1{
				Username: "john",
				Password: "password",
			},
			want: &authservice.LoginResponseV1{
				AccessToken: "token",
				ExpiresIn:   60,
			},
			wantErr: false,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			s := auth.NewService(tc.args.tp, tc.args.ttl)
			got, err := s.Login(context.Background(), tc.req)
			if (err != nil) != tc.wantErr {
				t.Fatalf("Login() error = %v, wantErr %v", err, tc.wantErr)
			}
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

type tokenProviderFunc func() string

func (t tokenProviderFunc) GetToken(ctx context.Context) string {
	return t()
}
