package middleware

import (
	"context"
	"ecommerce/internal/auth/utils"
	"net/http"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenStr := r.Header.Get("Authorization")

		if tokenStr == "" {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}

		claims, err := utils.ValidateJWT(tokenStr)

		if err != nil {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}

		ctx := context.WithValue(r.Context(), "username", claims.Username)
		ctx = context.WithValue(r.Context(), "role", claims.Role)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
