package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"unicode"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/jpfigueredo/cep-clima-distributed/internal/shared/tracing"
)

type CEPInput struct {
	CEP string `json:"cep" binding:"required"`
}

func main() {
	_ = godotenv.Load()
	defer tracing.InitTracer("service-a")() // Ou "service-b"
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong from A"})
	})

	r.POST("/clima", func(c *gin.Context) {
		var input CEPInput
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "invalid input"})
			return
		}

		// Validação simples: string com exatos 8 dígitos
		if len(input.CEP) != 8 || !isDigitsOnly(input.CEP) {
			c.JSON(http.StatusUnprocessableEntity, gin.H{"message": "invalid zipcode"})
			return
		}

		serviceBURL := os.Getenv("SERVICE_B_URL")
		if serviceBURL == "" {
			serviceBURL = "http://localhost:8081"
		}
		url := fmt.Sprintf("%s/clima/%s", serviceBURL, input.CEP)

		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "internal error"})
			return
		}

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "internal error"})
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			var errMsg map[string]string
			json.NewDecoder(resp.Body).Decode(&errMsg)
			c.JSON(resp.StatusCode, errMsg)
			return
		}

		var output map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&output)
		c.JSON(http.StatusOK, output)
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	r.Run(":" + port)
}

func isDigitsOnly(s string) bool {
	for _, r := range s {
		if !unicode.IsDigit(r) {
			return false
		}
	}
	return true
}
