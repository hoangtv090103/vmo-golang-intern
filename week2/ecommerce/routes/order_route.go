package routes

import (
	"ecommerce/internal/order/handler"
	"net/http"
)

func OrderRoute(handler *handler.OrderHandler) {
	http.HandleFunc("/orders", handler.Handler)
}
