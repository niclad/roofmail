package testutil

import (
	"context"
	"roofmail/database"
	"testing"
)

func SetupTestDB(t *testing.T) *database.Database {
	dsn := "file::memory:?cache=shared" // Use an in-memory SQLite database for testing
	db, err := database.NewDatabase(dsn)
	if err != nil {
		t.Fatalf("failed to open test database: %v", err)
	}

	return db
}

// Drop thhe tables for the given models in the test database.
func DropTables(t *testing.T, db *database.Database, models ...interface{}) {
	ctx := context.Background()
	for _, model := range models {
		_, err := db.DB.NewDropTable().Model(model).IfExists().Exec(ctx)
		if err != nil {
			t.Fatalf("failed to drop table for %T: %v", model, err)
		}
	}
	t.Log("Tables dropped successfully")
}
