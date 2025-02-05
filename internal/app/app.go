package app

import (
	"context"
	"fmt"

	"github.com/aspirin100/finapi/internal/config"
	"github.com/aspirin100/finapi/internal/handler"
	"github.com/aspirin100/finapi/internal/repository"
	"github.com/aspirin100/finapi/internal/service"
)

type App struct {
	requestHandler *handler.Handler
	repo           *repository.Repository
}

func New(ctx context.Context, cfg *config.Config) (*App, error) {
	repo, err := repository.NewConnection(ctx, cfg.PostgresDSN)
	if err != nil {
		return nil, fmt.Errorf("failed to create app instance: %w", err)
	}

	srvc := service.New(cfg.Timeout, repo)

	requestHandler := handler.New(cfg.Hostname, cfg.Port, srvc)

	return &App{
		requestHandler: requestHandler,
		repo:           repo,
	}, nil
}

func (app *App) Run() error {
	err := app.requestHandler.Run()
	if err != nil {
		return fmt.Errorf("failed to start application: %w", err)
	}

	return nil
}

func (app *App) Stop(ctx context.Context) error {
	app.repo.DB.Close()

	err := app.requestHandler.Shutdown(ctx)
	if err != nil {
		return fmt.Errorf("failed to shutdown http server: %w", err)
	}

	return nil
}
