package http

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func NewRouter(userHandler *UserHandler) http.Handler {

	r := chi.NewRouter()

	r.Post("/users", userHandler.CreateUser)
	r.Get("/users", userHandler.GetUsers)
	r.Get("/users/{id}", userHandler.GetUser)
	r.Put("/users/{id}", userHandler.UpdateUser)
	r.Delete("/users/{id}", userHandler.DeleteUser)

	return r
}