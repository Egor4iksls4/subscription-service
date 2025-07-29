package repository

import (
	"TestForWork/internal/entity"
	"context"
	"database/sql"
	"errors"
	"github.com/google/uuid"
	"log/slog"
	"time"
)

type SubscriptionRepository interface {
	Create(ctx context.Context, sub *entity.Subscription) error
	GetAll(ctx context.Context) ([]entity.Subscription, error)
	GetByID(ctx context.Context, id int) (*entity.Subscription, error)
	Delete(ctx context.Context, id int) error
	CalculateTotalCost(ctx context.Context, userID *string, serviceName *string, startDate, endDate time.Time) (int, error)
}

type SubscriptionRepo struct {
	db *sql.DB
}

func NewSubscriptionRepository(db *sql.DB) SubscriptionRepository {
	slog.Info("Initializing Subscription Repository")
	return &SubscriptionRepo{db: db}
}

func (r *SubscriptionRepo) Create(ctx context.Context, sub *entity.Subscription) error {
	query := `INSERT INTO subscriptions (service_name, price, user_id, start_date, end_date)
			VALUES ($1, $2, $3, $4, $5) RETURNING id`

	slog.Debug("Create new subscription",
		"service_name", sub.ServiceName,
		"user_id", sub.UserID,
		"start_date", sub.StartDate)

	return r.db.QueryRowContext(ctx, query,
		sub.ServiceName,
		sub.Price,
		sub.UserID,
		sub.StartDate,
		sub.EndDate).Scan(&sub.ID)
}

func (r *SubscriptionRepo) GetAll(ctx context.Context) ([]entity.Subscription, error) {
	var subs []entity.Subscription
	query := `SELECT id, service_name, price, user_id, start_date, end_date FROM subscriptions ORDER BY id`

	slog.Debug("Fetching all subscriptions")

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		slog.Error("Failed to fetch subscriptions", "error", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var sub entity.Subscription
		err := rows.Scan(&sub.ID, &sub.ServiceName, &sub.Price, &sub.UserID, &sub.StartDate, &sub.EndDate)
		if err != nil {
			slog.Error("Failed to scan subscription", "error", err)
			return nil, err
		}
		subs = append(subs, sub)
	}

	if err = rows.Err(); err != nil {
		slog.Error("Error iterating rows", "error", err)
		return nil, err
	}

	slog.Info("Successfully fetched subscriptions", "count", len(subs))
	return subs, nil
}

func (r *SubscriptionRepo) GetByID(ctx context.Context, id int) (*entity.Subscription, error) {
	var sub entity.Subscription
	query := `SELECT id, service_name, price, user_id, start_date, end_date FROM subscriptions WHERE id = $1`

	slog.Debug("Fetching subscription", "id", id)

	err := r.db.QueryRowContext(ctx, query, id).Scan(&sub.ID, &sub.ServiceName, &sub.Price, &sub.UserID, &sub.StartDate, &sub.EndDate)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			slog.Warn("Subscription not found", "subscription_id", id)
			return nil, err
		}
		slog.Error("Failed to fetch subscription", "subscription_id", id, "error", err)
	}
	slog.Info("Successfully fetched subscription", "subscription_id", id)
	return &sub, nil
}

func (r *SubscriptionRepo) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM subscriptions WHERE id = $1`

	slog.Debug("Deleting subscription", "subscription_id", id)

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		slog.Error("Failed to delete subscription", "subscription_id", id, "error", err)
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		slog.Error("Failed to delete subscription", "subscription_id", id, "error", err)
		return err
	}

	if rowsAffected == 0 {
		slog.Warn("Subscription not found", "subscription_id", id)
		return sql.ErrNoRows
	}

	slog.Info("Successfully deleted subscription", "subscription_id", id)
	return nil
}

func (r *SubscriptionRepo) CalculateTotalCost(ctx context.Context,
	userID *string, serviceName *string, startDate, endDate time.Time) (int, error) {

	query := `SELECT COALESCE(SUM(price), 0) FROM subscriptions
			WHERE start_date <= $1 AND (end_date IS NULL OR end_date >= $2)`
	args := []interface{}{endDate, startDate}

	slog.Debug("Calculating total cost of subscriptions",
		"start_date", startDate,
		"end_date", endDate,
		"user_id", userID,
		"service_name", serviceName)

	if userID != nil {
		parsedUUID, err := uuid.Parse(*userID)
		if err != nil {
			slog.Error("Failed to parse user ID", "user_id", userID, "error", err)
			return 0, err
		}
		query += ` AND user_id = $3`
		args = append(args, parsedUUID)
	}

	if serviceName != nil {
		query += ` AND service_name = $4`
		args = append(args, *serviceName)
	}

	var total int
	err := r.db.QueryRowContext(ctx, query, args...).Scan(&total)
	if err != nil {
		slog.Error("Failed to calculate total of subscriptions", "error", err)
		return 0, err
	}

	slog.Info("Successfully calculated total of subscriptions", "total_cost", total)
	return total, nil
}
