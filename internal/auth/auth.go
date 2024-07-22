package auth

import (
	"context"

	"github.com/google/wire"

	"github.com/andrii-lavreniuk/oapi-fiber-example/internal/interfaces"
)

var ProviderSet = wire.NewSet(New)

type Auth struct {
	authRepo interfaces.AuthRepo
}

func New(authRepo interfaces.AuthRepo) interfaces.Auth {
	return &Auth{
		authRepo: authRepo,
	}
}

func (a Auth) ValidateAPIKey(ctx context.Context, apikey string) (bool, error) {
	return a.authRepo.Exists(ctx, apikey)
}
