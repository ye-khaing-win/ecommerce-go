package middlewares

import (
	"context"
	"fmt"
	"net/http"
	"strings"
)

const (
	asc  SortDir = "ASC"
	desc SortDir = "DESC"
)

type SortDir string

func (o SortDir) isValid() bool {
	return o == asc || o == desc
}

type sortKey struct{}

func Sort(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query().Get("sort_by")
		entries := strings.Split(q, ",")
		sorted := make(map[string]SortDir)

		if len(entries) > 0 {
			for _, entry := range entries {
				parts := strings.Split(entry, ":")
				if len(parts) != 2 {
					continue
				}
				name := parts[0]
				order := SortDir(strings.ToUpper(parts[1]))
				fmt.Println("Order: ", order)

				if !order.isValid() || !isWhiteListed(name) {
					continue
				}

				sorted[name] = order
			}

		}

		ctx := context.WithValue(r.Context(), sortKey{}, sorted)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func isWhiteListed(key string) bool {
	_, ok := whitelist[key]

	return ok
}

func Sorted(ctx context.Context) map[string]SortDir {
	v, _ := ctx.Value(sortKey{}).(map[string]SortDir)
	return v
}
