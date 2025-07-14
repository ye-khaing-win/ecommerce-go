package middlewares

import (
	"context"
	"net/http"
	"strings"
)

type ctxKey string

const key ctxKey = "selectFields"

var whitelist = map[string]struct{}{
	"id":          {},
	"name":        {},
	"description": {},
}

func Select(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		raw := r.URL.Query().Get("select") // e.g. ?select=id,name
		var fields []string
		if raw != "" {

			for _, f := range strings.Split(raw, ",") {
				if _, ok := whitelist[f]; ok {
					fields = append(fields, f)
				}
			}
		} else {
			for k := range whitelist {
				fields = append(fields, k)
			}
		}

		ctx := context.WithValue(r.Context(), key, fields)

		next.ServeHTTP(w, r.WithContext(ctx)) // propagate new ctx
	})
}

func Selected(ctx context.Context) []string {
	v, _ := ctx.Value(key).([]string)
	return v
}
