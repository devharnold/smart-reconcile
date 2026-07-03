package transactions

import (
	"time"

	"github.com/shopspring/decimal"
)

type NormalizedTransaction struct {
	ID         string
	Provider   string
	Amount     decimal.Decimal
	Currency   string
	Reference  string
	OccurredAt time.Time
}

type Transaction struct {
	ID          string
	Amount      float64
	Currency    string
	Reference   string
	Description string
	Timestamp   string
}
