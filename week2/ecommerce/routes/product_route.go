package routes

import (
	"ecommerce/internal/product/handler"
	"net/http"
)

func ProductRoute(h *handler.ProductHandler) {
	http.HandleFunc("/products", h.Handler)
}
