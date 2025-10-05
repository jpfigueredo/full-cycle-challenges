package main

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/jpfigueredo/cep-clima-distributed/internal/shared/repository"
	"github.com/jpfigueredo/cep-clima-distributed/internal/shared/usecase"
)

func main() {
	_ = godotenv.Load()
	cepRepo := repository.NewViaCEPRepo()
	climaRepo := repository.NewWeatherAPIRepo()
	uc := usecase.NewCEPClimaUseCase(cepRepo, climaRepo)

	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})

	r.GET("/clima/:cep", func(c *gin.Context) {
		cep := c.Param("cep")
		out, err := uc.CheckCEPAndGetClima(c.Request.Context(), cep)
		if err != nil {
			switch err {
			case usecase.ErrInvalidCEP:
				c.JSON(http.StatusUnprocessableEntity, gin.H{"message": err.Error()})
			case usecase.ErrCEPNotFound:
				c.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
			default:
				c.JSON(http.StatusInternalServerError, gin.H{"message": "internal error"})
			}
			return
		}
		c.JSON(http.StatusOK, out)
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	r.Run(":" + port)
}
