package routes

import (
	"ecommerce/internal/order/handler"
	"net/http"
)

func OrderRoute(handler *handler.OrderHandler, mux *http.ServeMux) {
	mux.HandleFunc("/orders", handler.Handler)
}
