package handler

import (
	authDomain "ecommerce/internal/auth/domain"
	authUsecase "ecommerce/internal/auth/usecase"
	"encoding/json"
	"net/http"
)

type AuthHandler struct {
	authUseCase authUsecase.AuthUsecase
}

func NewAuthHandler(authUsecase *authUsecase.AuthUsecase) *AuthHandler {
	return &AuthHandler{
		authUseCase: *authUsecase,
	}
}

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

//func (h *AuthHandler) ForgetPassword(w http.ResponseWriter, r *http.Request) {
//	username := r.URL.Query().Get("username")
//
//	err := h.authUseCase.ForgetPassword(username)
//
//	if err != nil {
//		http.Error(w, err.Error(), http.StatusInternalServerError)
//		return
//	}
//}

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
