package interfaces

import (
	"context"
)

type UserProfile struct {
	ID        int32
	Username  string
	FirstName string
	LastName  string
	City      string
	School    string
}

type UsersRepo interface {
	GetProfiles(ctx context.Context, limit, offset int32) ([]UserProfile, error)
	GetProfile(ctx context.Context, username string) (*UserProfile, error)
}
