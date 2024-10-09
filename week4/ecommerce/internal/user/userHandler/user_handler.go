package userHandler

import (
	"ecommerce/internal/user/entity"
	"ecommerce/internal/user/usecase"
	"github.com/gofiber/fiber/v2"
	"strconv"
)

type UserHandler struct {
	uc usecase.UserUsecase
}

func NewUserHandler(uc usecase.UserUsecase) *UserHandler {
	return &UserHandler{
		uc: uc,
	}
}

func (uh *UserHandler) AddUser(c *fiber.Ctx) error {
	var user entity.User
	var err error

	// Decode request body to user struct
	if err := c.BodyParser(&user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	err = uh.uc.CreateUser(c.Context(), &user)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "User created successfully"})
}

// GetAllUsers handles retrieving all users
//
//	@Summary		Get all users
//	@Description	Retrieve a list of all users
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Success		200	{array}		entity.User
//	@Failure		500	{string}	string	"Internal Server Error"
//	@Router			/users [get]
func (uh *UserHandler) GetAllUsers(c *fiber.Ctx) error {
	var users []*entity.User
	var err error

	users, err = uh.uc.GetAllUsers(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(users)
}

// UpdateUser handles updating an existing user
//
//	@Summary		Update a user
//	@Description	Update details of an existing user
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			user	body		entity.User	true	"User data"
//	@Success		200		{string}	string		"User updated successfully"
//	@Failure		400		{string}	string		"Bad Request"
//	@Failure		500		{string}	string		"Internal Server Error"
//	@Router			/users [put]
func (uh *UserHandler) UpdateUser(c *fiber.Ctx) error {
	var user entity.User
	var err error
	idStr := c.Params("id")
	if idStr == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "ID is required"})
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid ID"})
	}

	user.ID = id
	// Decode request body to user struct
	if err = c.BodyParser(user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	err = uh.uc.UpdateUser(c.Context(), &user)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "Account updated successfully", "user": user})
}

// DeleteUser handles deleting a user by its ID
//
//	@Summary		Delete a user
//	@Description	Delete a user by its ID
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			id	query		int		true	"User ID"
//	@Success		200	{string}	string	"User deleted successfully"
//	@Failure		400	{string}	string	"Bad Request"
//	@Failure		500	{string}	string	"Internal Server Error"
//	@Router			/users [delete]
func (uh *UserHandler) DeleteUser(c *fiber.Ctx) error {
	var err error

	idStr := c.Query("id")
	if idStr == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "ID is required"})
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid ID"})
	}

	err = uh.uc.DeleteUser(c.Context(), id)
	if err != nil {
		c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		return nil
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "User deleted successfully"})
}
