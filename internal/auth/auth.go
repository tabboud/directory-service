package auth

import (
	"context"

	"github.com/tabboud/directory-service/rpc/authservice"
)

type TokenProvider interface {
	GetToken(context.Context) string
}

type Service struct {
	tokenProvider TokenProvider
	ttl           int

	authservice.UnimplementedAuthServiceV1Server
}

func NewService(tp TokenProvider, ttl int) *Service {
	return &Service{
		tokenProvider: tp,
		ttl:           ttl,
	}
}

func (s *Service) Login(ctx context.Context, req *authservice.LoginRequestV1) (*authservice.LoginResponseV1, error) {
	// TODO(tabboud): validate request and return twirp errors (although that couples us to twirp RPC)
	return &authservice.LoginResponseV1{
		AccessToken: s.tokenProvider.GetToken(ctx),
		ExpiresIn:   int64(s.ttl),
	}, nil
}
