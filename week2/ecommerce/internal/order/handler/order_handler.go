package handler

import (
	"ecommerce/internal/order/domain"
	"ecommerce/internal/order/usecase"
	orderUtils "ecommerce/internal/order/utils"
	utils "ecommerce/utils"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type OrderHandler struct {
	orderUsecase   usecase.OrderUsecase
	InvoiceService *usecase.InvoiceService
}

func NewOrderHandler(orderUsecase usecase.OrderUsecase) *OrderHandler {
	return &OrderHandler{
		orderUsecase: orderUsecase,
	}
}

func (oh *OrderHandler) Handler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		oh.AddOrder(w, r)
	case http.MethodGet:
		vars := mux.Vars(r)
		if id, ok := vars["id"]; ok {
			oh.GetOrderById(w, r, id)
		} else {
			oh.GetAllOrders(w, r)
		}
	case http.MethodPut:
		oh.UpdateOrder(w, r)
	case http.MethodDelete:
		oh.DeleteOrder(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (oh *OrderHandler) AddOrder(w http.ResponseWriter, r *http.Request) {
	var order domain.Order
	if err := json.NewDecoder(r.Body).Decode(&order); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := oh.orderUsecase.CreateOrder(order); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (oh *OrderHandler) GetAllOrders(w http.ResponseWriter, r *http.Request) {
	orders, err := oh.orderUsecase.GetAllOrders()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(orders)
}

func (oh *OrderHandler) GetOrderById(w http.ResponseWriter, r *http.Request, id string) {
	orderID, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "Invalid order ID", http.StatusBadRequest)
		return
	}

	order, err := oh.orderUsecase.GetOrderById(orderID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(order)
}

func (oh *OrderHandler) UpdateOrder(w http.ResponseWriter, r *http.Request) {
	var order domain.Order
	if err := json.NewDecoder(r.Body).Decode(&order); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := oh.orderUsecase.UpdateOrder(order); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (oh *OrderHandler) DeleteOrder(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	orderID, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid order ID", http.StatusBadRequest)
		return
	}

	if err := oh.orderUsecase.DeleteOrder(orderID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *OrderHandler) PrintInvoice(w http.ResponseWriter, r *http.Request) {
	orderId := utils.StrToInt(r.URL.Query().Get("order_id"))
	invoiceData, err := h.InvoiceService.GenerateInvoiceData(orderId)
	if err != nil {
		http.Error(w, "Failed to generate invoice data", http.StatusInternalServerError)
		return
	}

	pdfBytes, err := orderUtils.GenerateInvoicePDF(*invoiceData)
	if err != nil {
		http.Error(w, "Failed to generate PDF", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition", "attachment; filename=invoice.pdf")
	w.Write(pdfBytes)
}
