package database

import (
	"context"
	"database/sql"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/sqlitedialect"
	"github.com/uptrace/bun/driver/sqliteshim"
	"github.com/uptrace/bun/extra/bundebug"
)

type Database struct {
	DB *bun.DB
}

func NewDatabase(dbPath string) (*Database, error) {
	sqldb, err := sql.Open(sqliteshim.ShimName, dbPath)
	if err != nil {
		return nil, err
	}

	bundb := bun.NewDB(sqldb, sqlitedialect.New())
	return &Database{DB: bundb}, nil
}

// EnableVerbose enables verbose query logging.
func (d *Database) EnableVerbose() {
	d.DB.AddQueryHook(bundebug.NewQueryHook(
		bundebug.WithVerbose(true),
		bundebug.FromEnv("BUNDEBUG"),
	))
}

// CreateTables creates all tables
func (d *Database) CreateTables(ctx context.Context, models ...interface{}) error {
	for _, model := range models {
		query := d.DB.NewCreateTable().
			Model(model).
			IfNotExists()

		_, err := query.Exec(ctx)
		if err != nil {
			return err
		}
	}

	return nil
}

func (d *Database) Close() error {
	if d.DB == nil {
		return nil
	}
	if err := d.DB.Close(); err != nil {
		return err
	}
	return nil
}
