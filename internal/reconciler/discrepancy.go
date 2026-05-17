package reconciler

import (
	"time"

	"github.com/devharnold/smart-reconcile/internal/normalizer"
)

type DiscrepancyType string

const (
	MissingTransaction DiscrepancyType = "MISSING_TRANSACTION"
	AmountMismatch DiscrepancyType = "AMOUNT_MISMATCH"
	CurrencyMismatch DiscrepancyType = "CURRENCY_MISMATCH"
	DuplicateTransaction DiscrepancyType = "DUPLICATE_TRANSACTION"
)

type Discrepancy struct {
	Type DiscrepancyType
	Description string

	Excpected *normalizer.NormalizedTransaction
	Actual *normalizer.NormalizedTransaction

	DetectedAt time.Time
}