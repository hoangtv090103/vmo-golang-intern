package handler

import (
	"ecommerce/config"
	"ecommerce/internal/product/domain"
	"ecommerce/internal/product/usecase"
	"ecommerce/middleware"
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
		middleware.AdminOnlyMiddleware(http.HandlerFunc(ph.AddProduct)).ServeHTTP(w, r)
	case http.MethodGet:
		if idStr := r.URL.Query().Get("id"); idStr != "" {
			ph.GetProductByID(w, r)
		} else if name := r.URL.Query().Get("name"); name != "" {
			ph.GetProductsByName(w, r)
		} else {
			ph.GetAllProducts(w)
		}
	case http.MethodPut:
		middleware.AdminOnlyMiddleware(http.HandlerFunc(ph.UpdateProduct)).ServeHTTP(w, r)
	case http.MethodDelete:
		middleware.AdminOnlyMiddleware(http.HandlerFunc(ph.DeleteProduct)).ServeHTTP(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// AddProduct handles the addition of a new product
//
//	@Summary		Add a new product
//	@Description	Add a new product with details and image
//	@Tags			products
//	@Accept			multipart/form-data
//	@Produce		json
//	@Param			name		formData	string	true	"Product Name"
//	@Param			description	formData	string	true	"Product Description"
//	@Param			price		formData	number	true	"Product Price"
//	@Param			stock		formData	int		true	"Product Stock"
//	@Param			image		formData	file	true	"Product Image"
//	@Success		201			{string}	string	"Product created successfully"
//	@Failure		400			{string}	string	"Bad Request"
//	@Failure		500			{string}	string	"Internal Server Error"
//	@Router			/products [post]
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
	w.Write([]byte("Product created successfully"))
}

// GetAllProducts handles fetching all products
//
//	@Summary		Get all products
//	@Description	Get a list of all products
//	@Tags			products
//	@Accept			json
//	@Produce		json
//	@Success		200	{array}		domain.Product
//	@Failure		500	{string}	string	"Internal Server Error"
//	@Router			/products [get]
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

// GetProductByID godoc
//
//	@Summary		Get a product by ID
//	@Description	Get a product by ID
//	@Tags			products
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int				true	"Product ID"
//	@Success		200	{object}	domain.Product	"OK"
//	@Failure		404	{string}	string			"Not found"
//	@Failure		500	{string}	string			"Internal server error"
//	@Router			/products/{id} [get]
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

// GetProductsByName handles fetching products by name
//
//	@Summary		Get products by name
//	@Description	Get a list of products by name
//	@Tags			products
//	@Accept			json
//	@Produce		json
//	@Param			name	query		string	true	"Product Name"
//	@Success		200		{array}		domain.Product
//	@Failure		400		{string}	string	"Bad Request"
//	@Failure		500		{string}	string	"Internal Server Error"
//	@Router			/products [get]
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

// UpdateProduct handles updating an existing product
//
//	@Summary		Update a product
//	@Description	Update details of an existing product
//	@Tags			products
//	@Accept			multipart/form-data
//	@Produce		json
//	@Param			id			query		int		true	"Product ID"
//	@Param			name		formData	string	false	"Product Name"
//	@Param			description	formData	string	false	"Product Description"
//	@Param			price		formData	number	false	"Product Price"
//	@Param			stock		formData	int		false	"Product Stock"
//	@Param			image		formData	file	false	"Product Image"
//	@Success		200			{string}	string	"Updated Product with id"
//	@Failure		400			{string}	string	"Bad Request"
//	@Failure		500			{string}	string	"Internal Server Error"
//	@Router			/products [put]
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

// DeleteProduct handles deleting a product by its ID
//
//	@Summary		Delete a product
//	@Description	Delete a product by its ID
//	@Tags			products
//	@Accept			json
//	@Produce		json
//	@Param			id	query		int		true	"Product ID"
//	@Success		200	{string}	string	"Product deleted successfully"
//	@Failure		400	{string}	string	"Bad Request"
//	@Failure		500	{string}	string	"Internal Server Error"
//	@Router			/products [delete]
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
