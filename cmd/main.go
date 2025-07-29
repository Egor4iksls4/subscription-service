package main

import (
	"database/sql"
	"github.com/spf13/viper"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"TestForWork/internal/handler"
	"TestForWork/internal/repository"
	"TestForWork/internal/router"
	"TestForWork/internal/service"
	_ "github.com/lib/pq"
)

// @title Subscription Service API
// @version 1.0
// @description REST API for managing user subscriptions
// @host localhost:8080
// @BasePath /api/v1

func main() {
	setupLogging()

	cfg := loadConfig()

	slog.Info("Starting Subscription Service")

	db, err := initDB(cfg)
	if err != nil {
		slog.Error("Failed to connect to database", "error", err)
	}
	defer closeDB(db)

	repo := repository.NewSubscriptionRepository(db)
	svc := service.NewSubscriptionService(repo)
	h := handler.NewSubscriptionHandler(svc)

	r := router.NewRouter(h)

	serverAddr := ":" + cfg.GetString("server.port")
	slog.Info("Server starting", "port", serverAddr)

	go func() {
		if err := r.Run(serverAddr); err != nil {
			slog.Error("Failed to start server", "error", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	slog.Info("Shutting down server...")
	slog.Info("Server exited")
}

func setupLogging() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	slog.SetDefault(logger)
}

func loadConfig() *viper.Viper {
	v := viper.New()

	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath("./config")

	if err := v.ReadInConfig(); err != nil {
		slog.Error("Cannot read config file: ", err)
	}

	slog.Info("Config file loaded", "file", v.ConfigFileUsed())
	return v
}

func initDB(cfg *viper.Viper) (*sql.DB, error) {
	connStr := formatConnectionString(cfg)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		db.Close()
		return nil, err
	}

	slog.Info("Connected to database successfully", connStr)
	return db, nil
}

func formatConnectionString(cfg *viper.Viper) string {
	return "host=" + cfg.GetString("database.host") +
		" port=" + cfg.GetString("database.port") +
		" user=" + cfg.GetString("database.user") +
		" password=" + cfg.GetString("database.password") +
		" dbname=" + cfg.GetString("database.dbname") +
		" sslmode=disable"
}

func closeDB(db *sql.DB) {
	slog.Info("Closing database connection...")
	if err := db.Close(); err != nil {
		slog.Error("Error closing database connection", "error", err)
	} else {
		slog.Info("Database connection closed")
	}
}
