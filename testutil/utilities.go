package testutil

import (
	"roofmail/database"
	"testing"
)

func SetupTestDB(t *testing.T) *database.Database {
	dsn := "file::memory?cache=shared"
	db, err := database.NewDatabase(dsn)
	if err != nil {
		t.Fatalf("failed to open test database: %v", err)
	}

	return db
}
