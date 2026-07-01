package mpesa

import (
	"encoding/json"
	"time"

	"github.com/devharnold/smart-reconcile/internal/transactions"
)

type Mapper struct{}

func (m Mapper) Normalize(payload []byte) (*transactions.NormalizedTransaction, error) {

	var p TransactionPayload

	err := json.Unmarshal(payload, &p)
	if err != nil {
		return nil, err
	}

	occurredAt, err := time.Parse(time.RFC3339, p.Timestamp)
	if err != nil {
		return nil, err
	}

	return &transactions.NormalizedTransaction{
		ID:         p.TransactionID,
		Provider:   "mpesa",
		Amount:     p.Amount,
		Currency:   p.CurrencyCode,
		Reference:  p.Reference,
		OccurredAt: occurredAt,
	}, nil
}
