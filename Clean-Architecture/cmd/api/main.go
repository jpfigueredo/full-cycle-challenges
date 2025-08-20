package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jpfigueredo/full-cycle-challenges/Clean-Architecture/internal/api"
	gql "github.com/jpfigueredo/full-cycle-challenges/Clean-Architecture/internal/api/graphql"
	grp "github.com/jpfigueredo/full-cycle-challenges/Clean-Architecture/internal/api/grpc"
	"github.com/jpfigueredo/full-cycle-challenges/Clean-Architecture/internal/config"
	"github.com/jpfigueredo/full-cycle-challenges/Clean-Architecture/internal/domain"
	"github.com/jpfigueredo/full-cycle-challenges/Clean-Architecture/internal/repository"
	"github.com/jpfigueredo/full-cycle-challenges/Clean-Architecture/internal/service"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func connectWithRetry(dsn string, retries int, delay time.Duration) (*gorm.DB, error) {
	var db *gorm.DB
	var err error
	for i := 0; i < retries; i++ {
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err == nil {
			return db, nil
		}
		log.Printf("❌ failed to connect to DB, retrying in %s... (%d/%d)", delay, i+1, retries)
		time.Sleep(delay)
	}
	return nil, err
}

func main() {
	cfg := config.Load()

	// DSN for Postgres
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
	)

	db, err := connectWithRetry(dsn, 10, 2*time.Second)
	if err != nil {
		log.Fatalf("❌ failed to initialize DB: %v", err)
	}

	if err := db.AutoMigrate(&domain.Order{}, &domain.Patient{}); err != nil {
		log.Fatalf("❌ failed to migrate: %v", err)
	}

	orderRepo := repository.NewOrderRepository(db)
	orderService := service.NewOrderService(orderRepo)

	// REST
	go func() {
		r := api.SetupRouterWithServices(orderService, db)
		log.Printf("🚀 REST server running on port %s", cfg.ServerPort)
		if err := r.Run(":" + cfg.ServerPort); err != nil {
			log.Fatalf("❌ failed to start REST server: %v", err)
		}
	}()

	// gRPC
	go func() {
		grp.StartGRPCServer(orderService, "50051")
	}()

	// GraphQL
	go func() {
		resolver := &gql.Resolver{OrderService: orderService}
		gql.StartGraphQLServer(resolver, "8081")
	}()

	// block forever
	select {}
}
