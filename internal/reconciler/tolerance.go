package reconciler

import "math"

func WithinTolerance(variance float64, tolerance float64) bool {
	return math.Abs(variance) <= tolerance
}