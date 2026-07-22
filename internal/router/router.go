package router

import (
	"net/http"

	"github.com/devharnold/smart-reconcile/internal/auth"
	"github.com/devharnold/smart-reconcile/internal/authmiddleware"
	"github.com/devharnold/smart-reconcile/internal/merchants"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func setUpRouter(handler *merchants.Handler, jwtSvc auth.JWTService) http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Public Routes
	r.Post("/signup", handler.SignUp)
	r.Post("/login", handler.Login)

	// Protected Routes
	r.Group(func(r chi.Router) {
		r.Use(authmiddleware.JWTAuth(jwtSvc))

		r.Post("/search-user", handler.SearchUser)
	})

	return r
}
