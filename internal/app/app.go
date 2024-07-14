package app

import (
	"context"
	"github.com/RomanLevBy/BurgersAPI/internal/config"
)

type App struct {
	conf            *config.Config
	serviceProvider *serviceProvider
}

func NewApp(ctx context.Context) (*App, error) {
	a := &App{}

	err := a.initDeps(ctx)
	if err != nil {
		return nil, err
	}

	return a, nil
}

func (a *App) Run(ctx context.Context) error {
	_ = ctx
	if err := a.serviceProvider.server.ListenAndServe(); err != nil {
		return err
	}

	return nil
}

func (a *App) initDeps(ctx context.Context) error {
	inits := []func(context.Context) error{
		a.initConfig,
		a.initServiceProvider,
	}

	for _, f := range inits {
		err := f(ctx)
		if err != nil {
			return err
		}
	}

	return nil
}

func (a *App) initConfig(_ context.Context) error {
	conf, err := config.Load()
	if err != nil {
		return err
	}
	a.conf = conf

	return nil
}

func (a *App) initServiceProvider(ctx context.Context) error {
	a.serviceProvider = newServiceProvider(ctx, a.conf)

	return nil
}
