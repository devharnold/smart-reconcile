package normalizer

import "github.com/devharnold/smart-reconcile/internal/transactions"

type TransactionNormalizer interface {
	Normalize(payload []byte) (*transactions.NormalizedTransaction, error)
}
