package auth

import (
	"context"

	"github.com/palantir/pkg/uuid"
	"github.com/tabboud/directory-service/rpc/authservice"
)

type Service struct {
}

func NewService() *Service {
	return &Service{}
}

func (s *Service) Login(ctx context.Context, req *authservice.LoginRequestV1) (*authservice.LoginResponseV1, error) {
	return &authservice.LoginResponseV1{
		AccessToken: uuid.NewUUID().String(),
		ExpiresIn:   60,
	}, nil
}
