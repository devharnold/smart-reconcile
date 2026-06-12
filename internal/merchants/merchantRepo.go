package merchants

import (
	"context"
	"time"

	"github.com/devharnold/smart-reconcile/internal/storage"
)

type MerchantRepository struct {
	db *storage.DB
}

func NewMerchantsRepository(db *storage.DB) *MerchantRepository {
	return &MerchantRepository{db: db}
}

func (r *MerchantRepository) RegisterMerchant(ctx context.Context, merchants *Merchants) (int64, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	query := `
		INSERT INTO merchants (
			business_name,
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
	var u Merchants
	err := r.db.Pool.QueryRow(
		ctx,
		query,
		u.BusinessName,
		u.Email,
		u.PhoneNumber,
		u.Password,
	).Scan(&userID)

	if err != nil {
		return 0, err
	}

	return userID, nil
}

func (r *MerchantRepository) GetMerchantByEmail(ctx context.Context, email string) (*Merchants, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	query := `
		SELECT
			user_id,
			business_name,
			email,
			phone_number,
			password
		FROM merchants
		WHERE email = $1
	`

	var u Merchants
	err := r.db.Pool.QueryRow(ctx, query, email).Scan(
		&u.UserID,
		&u.BusinessName,
		&u.Email,
		&u.PhoneNumber,
		&u.Password,
	)

	if err != nil {
		return nil, err
	}

	return &u, nil
}

func (r *MerchantRepository) GetByID(ctx context.Context, merchantID int64) (*Merchants, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	query := `
		SELECT
			business_name,
			email,
			phone_number
		FROM merchants
		WHERE user_id = $1
	`

	var u Merchants
	err := r.db.Pool.QueryRow(ctx, query, merchantID).Scan(
		&u.UserID,
		&u.BusinessName,
		&u.Email,
		&u.PhoneNumber,
	)

	if err != nil {
		return nil, err
	}

	return &u, nil
}
