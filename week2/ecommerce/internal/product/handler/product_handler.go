package handler

import (
	"ecommerce/internal/product/domain"
	"ecommerce/internal/product/usecase"
	"ecommerce/utils"
	"encoding/json"
	"net/http"
	"strconv"
)

type ProductHandler struct {
	productUseCase usecase.ProductUseCase
}

func NewProductHandler(productUseCase *usecase.ProductUseCase) *ProductHandler {
	return &ProductHandler{
		productUseCase: *productUseCase,
	}
}

func (ph *ProductHandler) Handler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		ph.AddProduct(w, r)
	case http.MethodGet:
		if idStr := r.URL.Query().Get("id"); idStr != "" {
			ph.GetProductByID(w, r)
		} else {
			ph.GetAllProducts(w, r)
		}
	case http.MethodPut:
		ph.UpdateProduct(w, r)
	case http.MethodDelete:
		ph.DeleteProduct(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (ph *ProductHandler) AddProduct(w http.ResponseWriter, r *http.Request) {
	var product domain.Product
	var err error

	//Decode request body to product struct
	err = json.NewDecoder(r.Body).Decode(&product)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	err = ph.productUseCase.AddProduct(product)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (ph *ProductHandler) GetAllProducts(w http.ResponseWriter, r *http.Request) {
	var products []domain.Product
	var err error

	products, err = ph.productUseCase.GetAllProducts()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	//Encode response
	err = json.NewEncoder(w).Encode(products)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (ph *ProductHandler) GetProductByID(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	product, err := ph.productUseCase.GetProductByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	//Encode response
	err = json.NewEncoder(w).Encode(product)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
}

func (ph *ProductHandler) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id := utils.StrToInt(idStr)
	if id == 0 {
		http.Error(w, "ID is required", http.StatusBadRequest)
		return
	}

	var product domain.Product
	var err error

	//Decode request body to product struct
	err = json.NewDecoder(r.Body).Decode(&product)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = ph.productUseCase.UpdateProduct(product)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (ph *ProductHandler) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id := utils.StrToInt(idStr)
	if id == 0 {
		http.Error(w, "ID is required", http.StatusBadRequest)
		return
	}

	err := ph.productUseCase.DeleteProduct(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
