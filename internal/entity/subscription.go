package entity

import (
	"github.com/google/uuid"
	"time"
)

type Subscription struct {
	ID          int        `db:"id" json:"id"`
	ServiceName string     `db:"service_name" json:"service_name" binding:"required"`
	Price       int        `db:"price" json:"price" binding:"required"`
	UserID      uuid.UUID  `db:"user_id" json:"user_id" binding:"required"`
	StartDate   time.Time  `db:"start_date" json:"start_date" binding:"required"`
	EndDate     *time.Time `db:"end_date" json:"end_date,omitempty"`
}

type CreateSubscriptionRequest struct {
	ServiceName string    `json:"service_name" binding:"required"`
	Price       int       `json:"price" binding:"required,min=1"`
	UserID      uuid.UUID `json:"user_id" binding:"required"`
	StartDate   string    `json:"start_date" binding:"required"`
	EndDate     *string   `json:"end_date,omitempty"`
}

type SubscriptionCostRequest struct {
	UserID      *string `form:"user_id"`
	ServiceName *string `form:"service_name"`
	StartDate   string  `form:"start_date" binding:"required"`
	EndDate     string  `form:"end_date" binding:"required"`
}

type SubscriptionCostResponse struct {
	TotalCost int `json:"total_cost"`
}
