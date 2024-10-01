package handler

import (
	"encoding/json"
	"net/http"

	authDomain "ecommerce/internal/auth/domain"
	authUsecase "ecommerce/internal/auth/usecase"
)

// AuthHandler handles authentication related requests
type AuthHandler struct {
	authUseCase authUsecase.AuthUsecase
}

// NewAuthHandler creates a new AuthHandler
func NewAuthHandler(authUsecase *authUsecase.AuthUsecase) *AuthHandler {
	return &AuthHandler{
		authUseCase: *authUsecase,
	}
}

// Login godoc
//
//	@Summary		Login a user
//	@Description	Authenticate a user and return a JWT token
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Success		200	{string}	string	"token"
//	@Failure		405	{string}	string	"Method not allowed"
//	@Failure		500	{string}	string	"Internal server error"
//	@Router			/auth/login [post]
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var auth authDomain.Auth
	var err error

	// Decode request body to auth struct
	err = json.NewDecoder(r.Body).Decode(&auth)

	token, err := h.authUseCase.Login(auth.GetUsername(), auth.GetPassword())

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = w.Write([]byte(token))
	if err != nil {
		return
	}
}

// Register godoc
//
//	@Summary		Register a new user
//	@Description	Register a new user with the provided details
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Success		201	{string}	string	"Created"
//	@Failure		400	{string}	string	"Bad request"
//	@Failure		500	{string}	string	"Internal server error"
//	@Router			/auth/register [post]
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var auth authDomain.Auth
	var err error

	// Decode request body to auth struct
	err = json.NewDecoder(r.Body).Decode(&auth)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Register the new user
	err = h.authUseCase.Register(auth)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
