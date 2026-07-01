package main

import (
	"context"
	"log"

	"github.com/devharnold/smart-reconcile/internal/auth"
	"github.com/devharnold/smart-reconcile/internal/merchants"
	"github.com/devharnold/smart-reconcile/internal/reconciler"
	"github.com/devharnold/smart-reconcile/internal/router"
	"github.com/devharnold/smart-reconcile/internal/scheduler"
	"github.com/devharnold/smart-reconcile/internal/storage"
)

func main() {
	ctx := context.Background()

	reconciliation := reconciler.Engine{}
	s := scheduler.NewScheduler()

	err := s.Register(
		scheduler.Schedule{
			Frequency: scheduler.Daily,

			Hour:   6,
			Minute: 0,
		}, func() {
			reconciliation.Reconcile()
		},
	)
	// Init DB
	db, err := storage.NewDB(ctx)
	if err != nil {
		log.Fatal(err)
	}
	s.Start()
	log.Println("Scheduler started")
	defer db.Close()

	// Init services
	jwtSvc := auth.NewJWTService()

	merchantRepo := merchants.NewMerchantsRepository(db)
	merchantSvc := merchants.NewMerchantService(merchantRepo, jwtSvc)
	merchantHandler := merchants.NewMerchantHandler(merchantSvc)

	// Build router
	r := router.SetUpRouter(merchantHandler, jwtSvc)

	// Start server
	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
