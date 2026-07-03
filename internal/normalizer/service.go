package normalizer

import (
	"fmt"

	"github.com/devharnold/smart-reconcile/internal/transactions"
)

type Service struct{}

func (s Service) Normalize(provider string, payload []byte) (*transactions.NormalizedTransaction, error) {
	normalizer, exists := Providers[provider]
	if !exists {
		return nil, fmt.Errorf("Unsupported provider")
	}
	return normalizer.Normalize(payload)
}
