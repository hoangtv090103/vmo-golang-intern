package main

import (
	// "week2-clean-architecture/internal/module_infra/memory"
	"log"
	"net/http"
	"week2-clean-architecture/internal/module_user/infra/memory"
	"week2-clean-architecture/internal/module_user/usecases"
	"week2-clean-architecture/internal/module_user/handlers"
)

func main() {
	userRepo := memory.NewUserRepoMemory()
	userUsecase := usecases.NewUserUsecase(userRepo)
	userHandler := handlers.NewUserController(userUsecase)

	// Index route
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Welcome to Clean Architecture"))
	})

	// GET /users
	http.HandleFunc("/users", userHandler.Users)

	log.Println("Server running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))

}
