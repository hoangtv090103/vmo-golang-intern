package routes

import (
	"ecommerce/internal/product/handler"
	"net/http"
)

func ProductRoute(h *handler.ProductHandler, mux *http.ServeMux) {
	mux.HandleFunc("/products", h.Handler)
}
