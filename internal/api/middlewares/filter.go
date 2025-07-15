package middlewares

import (
	"context"
	"net/http"
)

type filterKey struct {
}

func Filter(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		filters := make(map[string]string)

		for k := range whitelist {
			value := r.URL.Query().Get(k)
			if value != "" {
				filters[k] = value
			}
		}

		ctx := context.WithValue(r.Context(), filterKey{}, filters)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func Filtered(ctx context.Context) map[string]string {
	v, _ := ctx.Value(filterKey{}).(map[string]string)
	return v
}
