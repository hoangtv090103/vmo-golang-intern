package routes

import (
	"ecommerce/internal/user/handler"
	"net/http"
)

func UserRoute(h *handler.UserHandler, mux *http.ServeMux) {
	mux.HandleFunc("/users", h.Handler)
}
