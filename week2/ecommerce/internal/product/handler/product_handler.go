package handler

import (
	"ecommerce/config"
	"ecommerce/internal/product/domain"
	"ecommerce/internal/product/usecase"
	"ecommerce/utils"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
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
		} else if name := r.URL.Query().Get("name"); name != "" {
			ph.GetProductsByName(w, r)

		} else {
			ph.GetAllProducts(w)
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

	// Parse multipart form (2MB limit)
	err = r.ParseMultipartForm(2 << 20) // 2MB
	if err != nil {
		http.Error(w, "Failed to parse form data", http.StatusBadRequest)
		return
	}

	// Retrieve image file from the form
	file, handler, err := r.FormFile("image")
	if err != nil {
		http.Error(w, "Failed to retrieve file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Validate the file type (only jpg and png allowed)
	//ext := strings.ToLower(filepath.Ext(handler.Filename))
	//if ext != ".jpg" && ext != ".jpeg" && ext != ".png" {
	//	http.Error(w, "Invalid file type. Only .jpg and .png are allowed", http.StatusUnsupportedMediaType)
	//	return
	//}
	//
	//// Validate file size (max 2MB)
	//if handler.Size > 2*1024*1024 {
	//	http.Error(w, "File too large (maximum 2MB)", http.StatusBadRequest)
	//	return
	//}
	//
	//// Ensure the uploads directory exists
	//uploadsDir := "uploads"
	//if _, err := os.Stat(uploadsDir); os.IsNotExist(err) {
	//	err = os.MkdirAll(uploadsDir, 0755)
	//	if err != nil {
	//		http.Error(w, "Failed to create uploads directory", http.StatusInternalServerError)
	//		return
	//	}
	//}
	//
	//// Save the image with a unique filename
	//filePath := filepath.Join(uploadsDir, handler.Filename)
	//dst, err := os.Create(filePath)
	//if err != nil {
	//	http.Error(w, "Failed to create file", http.StatusInternalServerError)
	//	return
	//}
	//defer dst.Close()
	//
	//// Copy the uploaded file to the destination file
	//_, err = io.Copy(dst, file)
	//if err != nil {
	//	http.Error(w, "Failed to save file", http.StatusInternalServerError)
	//	return
	//}

	s3Config := config.NewS3Config(
		os.Getenv("AWS_REGION"),
		os.Getenv("AWS_BUCKET_NAME"),
		os.Getenv("AWS_ACCESS_KEY_ID"),
		os.Getenv("AWS_SECRET_ACCESS_KEY"),
	)

	var filePath string
	filePath, err = utils.UploadFileToS3(*s3Config, file, handler, "product_images")

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	// Store the file path in the product struct
	product.SetImagePath(filePath)

	// Set other product fields
	product.SetName(r.FormValue("name"))
	product.SetDescription(r.FormValue("description"))
	price, _ := strconv.ParseFloat(r.FormValue("price"), 64)
	product.SetPrice(price)
	stock, _ := strconv.Atoi(r.FormValue("stock"))
	product.SetStock(stock)

	// Save product to the database
	err = ph.productUseCase.AddProduct(product)
	if err != nil {
		http.Error(w, "Failed to save product", http.StatusInternalServerError)
		return
	}

	// Respond with success
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Product and image uploaded successfully"))
}

func (ph *ProductHandler) GetAllProducts(w http.ResponseWriter) {
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

func (ph *ProductHandler) GetProductsByName(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")

	var products []domain.Product
	var err error

	products, err = ph.productUseCase.GetProductsByName(name)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = json.NewEncoder(w).Encode(products)

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

	product.ID = id
	product.SetName(r.FormValue("name"))
	product.SetDescription(r.FormValue("description"))
	price, _ := strconv.ParseFloat(r.FormValue("price"), 64)
	product.SetPrice(price)
	stock, _ := strconv.Atoi(r.FormValue("stock"))
	product.SetStock(stock)

	file, handler, err := r.FormFile("image")

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if handler == nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	defer file.Close()

	// Check file size
	if handler.Size > 2*1024*1024 {
		http.Error(w, "File size exceeds 2MB", http.StatusBadRequest)
		return
	}

	/// Create a new file in the uploads directory
	s3Config := config.NewS3Config(
		os.Getenv("AWS_REGION"),
		os.Getenv("AWS_BUCKET_NAME"),
		os.Getenv("AWS_ACCESS_KEY_ID"),
		os.Getenv("AWS_SECRET_ACCESS_KEY"),
	)

	var filePath string
	filePath, err = utils.UploadFileToS3(*s3Config, file, handler, "product_image")

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	// Store the file path in the product struct
	product.SetImagePath(filePath)
	err = ph.productUseCase.UpdateProduct(product)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("Updated Product with id (%d)", id)))
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
