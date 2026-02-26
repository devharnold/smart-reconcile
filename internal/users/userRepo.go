package users

import (
	"context"
	"time"

	"github.com/devharnold/smart-reconcile/internal/storage"
)


type UsersRepository struct {
	db *storage.DB
}

// constructor
func NewUsersRepository(db *storage.DB) *UsersRepository {
	return &UsersRepository{db: db}
}

func (r *UsersRepository) RegisterUser(ctx context.Context, users *Users) (int64, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	query := `
		INSERT INTO users (
			first_name,
			last_name,
			email,
			phone_number,
			password,
			created_at,
			updated_at
		)
		VALUES ($1, $2, $3, $4, $5, NOW(), NOW())
		RETURNING user_id
	`

	var userID int64
	var u Users
	err := r.db.Pool.QueryRow(
		ctx,
		query,
		u.FirstName,
		u.LastName,
		u.Email,
		u.PhoneNumber,
		u.Password,
	).Scan(&userID)

	if err != nil {
		return 0, err
	}

	return userID, nil
}

func (r *UsersRepository) GetUserByEmail(ctx context.Context, email string) (*Users, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	query := `
		SELECT
			user_id,
			first_name,
			last_name,
			email,
			phone_number,
			password
		FROM users
		WHERE email = $1
	`

	var u Users
	err := r.db.Pool.QueryRow(ctx, query, email).Scan(
		&u.UserID,
		&u.FirstName,
		&u.LastName,
		&u.Email,
		&u.PhoneNumber,
		&u.Password,
	)

	if err != nil {
		return nil, err
	}

	return &u, nil
}

func (r *UsersRepository) GetByID(ctx context.Context, userID int64) (*Users, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	query := `
		SELECT
			user_id,
			first_name,
			last_name,
			email,
			phone_number
		FROM users
		WHERE user_id = $1
	`

	var u Users
	err := r.db.Pool.QueryRow(ctx, query, userID).Scan(
		&u.UserID,
		&u.FirstName,
		&u.LastName,
		&u.Email,
		&u.PhoneNumber,
	)

	if err != nil {
		return nil, err
	}

	return &u, nil
}
