package preferences

import (
	"context"
	"fmt"

	"github.com/uptrace/bun"
)

// Repository provides access to preference storage.
type PrefRepository struct {
	db *bun.DB
}

// NewRepository creates a new preference repository.
func NewRepository(db *bun.DB) (*PrefRepository, error) {
	err := createPrefsTable(db)
	if err != nil {
		return nil, fmt.Errorf("error creating preferences table: %w", err)
	}

	return &PrefRepository{db}, nil
}

// createPrefsTable creates the preferences table in the database.
func createPrefsTable(db *bun.DB) error {
	ctx := context.Background()
	_, err := db.NewCreateTable().
		Model((*Preference)(nil)).
		IfNotExists().
		ColumnExpr("units TEXT NOT NULL CHECK (units IN ('us', 'si')) DEFAULT 'us'").
		ForeignKey(`("user_id") REFERENCES "users" ("id") ON DELETE CASCADE`).
		Exec(ctx)
	if err != nil {
		return err
	}

	return nil
}

// Create adds a user's preference
func (r *PrefRepository) Create(pref *Preference) error {
	ctx := context.Background()
	_, err := r.db.NewInsert().Model(pref).Exec(ctx)
	return err
}

// Update updates a user's preference
func (r *PrefRepository) Update(pref *Preference) error {
	ctx := context.Background()
	_, err := r.db.NewUpdate().Model(pref).Where("user_id = ?", pref.UserID).Exec(ctx)
	return err
}

func (r *PrefRepository) Delete(userID int64) error {
	ctx := context.Background()
	_, err := r.db.NewDelete().Model((*Preference)(nil)).Where("user_id = ?", userID).Exec(ctx)
	return err
}

func (r *PrefRepository) FindByUserID(userID int64) (*Preference, error) {
	ctx := context.Background()
	pref := new(Preference)
	err := r.db.NewSelect().Model(pref).Where("user_id = ?", userID).Scan(ctx)
	if err != nil {
		return nil, err
	}
	return pref, nil
}
