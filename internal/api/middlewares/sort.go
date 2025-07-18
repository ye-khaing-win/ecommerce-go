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

func Sort(whitelist map[string]struct{}) func(handler http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			q := r.URL.Query().Get("sort_by")
			entries := strings.Split(q, ",")
			//sorted := make(map[string]SortDir)
			var sorted []string

			if len(entries) > 0 {
				for _, entry := range entries {
					parts := strings.Split(entry, ":")
					if len(parts) != 2 {
						continue
					}
					name := parts[0]
					order := SortDir(strings.ToUpper(parts[1]))
					fmt.Println("Order: ", order)

					if !order.isValid() || !isWhiteListed(name, whitelist) {
						continue
					}

					sorted = append(sorted, fmt.Sprintf("%s %s", name, order))
				}

			}

			ctx := context.WithValue(r.Context(), sortKey{}, sorted)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func isWhiteListed(key string, whitelist map[string]struct{}) bool {
	_, ok := whitelist[key]

	return ok
}

func GetSorts(ctx context.Context) []string {
	v, _ := ctx.Value(sortKey{}).([]string)
	return v
}
