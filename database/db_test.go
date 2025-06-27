package database

import (
	"context"
	"testing"
)

var testDB *Database

func setupTestDB(t *testing.T) *Database {
	dsn := "file::memory:?cache=shared"
	db, err := NewDatabase(dsn)
	if err != nil {
		t.Fatalf("failed to open test database: %v", err)
	}
	return db
}

type User struct {
	ID   int64  `bun:",pk,autoincrement"`
	Name string `bun:",notnull"`
}

func TestNewDatabase(t *testing.T) {
	db := setupTestDB(t)
	if db.DB == nil {
		t.Error("expected bun.DB to be initialized")
	}
}

func TestEnableVerbose(t *testing.T) {
	db := setupTestDB(t)
	// Should not panic or error
	db.EnableVerbose()
}

func TestCreateTables(t *testing.T) {
	db := setupTestDB(t)
	ctx := context.Background()
	err := db.CreateTables(ctx, (*User)(nil))
	if err != nil {
		t.Errorf("CreateTables failed: %v", err)
	}

	// Try creating again to test IfNotExists
	err = db.CreateTables(ctx, (*User)(nil))
	if err != nil {
		t.Errorf("CreateTables failed on second call: %v", err)
	}
}

func TestInsertAndQueryUser(t *testing.T) {
	db := setupTestDB(t)
	ctx := context.Background()
	err := db.CreateTables(ctx, (*User)(nil))
	if err != nil {
		t.Fatalf("CreateTables failed: %v", err)
	}

	user := &User{Name: "Alice"}
	_, err = db.DB.NewInsert().Model(user).Exec(ctx)
	if err != nil {
		t.Fatalf("failed to insert user: %v", err)
	}

	var users []User
	err = db.DB.NewSelect().Model(&users).Where("name = ?", "Alice").Scan(ctx)
	if err != nil {
		t.Fatalf("failed to query user: %v", err)
	}
	if len(users) != 1 || users[0].Name != "Alice" {
		t.Errorf("unexpected user query result: %+v", users)
	}
}
