package middleware

import (
	"context"
	"ecommerce/internal/auth/utils"
	"fmt"
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
		r = r.WithContext(ctx)

		role, ok := r.Context().Value("role").(string)
		if !ok {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}
		fmt.Println(role)
		next.ServeHTTP(w, r)
	})
}
