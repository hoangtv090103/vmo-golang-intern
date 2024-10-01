package middleware

import (
	"net/http"
)

// AdminOnlyMiddleware restricts access to admins
func AdminOnlyMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		role, ok := r.Context().Value("role").(string)

		if !ok || role != "admin" {
			http.Error(w, "Forbidden: Admin access required", http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// UserOnlyMiddleware restricts access to users
func UserOnlyMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		role := r.Context().Value("role").(string)
		if role != "user" {
			http.Error(w, "Forbidden: User access required", http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r)
	})
}
