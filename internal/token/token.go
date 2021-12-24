package token

import (
	"context"

	"github.com/palantir/pkg/uuid"
)

type UUIDProvider struct{}

func NewUUIDProvider() *UUIDProvider {
	return &UUIDProvider{}
}

func (u *UUIDProvider) GetToken(ctx context.Context) string {
	return uuid.NewUUID().String()
}
