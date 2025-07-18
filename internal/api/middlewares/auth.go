package middlewares

import (
	"ecommerce-go/pkg/utils"
	"fmt"
	"net/http"
)

func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("access_token")
		if err != nil {
			utils.WriteError(w, http.StatusUnauthorized, "Unauthorized: missing token")
			return
		}
		fmt.Println("Token: ", cookie.Value)

		next.ServeHTTP(w, r)
	})
}
