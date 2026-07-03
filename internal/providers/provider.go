package providers

import "github.com/devharnold/smart-reconcile/internal/transactions"

type Provider interface {
	Name() string
	FetchTransactions(accoundID string) ([]transactions.Transaction, error)
}
