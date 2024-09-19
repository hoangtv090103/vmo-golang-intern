package handler

import (
	"ecommerce/internal/user/domain"
	"ecommerce/internal/user/usecase"
	"ecommerce/utils"
	"encoding/json"
	"net/http"
)

type UserHandler struct {
	userUseCase usecase.UserUseCase
}

func NewUserHandler(userUseCase *usecase.UserUseCase) *UserHandler {
	return &UserHandler{
		userUseCase: *userUseCase,
	}
}

func (uh *UserHandler) Handler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		uh.AddUser(w, r)
	case http.MethodGet:
		uh.GetAllUsers(w, r)
	case http.MethodPut:
		uh.UpdateUser(w, r)
	case http.MethodDelete:
		uh.DeleteUser(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (uh *UserHandler) AddUser(w http.ResponseWriter, r *http.Request) {
	var user domain.User
	var err error

	// Decode request body to user struct
	err = json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	err = uh.userUseCase.AddUser(user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (uh *UserHandler) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	var users []domain.User
	var err error

	users, err = uh.userUseCase.GetAllUsers()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// If no users found: users: null
	if len(users) == 0 {
		users = []domain.User{}
	}
	// Encode users to JSON
	err = json.NewEncoder(w).Encode(users)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
}

func (uh *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	var user domain.User
	var err error

	// Decode request body to user struct
	err = json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = uh.userUseCase.UpdateUser(user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (uh *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	var err error

	// Get id from URL
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		http.Error(w, "ID is required", http.StatusBadRequest)
		return
	}

	// Convert id to int
	id := utils.StrToInt(idStr)

	err = uh.userUseCase.DeleteUser(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
