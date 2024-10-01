package handler

import (
	"ecommerce/config"
	"ecommerce/internal/order/domain"
	"ecommerce/internal/order/usecase"
	orderUtils "ecommerce/internal/order/utils"
	"ecommerce/middleware"
	utils "ecommerce/utils"
	"encoding/json"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
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
		middleware.UserOnlyMiddleware(http.HandlerFunc(oh.AddOrder)).ServeHTTP(w, r)
	case http.MethodGet:
		vars := mux.Vars(r)
		if role, ok := r.Context().Value("role").(string); ok && role == "admin" {
			if id, ok := vars["id"]; ok {
				oh.GetOrderById(w, id)
			} else {
				oh.GetAllOrders(w)
			}
		} else {
			username, ok := r.Context().Value("username").(string)
			if !ok {
				return
			}

			oh.GetUserOrders(w, username)
		}
	case http.MethodPut:
		middleware.AdminOnlyMiddleware(http.HandlerFunc(oh.UpdateOrder)).ServeHTTP(w, r)
	case http.MethodDelete:
		middleware.AdminOnlyMiddleware(http.HandlerFunc(oh.DeleteOrder)).ServeHTTP(w, r)
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

func (oh *OrderHandler) GetAllOrders(w http.ResponseWriter) {
	orders, err := oh.orderUsecase.GetAllOrders()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(orders)
}

func (oh *OrderHandler) GetOrderById(w http.ResponseWriter, id string) {
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

func (oh *OrderHandler) GetUserOrders(w http.ResponseWriter, username string) {
	orders, err := oh.orderUsecase.GetUserOrders(username)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(orders)
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

	// Create a temporary file to store the PDF
	tempFile, err := os.CreateTemp("", "invoice-*.pdf")
	if err != nil {
		http.Error(w, "Failed to create temp file", http.StatusInternalServerError)
		return
	}
	defer os.Remove(tempFile.Name())
	defer tempFile.Close()

	// Write the PDF bytes to the temporary file
	if _, err := tempFile.Write(pdfBytes); err != nil {
		http.Error(w, "Failed to write PDF to temp file", http.StatusInternalServerError)
		return
	}

	// Reopen the file for reading
	tempFile, err = os.Open(tempFile.Name())
	if err != nil {
		http.Error(w, "Failed to reopen temp file", http.StatusInternalServerError)
		return
	}
	defer tempFile.Close()

	// Upload Invoice to S3
	s3Config := config.NewS3Config(
		os.Getenv("AWS_REGION"),
		os.Getenv("AWS_BUCKET_NAME"),
		os.Getenv("AWS_ACCESS_KEY_ID"),
		os.Getenv("AWS_SECRET_ACCESS_KEY"),
	)

	fileHeader := &multipart.FileHeader{
		Filename: filepath.Base(tempFile.Name()),
		Size:     int64(len(pdfBytes)),
	}

	pdfPath, err := utils.UploadFileToS3(
		*s3Config,
		tempFile,
		fileHeader,
		"invoices",
	)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write([]byte(pdfPath))
}
