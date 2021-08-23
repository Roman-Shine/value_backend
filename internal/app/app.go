package app

import (
	"context"
	"errors"
	"github.com/Roman-Shine/value_backend/internal/config"
	delivery "github.com/Roman-Shine/value_backend/internal/delivery/http"
	"github.com/Roman-Shine/value_backend/internal/repository"
	"github.com/Roman-Shine/value_backend/internal/server"
	"github.com/Roman-Shine/value_backend/internal/service"
	"github.com/Roman-Shine/value_backend/pkg/auth"
	"github.com/Roman-Shine/value_backend/pkg/cache"
	"github.com/Roman-Shine/value_backend/pkg/database/mongodb"
	"github.com/Roman-Shine/value_backend/pkg/hash"
	"github.com/Roman-Shine/value_backend/pkg/logger"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func Run(configPath string) {
	cfg, err := config.Init(configPath)
	if err != nil {
		logger.Error(err)
		return
	}

	// Dependencies
	mongoClient, err := mongodb.NewClient(cfg.Mongo.URI, cfg.Mongo.User, cfg.Mongo.Password)
	if err != nil {
		logger.Error(err)

		return
	}

	db := mongoClient.Database(cfg.Mongo.Name)

	memCache := cache.NewMemoryCache()
	hasher := hash.NewSHA1Hasher(cfg.Auth.PasswordSalt)

	tokenManager, err := auth.NewManager(cfg.Auth.JWT.SigningKey)
	if err != nil {
		logger.Error(err)

		return
	}

	// Services, Repos & API Handlers
	repos := repository.NewRepositories(db)
	services := service.NewServices(service.Deps{
		Repos:                  repos,
		Cache:                  memCache,
		Hasher:                 hasher,
		AccessTokenTTL:         cfg.Auth.JWT.AccessTokenTTL,
		RefreshTokenTTL:        cfg.Auth.JWT.RefreshTokenTTL,
		CacheTTL:               int64(cfg.CacheTTL.Seconds()),
		VerificationCodeLength: cfg.Auth.VerificationCodeLength,
		Environment:            cfg.Environment,
		Domain:                 cfg.HTTP.Host,
	})
	handlers := delivery.NewHandler(services, tokenManager)

	// HTTP Server
	srv := server.NewServer(cfg, handlers.Init(cfg))

	go func() {
		if err := srv.Run(); !errors.Is(err, http.ErrServerClosed) {
			logger.Errorf("error occurred while running http server: %s\n", err.Error())
		}
	}()

	logger.Info("Server started")

	// Graceful Shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	<-quit

	const timeout = 5 * time.Second

	ctx, shutdown := context.WithTimeout(context.Background(), timeout)
	defer shutdown()

	if err := srv.Stop(ctx); err != nil {
		logger.Errorf("failed to stop server: %v", err)
	}

	if err := mongoClient.Disconnect(context.Background()); err != nil {
		logger.Error(err.Error())
	}
}
