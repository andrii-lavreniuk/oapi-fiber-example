package handlers_test

import (
	"context"
	"errors"
	"testing"

	"github.com/aws/smithy-go/ptr"
	"github.com/go-faker/faker/v4"
	"github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/suite"

	api "github.com/andrii-lavreniuk/oapi-fiber-example/gen/api/profiles"
	mocks "github.com/andrii-lavreniuk/oapi-fiber-example/gen/mocks/interfaces"
	"github.com/andrii-lavreniuk/oapi-fiber-example/internal/handlers"
	"github.com/andrii-lavreniuk/oapi-fiber-example/internal/interfaces"
)

func TestProfilesSuite(t *testing.T) {
	suite.Run(t, new(ProfilesTestSuite))
}

type ProfilesTestSuite struct {
	suite.Suite
}

func (s *ProfilesTestSuite) TestGetProfiles() {
	ctx := context.Background()

	s.Run("GetList", func() {
		s.Run("Success", func() {
			tc := []struct {
				title     string
				reqLimit  *int32
				reqOffset *int32
				expLimit  int32
				expOffset int32
			}{
				{
					title:    "Default",
					expLimit: 10,
				},
				{
					title:     "Pagination",
					reqLimit:  ptr.Int32(20),
					reqOffset: ptr.Int32(10),
					expLimit:  20,
					expOffset: 10,
				},
			}

			for _, tt := range tc {
				s.Run(tt.title, func() {
					expectedProfiles := []interfaces.UserProfile{
						{
							ID:        1,
							Username:  faker.Username(),
							FirstName: faker.FirstName(),
							LastName:  faker.LastName(),
							City:      faker.GetRealAddress().City,
							School:    faker.Word(),
						},
						{
							ID:        2,
							Username:  faker.Username(),
							FirstName: faker.FirstName(),
							LastName:  faker.LastName(),
							City:      faker.GetRealAddress().City,
							School:    faker.Word(),
						},
					}

					usersRepo := mocks.NewMockUsersRepo(s.T())
					usersRepo.On("GetProfiles", ctx, tt.expLimit, tt.expOffset).Return(expectedProfiles, nil)

					logger, hook := test.NewNullLogger()
					logger.SetLevel(logrus.DebugLevel)

					profiles := handlers.NewProfilesHandler(usersRepo, logger)

					result, err := profiles.GetProfiles(ctx, api.GetProfilesRequestObject{
						Params: api.GetProfilesParams{
							Limit:  tt.reqLimit,
							Offset: tt.reqOffset,
						},
					})
					s.Require().NoError(err)

					s.Require().IsType(api.GetProfiles200JSONResponse{}, result)

					response, _ := result.(api.GetProfiles200JSONResponse)

					s.Require().Equal(len(expectedProfiles), len(response.Data))

					for i, u := range expectedProfiles {
						s.Equal(u.ID, response.Data[i].ID)
						s.Equal(u.Username, response.Data[i].Username)
						s.Equal(u.FirstName, response.Data[i].FirstName)
						s.Equal(u.LastName, response.Data[i].LastName)
						s.Equal(u.City, response.Data[i].City)
						s.Equal(u.School, response.Data[i].School)
					}

					s.Nil(hook.LastEntry())
				})
			}
		})

		s.Run("Fail", func() {
			dbError := errors.New("db error")

			usersRepo := mocks.NewMockUsersRepo(s.T())
			usersRepo.On("GetProfiles", ctx, int32(10), int32(0)).Return(nil, dbError)

			logger, hook := test.NewNullLogger()
			logger.SetLevel(logrus.DebugLevel)

			profiles := handlers.NewProfilesHandler(usersRepo, logger)

			result, err := profiles.GetProfiles(ctx, api.GetProfilesRequestObject{})
			s.Require().NoError(err)

			s.Require().IsType(api.GetProfiles500JSONResponse{}, result)

			response, _ := result.(api.GetProfiles500JSONResponse)

			s.Equal("Failed to get profiles", response.Message)
			s.Equal(dbError.Error(), response.Reason)

			s.Require().NotNil(hook.LastEntry())
			s.Require().NotEmpty(hook.LastEntry().Data)
			s.Equal("profiles", hook.LastEntry().Data["handler"])
			s.Equal(dbError, hook.LastEntry().Data["error"])
			s.Equal("failed to get profiles", hook.LastEntry().Message)
		})
	})

	s.Run("GetByUsername", func() {
		s.Run("Success", func() {
			tc := []struct {
				title    string
				username string
				profile  *interfaces.UserProfile
				response api.GetProfiles200JSONResponse
			}{
				{
					title:    "Default",
					username: "username",
					profile: &interfaces.UserProfile{
						ID:        1,
						Username:  "username",
						FirstName: "firstName",
						LastName:  "lastName",
						City:      "city",
						School:    "school",
					},
					response: api.GetProfiles200JSONResponse{
						Data: []api.UserProfile{
							{
								ID:        1,
								Username:  "username",
								FirstName: "firstName",
								LastName:  "lastName",
								City:      "city",
								School:    "school",
							},
						},
					},
				},
				{
					title:    "Empty",
					username: "username",
					response: api.GetProfiles200JSONResponse{
						Data: []api.UserProfile{},
					},
				},
			}

			for _, tt := range tc {
				s.Run(tt.title, func() {
					usersRepo := mocks.NewMockUsersRepo(s.T())
					usersRepo.On("GetProfile", ctx, tt.username).Return(tt.profile, nil)

					logger, hook := test.NewNullLogger()
					logger.SetLevel(logrus.DebugLevel)

					profiles := handlers.NewProfilesHandler(usersRepo, logger)

					result, err := profiles.GetProfiles(ctx, api.GetProfilesRequestObject{
						Params: api.GetProfilesParams{
							Username: &tt.username,
						},
					})
					s.Require().NoError(err)

					s.Require().IsType(api.GetProfiles200JSONResponse{}, result)

					response, _ := result.(api.GetProfiles200JSONResponse)

					s.Equal(tt.response, response)
					s.Nil(hook.LastEntry())
				})
			}
		})

		s.Run("Fail", func() {
			dbError := errors.New("db error")
			username := "username"

			usersRepo := mocks.NewMockUsersRepo(s.T())
			usersRepo.On("GetProfile", ctx, username).Return(nil, dbError)

			logger, hook := test.NewNullLogger()
			logger.SetLevel(logrus.DebugLevel)

			profiles := handlers.NewProfilesHandler(usersRepo, logger)

			result, err := profiles.GetProfiles(ctx, api.GetProfilesRequestObject{
				Params: api.GetProfilesParams{
					Username: &username,
				},
			})
			s.Require().NoError(err)

			s.Require().IsType(api.GetProfiles500JSONResponse{}, result)

			response, _ := result.(api.GetProfiles500JSONResponse)

			s.Equal("Failed to get profile", response.Message)
			s.Equal(dbError.Error(), response.Reason)

			s.Require().NotNil(hook.LastEntry())
			s.Require().NotEmpty(hook.LastEntry().Data)
			s.Equal("profiles", hook.LastEntry().Data["handler"])
			s.Equal(dbError, hook.LastEntry().Data["error"])
			s.Equal("failed to get profile", hook.LastEntry().Message)
		})
	})
}
