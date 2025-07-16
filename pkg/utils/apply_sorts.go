package utils

import "strings"

func ApplySorts(sorts []string) string {
	if len(sorts) == 0 {
		return ""
	}

	return "ORDER BY " + strings.Join(sorts, ", ")
}
