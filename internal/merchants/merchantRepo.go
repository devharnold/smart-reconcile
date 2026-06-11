package merchants

import (
	"context"
	"time"

	"github.com/devharnold/smart-reconcile/internal/storage"
)


type MerchantsRepository struct {
	db *storage.DB
}

// constructor
func NewMerchantsRepository(db *storage.DB) *MerchantsRepository {
	return &MerchantsRepository{db: db}
}

func (r *MerchantsRepository) RegisterMerchant(ctx context.Context, merchants *Merchants) (int64, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	query := `
		INSERT INTO merchants (
			first_name,
			last_name,
			email,
			phone_number,
			password,
			created_at,
			updated_at
		)
		VALUES ($1, $2, $3, $4, $5, NOW(), NOW())
		RETURNING merchant_id
	`

	var merchantID int64
	var u Merchants
	err := r.db.Pool.QueryRow(
		ctx,
		query,
		u.FirstName,
		u.LastName,
		u.Email,
		u.PhoneNumber,
		u.Password,
	).Scan(&merchantID)

	if err != nil {
		return 0, err
	}

	return merchantID, nil
}

func (r *MerchantsRepository) GetMerchantByEmail(ctx context.Context, email string) (*Merchants, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	query := `
		SELECT
			merchant_id,
			first_name,
			last_name,
			email,
			phone_number,
			password
		FROM merchants
		WHERE email = $1
	`

	var u Merchants
	err := r.db.Pool.QueryRow(ctx, query, email).Scan(
		&u.MerchantID,
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

func (r *MerchantsRepository) GetByID(ctx context.Context, merchantID int64) (*Merchants, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	query := `
		SELECT
			merchant_id,
			first_name,
			last_name,
			email,
			phone_number
		FROM merchants
		WHERE merchant_id = $1
	`

	var u Merchants
	err := r.db.Pool.QueryRow(ctx, query, merchantID).Scan(
		&u.MerchantID,
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
