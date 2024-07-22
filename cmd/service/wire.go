//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.

package main

import (
	// "context"

	"github.com/gofiber/fiber/v2"
	"github.com/google/wire"
	"github.com/sirupsen/logrus"

	"github.com/andrii-lavreniuk/oapi-fiber-example/internal/auth"
	"github.com/andrii-lavreniuk/oapi-fiber-example/internal/config"
	"github.com/andrii-lavreniuk/oapi-fiber-example/internal/data"
	"github.com/andrii-lavreniuk/oapi-fiber-example/internal/handlers"
	"github.com/andrii-lavreniuk/oapi-fiber-example/internal/openapi"
	"github.com/andrii-lavreniuk/oapi-fiber-example/internal/server"
)

// wireApp init application.
func wireApp(config.Server, config.Data, *logrus.Logger) (*fiber.App, func(), error) {
	panic(wire.Build(
		handlers.ProviderSet,
		auth.ProviderSet,
		server.ProviderSet,
		openapi.ProviderSet,
		data.ProviderSet,
	))
}
