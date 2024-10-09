package handler

import (
	"ecommerce/internal/product/entity"
	"ecommerce/internal/product/usecase"
	"github.com/gofiber/fiber/v2"
	"strconv"
)

type ProductHandler struct {
	uc usecase.ProductUsecase
}

func NewProductHandler(uc usecase.ProductUsecase) *ProductHandler {
	return &ProductHandler{
		uc: uc,
	}
}

func (ph *ProductHandler) AddProduct(c *fiber.Ctx) error {
	var product entity.Product
	var err error

	// Decode request body to product struct
	if err := c.BodyParser(&product); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	err = ph.uc.CreateProduct(c.Context(), &product)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "Product created successfully"})
}

func (ph *ProductHandler) GetAllProducts(c *fiber.Ctx) error {
	var products []*entity.Product
	var err error

	products, err = ph.uc.GetAllProducts(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(products)
}

func (ph *ProductHandler) GetProductByID(c *fiber.Ctx) error {
	var product *entity.Product
	var err error

	idStr := c.Params("id")
	if idStr == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "ID is required"})
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid ID"})
	}

	product, err = ph.uc.GetByProductID(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(product)
}

func (ph *ProductHandler) UpdateProduct(c *fiber.Ctx) error {
	var product *entity.Product
	var err error

	// Decode request body to product struct
	if err = c.BodyParser(product); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	err = ph.uc.UpdateProduct(c.Context(), product)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "Product updated successfully", "product": product})
}

func (ph *ProductHandler) DeleteProduct(c *fiber.Ctx) error {
	var err error

	idStr := c.Query("id")
	if idStr == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "ID is required"})
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid ID"})
	}

	err = ph.uc.DeleteProduct(c.Context(), id)
	if err != nil {
		c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		return nil
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Product deleted successfully"})
}
