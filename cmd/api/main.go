package main

import (
	"context"
	"log"

	"github.com/devharnold/smart-reconcile/internal/auth"
	"github.com/devharnold/smart-reconcile/internal/storage"
	"github.com/devharnold/smart-reconcile/internal/users"
)

func main() {
	ctx := context.Background()

	// Init DB
	db, err := storage.NewDB(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Init services
	jwtSvc := auth.NewJWTService()

	userRepo := users.NewUsersRepository(db)
	userSvc := users.NewUserService(userRepo, jwtSvc)
	userHandler := users.NewUserHandler(userSvc)

	// Build router
	r := router.SetUpRouter(userHandler, jwtSvc)

	// Start server
	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
