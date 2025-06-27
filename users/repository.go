package users

import (
	"context"

	"github.com/uptrace/bun"
)

// Repository provides access to user storage.
type UserRepository struct {
	db *bun.DB
}

// NewRepository creates a new user repository.
func NewRepository(db *bun.DB) *UserRepository {
	return &UserRepository{db}
}

// Create makes a new user from User
func (r *UserRepository) Create(user *User) error {
	ctx := context.Background()
	_, err := r.db.NewInsert().Model(user).Exec(ctx)
	return err
}

// Update updates a user with user.ID with User
func (r *UserRepository) Update(user *User) error {
	ctx := context.Background()
	_, err := r.db.NewUpdate().Model(user).Where("id = ?", user.ID).Exec(ctx)
	return err
}

// Delete deletes a user with ID
func (r *UserRepository) Delete(id int64) error {
	ctx := context.Background()
	_, err := r.db.NewDelete().Model((*User)(nil)).Where("id = ?", id).Exec(ctx)
	return err
}

// GetByUsername retrieves a user by username.
func (r *UserRepository) FindByUsername(username string) (*User, error) {
	ctx := context.Background()
	user := new(User)
	err := r.db.NewSelect().Model(user).Where("username = ?", username).Scan(ctx)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *UserRepository) FindByID(id int64) (*User, error) {
	ctx := context.Background()
	user := new(User)
	err := r.db.NewSelect().Model(user).Where("id = ?", id).Scan(ctx)
	if err != nil {
		return nil, err
	}

	return user, nil
}
