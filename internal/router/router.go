package router

import (
	"TestForWork/internal/handler"
	"github.com/gin-gonic/gin"
	"net/http"
)

func NewRouter(h *handler.SubscriptionHandler) *gin.Engine {
	router := gin.New()

	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	api := router.Group("/api/v1")
	{
		api.POST("/subscriptions", h.CreateSubscription)
		api.GET("/subscriptions", h.GetAllSubscriptions)
		api.GET("/subscriptions/:id", h.GetSubscriptionByID)
		api.DELETE("/subscriptions/:id", h.DeleteSubscription)
		api.GET("/subscriptions/cost", h.CalculateTotalCost)
	}

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"message": "Subscription service is running",
		})
	})

	return router
}
