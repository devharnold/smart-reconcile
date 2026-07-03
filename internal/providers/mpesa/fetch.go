package mpesa

import (
	"context"

	"github.com/devharnold/smart-reconcile/internal/transactions"
)

type Mpesa struct {
	client *Client
}

func New(client *Client) *Mpesa {
	return &Mpesa{client: client}
}

func (m *Mpesa) FetchTransactions(ctx context.Context) ([]transactions.NormalizedTransaction, error) {
	rawTransactions, err := m.client.PullTransactions(ctx)
	if err != nil {
		return nil, err
	}

	mapper := Mapper{}

	var result []transactions.NormalizedTransaction
	for _, raw := range rawTransactions {
		tx, err := mapper.Normalize(raw)
		if err != nil {
			return nil, err
		}

		result = append(result, *tx)
	}

	return result, nil
}
