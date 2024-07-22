package data

import (
	"context"
	"database/sql"
	"errors"

	"github.com/andrii-lavreniuk/oapi-fiber-example/internal/interfaces"
	"github.com/andrii-lavreniuk/oapi-fiber-example/internal/models"

	"github.com/uptrace/bun"
)

type usersRepo struct {
	db *bun.DB
}

func NewUsersRepo(data *Data) interfaces.UsersRepo {
	return &usersRepo{
		db: data.db,
	}
}

func (r *usersRepo) GetProfiles(ctx context.Context, limit, offset int32) ([]interfaces.UserProfile, error) {
	var users []models.User

	err := r.db.NewSelect().
		Model(&users).
		Relation("Profile").
		Relation("Data").
		Limit(int(limit)).
		Offset(int(offset)).
		Scan(ctx)
	if err != nil {
		return nil, err
	}

	profiles := make([]interfaces.UserProfile, len(users))

	for i, user := range users {
		profiles[i] = profileFromUser(user)
	}

	return profiles, nil
}

func (r *usersRepo) GetProfile(ctx context.Context, username string) (*interfaces.UserProfile, error) {
	var user models.User

	err := r.db.NewSelect().Model(&user).
		Relation("Profile").
		Relation("Data").
		Where("username = ?", username).
		Scan(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil //nolint:nilnil // expected result
		}

		return nil, err
	}

	profile := profileFromUser(user)

	return &profile, nil
}

func profileFromUser(user models.User) interfaces.UserProfile {
	profile := interfaces.UserProfile{
		ID:       user.ID,
		Username: user.Username,
	}

	if user.Profile != nil {
		profile.FirstName = user.Profile.FirstName
		profile.LastName = user.Profile.LastName
		profile.City = user.Profile.City
	}

	if user.Data != nil {
		profile.School = user.Data.School
	}

	return profile
}
