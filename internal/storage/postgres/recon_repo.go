package postgres

import (
	"context"
	"database/sql"

	"github.com/devharnold/smart-reconcile/internal/reconciler"
	"github.com/google/uuid"
)

type ReconciliationRepo struct {
	db *sql.DB
}

func NewReconciliationRepository(db *sql.DB,) *ReconciliationRepo {
	return &ReconciliationRepo{
		db: db,
	}
}

func (r *ReconciliationRepo) CreateResult(ctx context.Context, result *reconciler.Result,) error {
	query := `
	INSERT INTO reconciliation_results (
		id,
		internal_tx_id,
		external_tx_id,
		status,
		variance,
		reason,
		reconciled_at
	)
	VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	_, err := r.db.ExecContext(
		ctx,
		query,
		result.ID,
		result.InternalTxID,
		result.ExternalTxID,
		result.Status,
		result.Variance,
		result.Reason,
		result.ReconciledAt,
	)
	return err
}

func (r *ReconciliationRepo) UpdateTransactionStatus (ctx context.Context, transactionID uuid.UUID, status string,) error {
	query := `
	UPDATE transactions
	SET status = $1
	WHERE id = $2
	`

	_, err := r.db.ExecContext(
		ctx,
		query,
		status,
		transactionID,
	)

	return err
}