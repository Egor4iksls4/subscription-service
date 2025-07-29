package service

import (
	"TestForWork/internal/entity"
	"TestForWork/internal/repository"
	"context"
	"errors"
	"log/slog"
	"time"
)

type SubscriptionServiceInterface interface {
	Create(ctx context.Context, req entity.CreateSubscriptionRequest) (*entity.Subscription, error)
	GetAll(ctx context.Context) ([]entity.Subscription, error)
	GetByID(ctx context.Context, id int) (*entity.Subscription, error)
	Delete(ctx context.Context, id int) error
	CalculateTotalCost(ctx context.Context, req entity.SubscriptionCostRequest) (int, error)
}

type SubService struct {
	repo repository.SubscriptionRepository
}

func NewSubscriptionService(repo repository.SubscriptionRepository) SubscriptionServiceInterface {
	slog.Info("New Subscription Service")
	return &SubService{repo: repo}
}

func (s *SubService) Create(ctx context.Context, req entity.CreateSubscriptionRequest) (*entity.Subscription, error) {
	slog.Info("Creating new subscription",
		"service_name", req.ServiceName,
		"user_id", req.UserID)

	startDate, err := time.Parse("01-2006", req.StartDate)
	if err != nil {
		slog.Error("Invalid start date format",
			"start_date", req.StartDate,
			"error", err)
		return nil, errors.New("invalid start_date format, expected MM-YYYY")
	}

	var endDate *time.Time
	if req.EndDate != nil {
		ed, err := time.Parse("01-2006", *req.EndDate)
		if err != nil {
			slog.Error("Invalid end date format",
				"end_date", *req.EndDate,
				"error", err)
			return nil, errors.New("invalid end_date format, expected MM-YYYY")
		}
		endDate = &ed
	}

	sub := &entity.Subscription{
		ServiceName: req.ServiceName,
		Price:       req.Price,
		UserID:      req.UserID,
		StartDate:   startDate,
		EndDate:     endDate,
	}

	err = s.repo.Create(ctx, sub)
	if err != nil {
		slog.Error("Failed to create subscription in repository",
			"error", err)
		return nil, err
	}

	slog.Info("Successfully created subscription",
		"subscription_id", sub.ID)
	return sub, nil
}

func (s *SubService) GetAll(ctx context.Context) ([]entity.Subscription, error) {
	slog.Info("Fetching all subscriptions")
	return s.repo.GetAll(ctx)
}

func (s *SubService) GetByID(ctx context.Context, id int) (*entity.Subscription, error) {
	slog.Info("Fetching subscription by ID",
		"subscription_id", id)
	return s.repo.GetByID(ctx, id)
}

func (s *SubService) Delete(ctx context.Context, id int) error {
	slog.Info("Deleting subscription",
		"subscription_id", id)
	return s.repo.Delete(ctx, id)
}

func (s *SubService) CalculateTotalCost(ctx context.Context, req entity.SubscriptionCostRequest) (int, error) {
	slog.Info("Calculating total cost of subscription",
		"start_date", req.StartDate,
		"end_date", req.EndDate,
		"user_id", req.UserID,
		"service_name", req.ServiceName)

	startDate, err := time.Parse("01-2006", req.StartDate)
	if err != nil {
		slog.Error("Invalid start date format",
			"start_date", req.StartDate,
			"error", err)
		return 0, errors.New("invalid start_date format, expected MM-YYYY")
	}

	endDate, err := time.Parse("01-2006", req.EndDate)
	if err != nil {
		slog.Error("Invalid end date format",
			"end_date", req.EndDate,
			"error", err)
		return 0, errors.New("invalid end_date format, expected MM-YYYY")
	}

	total, err := s.repo.CalculateTotalCost(ctx, req.UserID, req.ServiceName, startDate, endDate)
	if err != nil {
		slog.Error("Failed to calculate total cost",
			"error", err)
		return 0, err
	}

	slog.Info("Successfully calculated total cost",
		"total_cost", total)
	return total, nil
}
