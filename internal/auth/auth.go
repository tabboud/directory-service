package auth

import (
	"context"

	"github.com/tabboud/directory-service/rpc/authservice"
	"github.com/twitchtv/twirp"
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
	if req == nil {
		return nil, twirp.NewError(twirp.InvalidArgument, "login request cannot be empty")
	}
	if err := validateRequest(req); err != nil {
		return nil, err
	}
	return &authservice.LoginResponseV1{
		AccessToken: s.tokenProvider.GetToken(ctx),
		ExpiresIn:   int64(s.ttl),
	}, nil
}

func validateRequest(req *authservice.LoginRequestV1) error  {
	if req.Username == "" {
		return twirp.NewError(twirp.InvalidArgument, "login username cannot be empty")
	}
	if req.Password == "" {
		return twirp.NewError(twirp.InvalidArgument, "login password cannot be empty")
	}
	return nil
}