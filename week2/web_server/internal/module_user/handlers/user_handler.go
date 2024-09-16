package handlers

import (
	"encoding/json"
	"net/http"
	"week2-clean-architecture/internal/module_user/domain"
	"week2-clean-architecture/internal/module_user/usecases"
)

type UserHanlder struct {
	usecase *usecases.UserUseCase
}

func NewUserController(u *usecases.UserUseCase) *UserHanlder {
	return &UserHanlder{usecase: u}
}

func (h *UserHanlder) Users(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		{
			if r.URL.Query().Get("username") != "" {
				h.GetUserByUsername(w, r)
			} else {
				h.GetUsers(w, r)
			}
		}
	case "POST":
		h.CreateUser(w, r)
	}
}

func (h *UserHanlder) GetUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.usecase.GetUsers()

	if err != nil {
		http.Error(w, "Unable to get users", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(users)
}

func (h *UserHanlder) GetUserByUsername(w http.ResponseWriter, r *http.Request) {
	username := r.URL.Query().Get("username")

	user, err := h.usecase.GetUserByUsername(username)

	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(user)
}

func (h *UserHanlder) CreateUser(w http.ResponseWriter, r *http.Request) {
	var user domain.User

	json.NewDecoder(r.Body).Decode(&user)

	err := h.usecase.CreateUser(&user)

	if err != nil {
		http.Error(w, "Unable to create user", http.StatusInternalServerError)
		return
	}
}
