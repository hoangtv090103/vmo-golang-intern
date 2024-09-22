package routes

import (
	"ecommerce/internal/reports"
	"net/http"
)

func ReportRoute(handler *reports.ReportHandler, mux *http.ServeMux) {
	mux.HandleFunc("/reports", handler.HandleReport)
}
