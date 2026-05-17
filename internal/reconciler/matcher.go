package reconciler

import (
	"regexp"
	"strings"
)

var specialChars = regexp.MustCompile(`[^A-Z0-9]`)

func MatchByReference(internal string, external string,) bool {
	internal = normalizeReference(internal)
	external = normalizeReference(external)

	return internal == external
}

func normalizeReference(value string) string {
	value = strings.TrimSpace(value)
	value = strings.ToUpper(value)

	value = specialChars.ReplaceAllString(
		value,
		"",
	)
	return value
}