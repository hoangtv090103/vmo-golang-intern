package routes

import (
	"ecommerce/internal/auth/handler"
	"net/http"
)

func AuthRoute(h *handler.AuthHandler, mux *http.ServeMux) {
	mux.HandleFunc("/login", h.Login)
	mux.HandleFunc("/register", h.Register)
}
