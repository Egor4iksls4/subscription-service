package handler

import (
	"TestForWork/internal/entity"
	"TestForWork/internal/service"
	"github.com/gin-gonic/gin"
	"log/slog"
	"net/http"
	"strconv"
)

// @title Subscription Service API
// @version 1.0
// @description REST API for managing user subscriptions
// host localhost:8080
// @BasePath /api/v1

type SubscriptionHandler struct {
	service service.SubscriptionServiceInterface
}

func NewSubscriptionHandler(s service.SubscriptionServiceInterface) *SubscriptionHandler {
	slog.Info("Initializing Subscription Handler")
	return &SubscriptionHandler{service: s}
}

// CreateSubscription godoc
// @Summary Create a new subscription
// @Description Create a new subscription for a user
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param subscription body model.CreateSubscriptionRequest true "Subscription data"
// @Success 201 {object} entity.Subscription
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /subscriptions [post]
func (h *SubscriptionHandler) CreateSubscription(c *gin.Context) {
	slog.Info("Handling create subscription request")

	var req entity.CreateSubscriptionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		slog.Warn("Invalid request body for create subscription",
			"error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	sub, err := h.service.Create(c.Request.Context(), req)
	if err != nil {
		slog.Error("Failed to create subscription",
			"error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	slog.Info("Successfully created subscription",
		"subscription_id", sub.ID)
	c.JSON(http.StatusCreated, sub)
}

// GetAllSubscriptions godoc
// @Summary Get all subscriptions
// @Description Get a list of all subscriptions
// @Tags subscriptions
// @Produce json
// @Success 200 {array} entity.Subscription
// @Failure 500 {object} map[string[string]
// @Router /subscriptions [get]
func (h *SubscriptionHandler) GetAllSubscriptions(c *gin.Context) {
	slog.Info("Handling get all subscriptions request")

	subs, err := h.service.GetAll(c.Request.Context())
	if err != nil {
		slog.Error("Failed to fetch subscriptions",
			"error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	slog.Info("Successfully fetched all subscriptions",
		"count", len(subs))
	c.JSON(http.StatusOK, subs)
}

// GetSubscriptionByID godoc
// @Summary Get subscription by ID
// @Description Get a specific subscription by ID
// @Tags subscription
// @Produce json
// @Param id path int true "Subscription ID"
// @Success 200 {object} entity.Subscription
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /subscriptions/{id} [get]
func (h *SubscriptionHandler) GetSubscriptionByID(c *gin.Context) {
	slog.Info("Handling get subscription by ID request")

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		slog.Warn("Invalid subscription ID format",
			"id", c.Param("id"))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	sub, err := h.service.GetByID(c.Request.Context(), id)
	if err != nil {
		slog.Error("Subscription not found",
			"subscription_id", id,
			"error", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "subscription not found"})
		return
	}

	slog.Info("Successfully fetched subscription",
		"subscription_id", sub.ID)
	c.JSON(http.StatusOK, sub)
}

// DeleteSubscription godoc
// @Summary Delete subscription
// @Description Delete a subscription by its ID
// @Tags subscription
// @Param id path int true "Subscription ID"
// @Success 204
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /subscriptions/{id} [delete]
func (h *SubscriptionHandler) DeleteSubscription(c *gin.Context) {
	slog.Info("Handling delete subscription request")

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		slog.Warn("Invalid subscription ID format",
			"id", c.Param("id"))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	err = h.service.Delete(c.Request.Context(), id)
	if err != nil {
		slog.Error("Failed to delete subscription",
			"subscription_id", id,
			"error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	slog.Info("Successfully deleted subscription",
		"subscription_id", id)
	c.Status(http.StatusNoContent)
}

// CalculateTotalCost godoc
// @Summary Calculate total cost of subscriptions
// @Description Calculate total cost with optional filters by user_id and service_name
// @Tags subscription
// @Produce json
// @Param user_id query string false "User ID (UUID format)"
// @Param service_name query string false "Service name"
// @Param start_date query string true "Start date (MM-YYYY format)"
// @Param end_date query string true "End date (MM-YYYY format)"
// @Success 200 {object} entity.SubscriptionCostResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /subscriptions/cost [get]
func (h *SubscriptionHandler) CalculateTotalCost(c *gin.Context) {
	slog.Info("Handling calculate total cost request")

	var req entity.SubscriptionCostRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		slog.Warn("Invalid query parameters for calculate cost",
			"error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	total, err := h.service.CalculateTotalCost(c.Request.Context(), req)
	if err != nil {
		slog.Error("Failed to calculate total cost",
			"error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response := entity.SubscriptionCostResponse{
		TotalCost: total,
	}

	slog.Info("Successfully calculated total cost",
		"total_cost", total)
	c.JSON(http.StatusOK, response)
}
