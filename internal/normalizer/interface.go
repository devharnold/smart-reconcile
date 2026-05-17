package normalizer

import (
	"time"

	"github.com/shopspring/decimal"
)

type NormalizedTransaction struct {
	ID string
	Provider string
	Amount decimal.Decimal
	Currency string
	Reference string
	OccurredAt time.Time
}

type TransactionNormalizer interface {
	Normalize(payload []byte) (*NormalizedTransaction, error)
}