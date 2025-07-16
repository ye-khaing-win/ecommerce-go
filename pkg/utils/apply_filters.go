package utils

import (
	"fmt"
	"strings"
)

func ApplyFilters(filters map[string]string) (string, []any) {
	var args []any
	var clauses []string

	for k, v := range filters {
		clauses = append(clauses, fmt.Sprintf("%s = ?", k))
		args = append(args, v)
	}

	if len(clauses) == 0 {
		return "", nil
	}

	return " AND " + strings.Join(clauses, " AND "), args
}
