package middlewares

import (
	"context"
	"net/http"
	"strings"
)

type selectKey struct {
}

func Select(whitelist map[string]struct{}) func(handler http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			raw := r.URL.Query().Get("select") // e.g. ?select=id,name
			fields := make(map[string]struct{})
			if raw != "" {

				for _, f := range strings.Split(raw, ",") {
					f = strings.TrimSpace(f)
					if _, ok := whitelist[f]; ok {
						fields[f] = struct{}{}
					}
				}
			} else {
				for k := range whitelist {
					fields[k] = struct{}{}
				}
			}

			ctx := context.WithValue(r.Context(), selectKey{}, fields)

			next.ServeHTTP(w, r.WithContext(ctx)) // propagate new ctx
		})
	}
}

func Selected(ctx context.Context) map[string]struct{} {
	if m, ok := ctx.Value(selectKey{}).(map[string]struct{}); ok {
		return m
	}
	return nil
}
