package handlers

import (
	"context"

	"github.com/google/wire"
	"github.com/sirupsen/logrus"

	api "github.com/andrii-lavreniuk/oapi-fiber-example/gen/api/profiles"
	"github.com/andrii-lavreniuk/oapi-fiber-example/internal/interfaces"
)

var (
	ProviderSet = wire.NewSet(NewProfilesHandler)

	// Compile-time proof of interface implementation.
	_ api.StrictServerInterface = (*ProfilesHandler)(nil)
)

type ProfilesHandler struct {
	usersRepo interfaces.UsersRepo
	lg        *logrus.Entry
}

func NewProfilesHandler(usersRepo interfaces.UsersRepo, lg *logrus.Logger) *ProfilesHandler {
	return &ProfilesHandler{
		usersRepo: usersRepo,
		lg:        lg.WithField("handler", "profiles"),
	}
}

func (h *ProfilesHandler) GetProfiles(
	ctx context.Context, req api.GetProfilesRequestObject,
) (api.GetProfilesResponseObject, error) {
	lg := h.lg.WithContext(ctx)

	var profiles []interfaces.UserProfile

	if req.Params.Username != nil { //nolint:nestif // readability
		profile, err := h.usersRepo.GetProfile(ctx, *req.Params.Username)
		if err != nil {
			lg.WithError(err).Error("failed to get profile")

			return api.GetProfiles500JSONResponse{
				Message: "Failed to get profile",
				Reason:  err.Error(),
			}, nil
		}

		if profile != nil {
			profiles = append(profiles, *profile)
		}
	} else {
		var (
			limit  int32 = 10
			offset int32

			err error
		)

		if req.Params.Limit != nil {
			limit = *req.Params.Limit
		}

		if req.Params.Offset != nil {
			offset = *req.Params.Offset
		}

		profiles, err = h.usersRepo.GetProfiles(ctx, limit, offset)
		if err != nil {
			lg.WithError(err).Error("failed to get profiles")

			return api.GetProfiles500JSONResponse{
				Message: "Failed to get profiles",
				Reason:  err.Error(),
			}, nil
		}
	}

	response := api.GetProfiles200JSONResponse{
		Data: make([]api.UserProfile, len(profiles)),
	}

	for i, profile := range profiles {
		response.Data[i] = api.UserProfile{
			ID:        profile.ID,
			Username:  profile.Username,
			FirstName: profile.FirstName,
			LastName:  profile.LastName,
			City:      profile.City,
			School:    profile.School,
		}
	}

	return response, nil
}
