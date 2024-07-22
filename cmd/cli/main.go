package main

import (
	"database/sql"
	"log"
	"os"
	"strings"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/mysqldialect"
	"github.com/uptrace/bun/extra/bundebug"
	"github.com/uptrace/bun/migrate"
	"github.com/urfave/cli/v2"

	"github.com/andrii-lavreniuk/oapi-fiber-example/migrations"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/joho/godotenv/autoload"
)

func main() {
	app := &cli.App{
		Name: "CLI",
		Commands: []*cli.Command{
			newDBCommand(migrations.Migrations),
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func newDBCommand(migrations *migrate.Migrations) *cli.Command { //nolint:funlen,gocognit // no need to split
	return &cli.Command{
		Name:  "db",
		Usage: "manage database migrations",
		Subcommands: []*cli.Command{
			{
				Name:  "init",
				Usage: "create migration tables",
				Action: func(c *cli.Context) error {
					db, err := getDB()
					if err != nil {
						return err
					}

					migrator := migrate.NewMigrator(db, migrations)
					return migrator.Init(c.Context)
				},
			},
			{
				Name:  "migrate",
				Usage: "migrate database",
				Action: func(c *cli.Context) error {
					db, err := getDB()
					if err != nil {
						return err
					}

					migrator := migrate.NewMigrator(db, migrations, migrate.WithMarkAppliedOnSuccess(true))

					group, err := migrator.Migrate(c.Context)
					if err != nil {
						return err
					}

					if group.ID == 0 {
						log.Printf("there are no new migrations to run\n")
						return nil
					}

					log.Printf("migrated to %s\n", group)
					return nil
				},
			},
			{
				Name:  "rollback",
				Usage: "rollback the last migration group",
				Action: func(c *cli.Context) error {
					db, err := getDB()
					if err != nil {
						return err
					}

					migrator := migrate.NewMigrator(db, migrations)

					group, err := migrator.Rollback(c.Context)
					if err != nil {
						return err
					}

					if group.ID == 0 {
						log.Printf("there are no groups to roll back\n")
						return nil
					}

					log.Printf("rolled back %s\n", group)
					return nil
				},
			},
			{
				Name:  "lock",
				Usage: "lock migrations",
				Action: func(c *cli.Context) error {
					db, err := getDB()
					if err != nil {
						return err
					}

					migrator := migrate.NewMigrator(db, migrations)
					return migrator.Lock(c.Context)
				},
			},
			{
				Name:  "unlock",
				Usage: "unlock migrations",
				Action: func(c *cli.Context) error {
					db, err := getDB()
					if err != nil {
						return err
					}

					migrator := migrate.NewMigrator(db, migrations)
					return migrator.Unlock(c.Context)
				},
			},
			{
				Name:  "create_go",
				Usage: "create Go migration",
				Action: func(c *cli.Context) error {
					db, err := getDB()
					if err != nil {
						return err
					}

					migrator := migrate.NewMigrator(db, migrations)

					name := strings.Join(c.Args().Slice(), "_")
					mf, err := migrator.CreateGoMigration(c.Context, name)
					if err != nil {
						return err
					}
					log.Printf("created migration %s (%s)\n", mf.Name, mf.Path)

					return nil
				},
			},
			{
				Name:  "create_sql",
				Usage: "create up and down SQL migrations",
				Action: func(c *cli.Context) error {
					db, err := getDB()
					if err != nil {
						return err
					}

					migrator := migrate.NewMigrator(db, migrations)

					name := strings.Join(c.Args().Slice(), "_")
					files, err := migrator.CreateSQLMigrations(c.Context, name)
					if err != nil {
						return err
					}

					for _, mf := range files {
						log.Printf("created migration %s (%s)\n", mf.Name, mf.Path)
					}

					return nil
				},
			},
			{
				Name:  "status",
				Usage: "print migrations status",
				Action: func(c *cli.Context) error {
					db, err := getDB()
					if err != nil {
						return err
					}

					migrator := migrate.NewMigrator(db, migrations)

					ms, err := migrator.MigrationsWithStatus(c.Context)
					if err != nil {
						return err
					}

					log.Printf("migrations: %s\n", ms)
					log.Printf("unapplied migrations: %s\n", ms.Unapplied())
					log.Printf("last migration group: %s\n", ms.LastGroup())

					return nil
				},
			},
			{
				Name:  "mark_applied",
				Usage: "mark migrations as applied without actually running them",
				Action: func(c *cli.Context) error {
					db, err := getDB()
					if err != nil {
						return err
					}

					migrator := migrate.NewMigrator(db, migrations)

					group, err := migrator.Migrate(c.Context, migrate.WithNopMigration())
					if err != nil {
						return err
					}

					if group.ID == 0 {
						log.Printf("there are no new migrations to mark as applied\n")
						return nil
					}

					log.Printf("marked as applied %s\n", group)
					return nil
				},
			},
		},
	}
}

func getDB() (*bun.DB, error) {
	sqldb, err := sql.Open("mysql", os.Getenv("SERVICE_DATA_MYSQL_URL"))
	if err != nil {
		return nil, err
	}

	db := bun.NewDB(sqldb, mysqldialect.New())

	db.AddQueryHook(bundebug.NewQueryHook(bundebug.WithVerbose(true)))

	return db, nil
}
