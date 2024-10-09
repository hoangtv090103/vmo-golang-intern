package handler

import (
	"ecommerce/internal/order/entity"
	"ecommerce/internal/order/usecase"
	utils "ecommerce/internal/order/utils"
	"ecommerce/pkg/config"
	globalUtils "ecommerce/pkg/utils"
	"github.com/gofiber/fiber/v2"
	"mime/multipart"
	"os"
	"path/filepath"
	"strconv"
)

type OrderHandler struct {
	orderUsecase *usecase.OrderUsecase
}

func NewOrderHandler(orderUsecase *usecase.OrderUsecase) *OrderHandler {
	return &OrderHandler{
		orderUsecase: orderUsecase,
	}
}

func (h *OrderHandler) CreateOrder(c *fiber.Ctx) error {
	var order entity.Order
	if err := c.BodyParser(&order); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	if err := h.orderUsecase.CreateOrder(c.Context(), &order); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(order)
}

func (h *OrderHandler) GetAllOrders(c *fiber.Ctx) error {
	orders, err := h.orderUsecase.GetAllOrders(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(orders)
}

func (h *OrderHandler) GetOrderByID(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid ID"})
	}

	order, err := h.orderUsecase.GetOrderByID(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(order)
}

func (h *OrderHandler) GetUserOrders(c *fiber.Ctx) error {
	username := c.Params("username")
	var (
		err    error
		orders []*entity.Order
	)
	if username == "" {
		orders, err = h.orderUsecase.GetAllOrders(c.Context())
	} else {
		orders, err = h.orderUsecase.GetUserOrders(c.Context(), username)

	}
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(orders)
}

func (h *OrderHandler) UpdateOrder(c *fiber.Ctx) error {
	var order entity.Order
	if err := c.BodyParser(&order); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	if err := h.orderUsecase.UpdateOrder(c.Context(), &order); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(order)
}

func (h *OrderHandler) DeleteOrder(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid ID"})
	}

	if err := h.orderUsecase.DeleteOrder(c.Context(), id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusNoContent).Send(nil)
}

func (h *OrderHandler) GetInvoice(c *fiber.Ctx) error {
	orderIDStr := c.Params("id")
	orderID, err := strconv.Atoi(orderIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid order ID"})
	}

	invoices, err := h.orderUsecase.GetInvoice(c.Context(), orderID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(invoices)
}

func (h *OrderHandler) PrintInvoice(c *fiber.Ctx) error {
	orderIdStr := c.Params("id")
	orderId, err := strconv.Atoi(orderIdStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	invoiceData, err := h.orderUsecase.GetInvoice(c.Context(), orderId)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to generate invoice data"})
	}

	pdfBytes, err := utils.GenerateInvoicePDF(invoiceData)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err,
		})
	}

	// Create a temporary file to store the PDF
	tempFile, err := os.CreateTemp("", "invoice-*.pdf")
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err,
		})
	}
	defer os.Remove(tempFile.Name())
	defer tempFile.Close()

	// Write the PDF bytes to the temporary file
	if _, err := tempFile.Write(pdfBytes); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to write PDF to temp file",
		})
	}

	// Reopen the file for reading
	tempFile, err = os.Open(tempFile.Name())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to reopen temp file",
		})
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

	pdfPath, err := globalUtils.UploadFileToS3(
		*s3Config,
		tempFile,
		fileHeader,
		"invoices",
	)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"pdf_path": pdfPath})
}
