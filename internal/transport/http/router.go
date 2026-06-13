package http

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/prajwal-huggi/backend_go/internal/auth"
	"github.com/prajwal-huggi/backend_go/internal/domain"
	"github.com/prajwal-huggi/backend_go/internal/middleware"

	httpSwagger "github.com/swaggo/http-swagger"
)

func NewRouter(userHandler *UserHandler, jwtService *auth.JWTService) http.Handler {

	r := chi.NewRouter()

	// -------- PUBLIC ROUTES --------
	r.Post("/users", userHandler.CreateUser) // signup
	r.Post("/login", userHandler.Login)
	r.Post("/refresh", userHandler.Refresh)

	// -------- PROTECTED ROUTES --------

	// r.Get("/users", middleware.JWTAuth(jwtService)(handler))
	// r.Get("/users/{id}", middleware.JWTAuth(jwtService)(handler))
	// OR
	r.Group(func(r chi.Router) {
		r.Use(middleware.JWTAuth(jwtService))

		r.With(middleware.RequireRole(domain.Manager)).Get("/users", userHandler.GetUsers)
		r.Get("/users/{id}", userHandler.GetUser)
		r.Put("/users/{id}", userHandler.UpdateUser)
		r.With(middleware.RequireRole(domain.Admin)).Delete("/users/{id}", userHandler.DeleteUser)
	})

	r.Get("/docs", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/docs/index.html", http.StatusMovedPermanently)
	})

	r.Get("/docs/*", httpSwagger.WrapHandler)
	return r
}
