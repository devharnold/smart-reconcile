package reconciler

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/devharnold/smart-reconcile/internal/storage/postgres"
)

const (
	StatusMatched = "MATCHED"
	StatusManualReview = "MANUAL_REVIEW"
	StatusFailed = "FAILED"
)

type Transaction struct {
	ID uuid.UUID
	BusinessID uuid.UUID
	Provider string
	ExternalID string
	Reference string
	Amount float64
	Currency string
	OccurredAt time.Time
	Status string
}

type Result struct {
	ID            uuid.UUID
	InternalTxID  uuid.UUID
	ExternalTxID  uuid.UUID
	Status        string
	Variance      float64
	Reason        string
	ReconciledAt  time.Time
}

type Repository interface {
	CreateResult(ctx context.Context, result *Result) error
}

type Engine struct {
	repo Repository
	tolerance float64
}

func NewEngine(repo Repository, tolerance float64) *Engine {
	return &Engine{
		repo: repo,
		tolerance: tolerance,
	}
}

func (e *Engine) Reconcile(ctx context.Context, internal Transaction, external Transaction) (*Result, error) {
	if internal.Currency != external.Currency {
		return nil, errors.New("Currency Mismatch")
	}

	if internal.Reference == "" || external.Reference == "" {
		return nil, errors.New("Missing References")
	}

	variance := internal.Amount - external.Amount
	result := &Result {
		ID: uuid.New(),
		InternalTxID: internal.ID,
		ExternalTxID: external.ID,
		Variance: variance,
		ReconciledAt: time.Now().UTC(),
	}

	if MatchByReference(internal.Reference, external.Reference) && WithinTolerance(variance, e.tolerance) {
		result.Status = StatusMatched

		if err := e.repo.UpdateTransactionStatus(
			ctx,
			internal.ID,
			StatusMatched,
		); err != nil {
			return nil, err
		}
	} else {
		result.Status = StatusManualReview
		result.Reason = "Variance exceeded tolerance"

		if err := e.repo.UpdateTransactionStatus(
			ctx,
			internal.ID,
			StatusManualReview,
		); err != nil {
			return nil, err
		}
	}

	if err := e.repo.CreateResult(ctx, result); err != nil {
		return nil, err
	}

	return result, nil
}