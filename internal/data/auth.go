package data

import (
	"context"

	"github.com/andrii-lavreniuk/oapi-fiber-example/internal/interfaces"
	"github.com/andrii-lavreniuk/oapi-fiber-example/internal/models"

	"github.com/uptrace/bun"
)

type authRepo struct {
	db *bun.DB
}

func NewAuthRepo(data *Data) interfaces.AuthRepo {
	return &authRepo{
		db: data.db,
	}
}

func (r *authRepo) Exists(ctx context.Context, apikey string) (bool, error) {
	auth := &models.Auth{}

	ok, err := r.db.NewSelect().Model(auth).Where("`api-key` = ?", apikey).Exists(ctx)
	if err != nil {
		return false, err
	}

	return ok, nil
}
