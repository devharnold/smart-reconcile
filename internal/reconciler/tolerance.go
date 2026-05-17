package reconciler

import (
	"github.com/shopspring/decimal"
)

func WithinTolerance(variance decimal.Decimal, tolerance decimal.Decimal) bool {
	return variance.LessThanOrEqual(tolerance)
}