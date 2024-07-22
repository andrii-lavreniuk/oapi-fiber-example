package migrations

import (
	"embed"

	"github.com/uptrace/bun/migrate"
)

//nolint:gochecknoglobals // required by bun/migrate
var Migrations = migrate.NewMigrations()

//go:embed *.sql
var sqlMigrations embed.FS

func init() { //nolint:gochecknoinits // required by bun/migrate
	if err := Migrations.Discover(sqlMigrations); err != nil {
		panic(err)
	}
}
