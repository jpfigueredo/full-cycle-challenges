package api

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/jpfigueredo/full-cycle-challenges/Clean-Architecture/internal/api/handler"
	"github.com/jpfigueredo/full-cycle-challenges/Clean-Architecture/internal/domain"
	"github.com/jpfigueredo/full-cycle-challenges/Clean-Architecture/internal/repository"
	"github.com/jpfigueredo/full-cycle-challenges/Clean-Architecture/internal/service"
)

func ConnectDatabaseWithRetry(dsn string, retries int, delay time.Duration) (*gorm.DB, error) {
	var db *gorm.DB
	var err error
	for i := 0; i < retries; i++ {
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err == nil {
			return db, nil
		}
		log.Printf("failed to connect to database, retrying in %s... (%d/%d)", delay, i+1, retries)
		time.Sleep(delay)
	}
	return nil, err
}

func SetupRouter() *gin.Engine {
	router := gin.Default()

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
	)

	db, err := ConnectDatabaseWithRetry(dsn, 10, 2*time.Second)
	if err != nil {
		log.Fatalf("failed to initialize database: %v", err)
	}

	if err := db.AutoMigrate(&domain.Order{}, &domain.Patient{}); err != nil {
		panic("failed to migrate database")
	}

	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	orderRepo := repository.NewOrderRepository(db)
	orderService := service.NewOrderService(orderRepo)
	orderHandler := handler.NewOrderHandler(orderService)

	router.GET("/orders", orderHandler.GetOrders)
	router.POST("/orders", orderHandler.CreateOrder)
	router.GET("/orders/:id", orderHandler.GetOrderByID)
	router.PUT("/orders/:id", orderHandler.UpdateOrder)
	router.DELETE("/orders/:id", orderHandler.DeleteOrder)

	patientRepo := repository.NewPatientRepository(db)
	patientService := service.NewPatientService(patientRepo)
	patientHandler := handler.NewPatientHandler(patientService)

	router.GET("/patients", patientHandler.GetPatients)
	router.POST("/patients", patientHandler.CreatePatient)
	router.GET("/patients/:id", patientHandler.GetPatientByID)
	router.PUT("/patients/:id", patientHandler.UpdatePatient)
	router.DELETE("/patients/:id", patientHandler.DeletePatient)

	return router
}
