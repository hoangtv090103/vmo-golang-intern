package handler

import (
	"ecommerce/internal/account/entity"
	"ecommerce/internal/account/usecase"

	"github.com/gofiber/fiber/v2"
)

type AccountHandler struct {
	au usecase.IAccountUsecase
}

func NewAccountHandler(au usecase.IAccountUsecase) *AccountHandler {
	return &AccountHandler{
		au: au,
	}
}

func (h *AccountHandler) Register(c *fiber.Ctx) error {
    var account entity.Account
	if err := c.BodyParser(&account); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
    
	err := h.au.Register(c.Context(), &account)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
    
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "Account created successfully"})
}

func (h *AccountHandler) Login(c *fiber.Ctx) error {
	var loginRequest struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := c.BodyParser(&loginRequest); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	account, err := h.au.Login(c.Context(), loginRequest.Username, loginRequest.Password)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(account)
}
