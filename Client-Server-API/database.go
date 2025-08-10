package main

import (
	"fmt"
	"log"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
	var err error
	DB, err = gorm.Open(sqlite.Open("cotacoes.db"), &gorm.Config{})
	if err != nil {
		log.Fatalf("Erro ao abrir banco: %v", err)
	}

	if err := DB.AutoMigrate(&Cotacao{}); err != nil {
		log.Fatalf("Erro ao migrar schema: %v", err)
	}

	sqlDB, err := DB.DB()
	if err != nil {
		log.Fatalf("Erro ao obter DB SQL: %v", err)
	}
	sqlDB.SetMaxOpenConns(1)
	sqlDB.SetConnMaxLifetime(time.Hour)
	fmt.Println("Banco inicializado com sucesso.")
}
