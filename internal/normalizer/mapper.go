// Normalization Registry
package normalizer

import (
	"github.com/devharnold/smart-reconcile/internal/providers/mpesa"
)

var Providers = map[string]TransactionNormalizer{
	"mpesa": mpesa.Mapper{},
}