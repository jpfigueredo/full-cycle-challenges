package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var db *gorm.DB

type Cotacao struct {
	ID        uint `gorm:"primaryKey"`
	Bid       string
	Timestamp time.Time
}

type CotacaoResponse struct {
	USDBRL struct {
		Bid string `json:"bid"`
	} `json:"USDBRL"`
}

func initDB() {
	var err error
	db, err = gorm.Open(sqlite.Open("cotacao.db"), &gorm.Config{})
	if err != nil {
		fmt.Println("Erro ao conectar ao banco de dados:", err)
	}
	db.AutoMigrate(&Cotacao{})
}

func GeraCotacao(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", "https://economia.awesomeapi.com.br/json/last/USD-BRL", nil)
	if err != nil {
		fmt.Println("Erro ao criar requisição:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "erro ao criar requisição"})
		return
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("Erro ao buscar cotação:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "erro ao buscar cotação"})
		return
	}
	defer resp.Body.Close()

	var cotacao CotacaoResponse
	if err := json.NewDecoder(resp.Body).Decode(&cotacao); err != nil {
		fmt.Println("Erro ao decodificar JSON:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "erro ao decodificar resposta"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"bid": cotacao.USDBRL.Bid})

	dbCtx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()

	registro := Cotacao{
		Bid:       cotacao.USDBRL.Bid,
		Timestamp: time.Now(),
	}

	done := make(chan error, 1)

	go func() {
		err := db.WithContext(dbCtx).Create(&registro).Error
		done <- err
	}()

	select {
	case <-dbCtx.Done():
		fmt.Println("Timeout ao salvar no banco")
	case err := <-done:
		if err != nil {
			fmt.Println("Erro ao salvar no banco:", err)
		}
	}

}

func main() {
	initDB()
	routes := gin.Default()
	routes.GET("/cotacao", GeraCotacao)
	routes.Run(":8080")
}
