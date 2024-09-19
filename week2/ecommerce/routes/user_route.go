package routes

import (
	"ecommerce/internal/user/handler"
	"net/http"
)

func UserRoute(h *handler.UserHandler) {
	http.HandleFunc("/users", h.Handler)
}
