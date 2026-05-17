package reconciler

import "strings"

func MatchByReference(internal string, external string,) bool {
	internal = normalizeReference(internal)
	external = normalizeReference(external)

	return internal == external
}

func normalizeReference(value string) string {
	value = strings.TrimSpace(value)
	value = strings.ToUpper(value)

	return value
}