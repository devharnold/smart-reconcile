package reconciler

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

const (
	StatusMatched      = "MATCHED"
	StatusManualReview = "MANUAL_REVIEW"
	StatusFailed       = "FAILED"
)

type Transaction struct {
	ID         uuid.UUID
	BusinessID uuid.UUID
	Provider   string
	ExternalID string
	Reference  string
	Amount   decimal.Decimal
	Currency string
	OccurredAt time.Time
	Status     string
}

type Result struct {
	ID uuid.UUID
	InternalTxID uuid.UUID
	ExternalTxID uuid.UUID
	Status   string
	Variance decimal.Decimal
	Reason   string
	ReconciledAt time.Time
}

type Repository interface {
	CreateResult(ctx context.Context, result *Result) error

	UpdateTransactionStatus(ctx context.Context, transactionID uuid.UUID, status string) error
}

type Engine struct {
	repo      Repository
	tolerance decimal.Decimal
}

func NewEngine(repo Repository, tolerance decimal.Decimal) *Engine {

	return &Engine{
		repo:      repo,
		tolerance: tolerance,
	}
}

func (e *Engine) Reconcile(ctx context.Context, internal Transaction, external Transaction) (*Result, error) {

	if internal.Currency != external.Currency {
		return nil, errors.New("currency mismatch")
	}

	if internal.Reference == "" || external.Reference == "" {
		return nil, errors.New("missing transaction reference")
	}

	variance := internal.Amount.Sub(external.Amount).Abs()

	result := &Result{
		ID:            uuid.New(),
		InternalTxID:  internal.ID,
		ExternalTxID:  external.ID,
		Variance:      variance,
		ReconciledAt:  time.Now().UTC(),
	}
	referencesMatch := MatchByReference(
		internal.Reference,
		external.Reference,
	)

	withinTolerance := variance.LessThanOrEqual(e.tolerance,)

	if referencesMatch && withinTolerance {
		result.Status = StatusMatched

		err := e.repo.UpdateTransactionStatus(
			ctx,
			internal.ID,
			StatusMatched,
		)
		if err != nil {
			return nil, err
		}

	} else {
		result.Status = StatusManualReview
		if !referencesMatch {
			result.Reason = "transaction references do not match"
		} else {
			result.Reason = "variance exceeded tolerance"
		}

		err := e.repo.UpdateTransactionStatus(
			ctx,
			internal.ID,
			StatusManualReview,
		)
		if err != nil {
			return nil, err
		}
	}

	err := e.repo.CreateResult(ctx, result)
	if err != nil {
		return nil, err
	}

	return result, nil
}