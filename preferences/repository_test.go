package preferences

import (
	"context"
	"fmt"
	"roofmail/database"
	"roofmail/testutil"
	"roofmail/users"
	"testing"
)

// testDB creates a test database for the preferences package.
func testDB(t *testing.T) *database.Database {
	db := testutil.SetupTestDB(t)
	dropTables(t, db) // just in case any previous tests left tables behind
	createPrefsTable(db.DB)
	db.CreateTables(context.Background(), (*users.User)(nil))
	return db
}

func dropTables(t *testing.T, db *database.Database) {
	testutil.DropTables(t, db, (*Preference)(nil), (*users.User)(nil))
	t.Log("Test tables dropped successfully")
}

func closeDB(t *testing.T, db *database.Database) {
	dropTables(t, db)
	// Close the database connection
	// This is a no-op for the in-memory database, but included for completeness.
	if err := db.Close(); err != nil {
		t.Fatalf("failed to close test database: %v", err)
	}
	t.Log("Test database closed successfully")
}

func createTestPreference(repo *PrefRepository) error {
	pref := &Preference{
		UserID:           1, // Assuming a test user with ID 1 exists
		Units:            "us",
		TemperatureMin:   0,
		TemperatureMax:   1,
		WindMin:          2,
		WindMax:          3,
		PrecipitationMin: 4,
		PrecipitationMax: 5,
		Locale:           "en-US",
	}

	if err := repo.Create(pref); err != nil {
		return fmt.Errorf("failed to create test preference: %w", err)
	}
	return nil
}

// TestCreatePrefsTable tests the creation of the preferences table.
func TestCreatPrefsTable(t *testing.T) {
	db := testutil.SetupTestDB(t)
	dropTables(t, db) // Ensure no previous tables exist

	// Create the preferences table
	if err := createPrefsTable(db.DB); err != nil {
		t.Fatalf("failed to create preferences table: %v", err)
	}

	var pref []Preference
	// Check if the preferences table exists
	db.DB.NewSelect().Model(&pref).Scan(context.Background())
	fmt.Println("Preferences table exists:", pref)

	// Verify the table exists
	var count int
	res, err := db.DB.NewSelect().Model((*Preference)(nil)).Count(context.Background())
	if err != nil {
		t.Fatalf("failed to count preferences table: %v", err)
	}
	count = res
	if err != nil {
		t.Fatalf("failed to count preferences table: %v", err)
	}
	if count != 0 {
		t.Fatalf("expected preferences table to be empty, got %d rows", count)
	}
	t.Log("Preferences table created successfully")
	if err := db.Close(); err != nil {
		t.Fatalf("failed to close test database: %v", err)
	}
	t.Log("Test database closed successfully")
}

func TestCreate(t *testing.T) {
	db := testDB(t)
	defer closeDB(t, db)

	repo, err := NewRepository(db.DB)
	if err != nil {
		t.Fatalf("failed to create repository: %v", err)
	}

	pref := &Preference{
		// UserID:           testUser.ID,
		Units:            "us",
		TemperatureMin:   1,
		TemperatureMax:   2,
		WindMin:          3,
		WindMax:          4,
		PrecipitationMin: 5,
		PrecipitationMax: 6,
		Locale:           "es-ES",
	}

	if err := repo.Create(pref); err != nil {
		// print the error
		t.Fatalf("failed to create preference: %v", err)
	}
	// Verify the preference was created
	var count int
	res, err := db.DB.NewSelect().Model((*Preference)(nil)).Count(context.Background())
	if err != nil {
		t.Fatalf("failed to count preferences: %v", err)
	}
	count = res
	if count != 1 {
		t.Fatalf("expected 1 preference, got %d", count)
	}
	t.Log("Preference created successfully")
}

func TestUpdate(t *testing.T) {
	db := testDB(t)
	defer closeDB(t, db)

	repo, err := NewRepository(db.DB)
	if err != nil {
		t.Fatalf("failed to create repository: %v", err)
	}

	// Create a test preference
	if err := createTestPreference(repo); err != nil {
		t.Fatalf("failed to create test preference: %v", err)
	}

	pref := &Preference{
		UserID:           1, // Assuming a test user with ID 1 exists
		Units:            "us",
		TemperatureMin:   10,
		TemperatureMax:   20,
		WindMin:          30,
		WindMax:          40,
		PrecipitationMin: 50,
		PrecipitationMax: 60,
		Locale:           "fr-FR",
	}

	if err := repo.Update(pref); err != nil {
		t.Fatalf("failed to update preference: %v", err)
	}
	// Verify the preference was updated
	var updatedPref Preference
	ctx := context.Background()
	err = db.DB.NewSelect().Model(&updatedPref).Where("user_id = ?", pref.UserID).Scan(ctx)
	if err != nil {
		t.Fatalf("failed to find updated preference: %v", err)
	}
	if updatedPref.TemperatureMin != pref.TemperatureMin ||
		updatedPref.TemperatureMax != pref.TemperatureMax ||
		updatedPref.WindMin != pref.WindMin ||
		updatedPref.WindMax != pref.WindMax ||
		updatedPref.PrecipitationMin != pref.PrecipitationMin ||
		updatedPref.PrecipitationMax != pref.PrecipitationMax ||
		updatedPref.Locale != pref.Locale {
		t.Fatalf("preference not updated correctly: got %+v, want %+v", updatedPref, pref)
	}
}

func TestDelete(t *testing.T) {
	db := testDB(t)
	defer closeDB(t, db)

	repo, err := NewRepository(db.DB)
	if err != nil {
		t.Fatalf("failed to create repository: %v", err)
	}

	// Create a test preference
	if err := createTestPreference(repo); err != nil {
		t.Fatalf("failed to create test preference: %v", err)
	}

	// Delete the preference
	if err := repo.Delete(1); err != nil {
		t.Fatalf("failed to delete preference: %v", err)
	}

	// Verify the preference was deleted
	var count int
	res, err := db.DB.NewSelect().Model((*Preference)(nil)).Count(context.Background())
	if err != nil {
		t.Fatalf("failed to count preferences: %v", err)
	}
	count = res
	if count != 0 {
		t.Fatalf("expected 0 preferences, got %d", count)
	}
	t.Log("Preference deleted successfully")
}

func TestFindByUserID(t *testing.T) {
	db := testDB(t)
	defer closeDB(t, db)

	repo, err := NewRepository(db.DB)
	if err != nil {
		t.Fatalf("failed to create repository: %v", err)
	}

	// Create a test preference
	if err := createTestPreference(repo); err != nil {
		t.Fatalf("failed to create test preference: %v", err)
	}

	// Find the preference by user ID
	pref, err := repo.FindByUserID(1)
	if err != nil {
		t.Fatalf("failed to find preference by user ID: %v", err)
	}
	if pref == nil {
		t.Fatal("expected preference to be found, got nil")
	}
	if pref.UserID != 1 {
		t.Fatalf("expected user ID 1, got %d", pref.UserID)
	}
	t.Log("Preference found successfully")
}
