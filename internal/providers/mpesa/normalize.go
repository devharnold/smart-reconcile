package mpesa

import (
	"time"

	"github.com/devharnold/smart-reconcile/internal/transactions"
)

type Mapper struct{}

func (m Mapper) Normalize(p TransactionPayload) (*transactions.NormalizedTransaction, error) {
	occurredAt, err := time.Parse(time.RFC3339, p.Timestamp)
	if err != nil {
		return nil, err
	}

	return &transactions.NormalizedTransaction{
		ID:         p.TransactionID,
		Provider:   "mpesa",
		Currency:   p.CurrencyCode,
		Reference:  p.Reference,
		OccurredAt: occurredAt,
	}, nil
}
