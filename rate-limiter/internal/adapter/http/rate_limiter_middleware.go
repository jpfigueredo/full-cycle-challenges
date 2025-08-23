package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jpfigueredo/rate-limiter-challenge/internal/usecase"
)

func RateLimiterMiddleware(uc *usecase.RateLimiterUseCase) gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		token := c.GetHeader("API_KEY")

		allowed, err := uc.CheckAndIncrement(c.Request.Context(), ip, token)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
			return
		}
		if !allowed {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"error": "you have reached the maximum number of requests or actions allowed within a certain time frame",
			})
			return
		}

		c.Next()
	}
}
