package normalizer

import (
	"fmt"
)

type Service struct {}

func(s Service) Normalize(provider string, payload []byte) (*NormalizedTransaction, error) {
	normalizer, exists := Providers[provider]
	if !exists {
		return nil, fmt.Errorf("Unsupported provider")
	}
	return normalizer.Normalize(payload)
}