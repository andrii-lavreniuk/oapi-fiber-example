package server

import (
	"net/http"
	"strings"

	"github.com/gofiber/adaptor/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/keyauth"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/google/wire"
	"github.com/sirupsen/logrus"

	api_profiles "github.com/andrii-lavreniuk/oapi-fiber-example/gen/api/profiles"
	"github.com/andrii-lavreniuk/oapi-fiber-example/internal/config"
	"github.com/andrii-lavreniuk/oapi-fiber-example/internal/handlers"
	"github.com/andrii-lavreniuk/oapi-fiber-example/internal/interfaces"
	"github.com/andrii-lavreniuk/oapi-fiber-example/internal/openapi"
)

var ProviderSet = wire.NewSet(New)

const openapiSpecPath = "/docs"

func New(
	cfg config.Server,
	auth interfaces.Auth,
	profilesHandler *handlers.ProfilesHandler,
	spec *openapi.Spec,
	lg *logrus.Logger,
) (*fiber.App, error) {
	app := fiber.New(fiber.Config{
		IdleTimeout:           cfg.IdleTimeout,
		DisableStartupMessage: true,
	})

	app.Hooks().OnListen(func(data fiber.ListenData) error {
		lg.WithField("port", data.Port).Info("server started")

		if cfg.OpenAPI {
			lg.Infof("OpenAPI available at http://%s:%s%s", data.Host, data.Port, openapiSpecPath)
		}

		return nil
	})

	app.Use(recover.New(recover.Config{
		EnableStackTrace: cfg.Debug,
	}))

	app.Use(requestid.New(
		requestid.Config{
			Header: "X-Request-ID",
		},
	))

	app.Use(keyauth.New(keyauth.Config{
		KeyLookup:  "header:X-API-Key",
		AuthScheme: "",
		Validator: func(c *fiber.Ctx, key string) (bool, error) {
			return auth.ValidateAPIKey(c.Context(), key)
		},
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			if err != nil {
				lg.WithError(err).Warn("authentication failed")
			}

			return c.Status(fiber.StatusForbidden).JSON(fiber.ErrForbidden)
		},
		Next: func(c *fiber.Ctx) bool {
			return c.Method() == http.MethodGet && strings.HasPrefix(c.Path(), openapiSpecPath)
		},
	}))

	// log requests
	app.Use(func(c *fiber.Ctx) error {
		// continue processing the request to get the status code
		_ = c.Next()

		status := c.Response().StatusCode()

		fields := logrus.Fields{
			"request.id": c.Locals("requestid"),
			"method":     c.Method(),
			"path":       c.Path(),
			"status":     status,
		}

		var lvl logrus.Level

		switch {
		case status >= http.StatusInternalServerError:
			lvl = logrus.ErrorLevel
		case status >= http.StatusBadRequest:
			lvl = logrus.WarnLevel
		default:
			lvl = logrus.InfoLevel
		}

		lg.WithFields(fields).Log(lvl, "request")

		return nil
	})

	// register profiles handlers
	api_profiles.RegisterHandlers(app, api_profiles.NewStrictHandler(profilesHandler, nil))

	// register OpenAPI spec
	if cfg.OpenAPI {
		app.Get(openapiSpecPath+"*", adaptor.HTTPHandler(spec.Redoc(openapiSpecPath).Handler()))
	}

	// 404 handler
	app.All("/*", func(c *fiber.Ctx) error {
		_ = c.SendStatus(fiber.StatusNotFound)
		return c.JSON(fiber.ErrNotFound)
	})

	return app, nil
}
