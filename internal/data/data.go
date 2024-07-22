package data

import (
	"github.com/andrii-lavreniuk/oapi-fiber-example/internal/config"

	"database/sql"

	"github.com/google/wire"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/mysqldialect"
	"github.com/uptrace/bun/extra/bundebug"
	"github.com/uptrace/bun/schema"
)

var ProviderSet = wire.NewSet(NewData, NewAuthRepo, NewUsersRepo)

type Data struct {
	db *bun.DB
}

func NewData(c config.Data) (*Data, error) {
	// return table names as is (without pluralization)
	schema.SetTableNameInflector(func(s string) string {
		return s
	})

	sqldb, err := sql.Open("mysql", c.MySQL.URL)
	if err != nil {
		return nil, err
	}

	if err = sqldb.Ping(); err != nil {
		return nil, err
	}

	db := bun.NewDB(sqldb, mysqldialect.New())
	db.AddQueryHook(bundebug.NewQueryHook(bundebug.WithVerbose(c.MySQL.Verbose)))

	return &Data{
		db: db,
	}, nil
}
