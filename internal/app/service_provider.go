package app

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/RomanLevBy/BurgersAPI/internal/config"
	getAllBurger "github.com/RomanLevBy/BurgersAPI/internal/http-server/handlers/burger/all"
	getBurger "github.com/RomanLevBy/BurgersAPI/internal/http-server/handlers/burger/get"
	saveBurger "github.com/RomanLevBy/BurgersAPI/internal/http-server/handlers/burger/save"
	getCategory "github.com/RomanLevBy/BurgersAPI/internal/http-server/handlers/category/get"
	burgerRepo "github.com/RomanLevBy/BurgersAPI/internal/repository/burger/postgres"
	categoryRepo "github.com/RomanLevBy/BurgersAPI/internal/repository/category/postgres"
	burgerService "github.com/RomanLevBy/BurgersAPI/internal/service/burger"
	categoryService "github.com/RomanLevBy/BurgersAPI/internal/service/category"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	_ "github.com/lib/pq"
	"log"
	"log/slog"
	"net/http"
	"os"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

type serviceProvider struct {
	conf       *config.Config
	logger     *slog.Logger
	postgresDB *sql.DB
	server     *http.Server

	categoryService *categoryService.Service
	categoryRepo    *categoryRepo.Repository

	burgerService *burgerService.Service
	burgerRepo    *burgerRepo.Repository
}

func newServiceProvider(ctx context.Context, config *config.Config) *serviceProvider {
	provider := serviceProvider{}
	provider.conf = config

	provider.InitLogger()
	provider.InitCategoryService()
	provider.InitBurgerService()
	provider.InitServer()

	_ = ctx //todo delete context or use for init databases?

	return &provider
}

func (s *serviceProvider) InitLogger() *slog.Logger {
	if s.logger == nil {
		switch s.conf.Env {
		case envLocal:
			s.logger = slog.New(
				slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
			)
		case envDev:
			s.logger = slog.New(
				slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
			)
		case envProd:
			s.logger = slog.New(
				slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
			)
		}
	}

	s.logger.Debug("Debug logging enabled")

	return s.logger
}

func (s *serviceProvider) InitCategoryService() *categoryService.Service {
	if s.categoryService == nil {
		s.categoryService = categoryService.New(
			s.InitCategoryRepository(),
			s.InitLogger(),
		)
	}

	return s.categoryService
}

func (s *serviceProvider) InitCategoryRepository() *categoryRepo.Repository {
	if s.categoryRepo == nil {
		s.categoryRepo = categoryRepo.New(
			s.InitPostgres(),
		)
	}

	return s.categoryRepo
}

func (s *serviceProvider) InitBurgerService() *burgerService.Service {
	if s.burgerService == nil {
		s.burgerService = burgerService.New(
			s.InitBurgerRepository(),
			s.InitLogger(),
		)
	}

	return s.burgerService
}

func (s *serviceProvider) InitBurgerRepository() *burgerRepo.Repository {
	if s.burgerRepo == nil {
		s.burgerRepo = burgerRepo.New(
			s.InitPostgres(),
		)
	}

	return s.burgerRepo
}

func (s *serviceProvider) InitPostgres() *sql.DB {
	const fn = "init.postgres-db"

	if s.postgresDB == nil {
		postgresDB, err := sql.Open("postgres",
			fmt.Sprintf("host=%s user=%s dbname=%s sslmode=disable password=%s",
				s.conf.Postgres.Host,
				s.conf.Postgres.User,
				s.conf.Postgres.DBName,
				s.conf.Postgres.Password,
			),
		)

		if err != nil {
			log.Fatalf("%s, %v", fn, err)
		}

		s.postgresDB = postgresDB
	}

	return s.postgresDB
}

func (s *serviceProvider) InitServer() *http.Server {
	if s.server == nil {
		router := chi.NewRouter()
		router.Use(middleware.RequestID)
		router.Use(middleware.Recoverer)
		router.Use(middleware.URLFormat)

		//TODO: delete
		router.Get("/", func(w http.ResponseWriter, r *http.Request) {
			_, err := w.Write([]byte("welcome"))
			if err == nil {
				return
			}
		})

		router.Route("/v1", func(router chi.Router) {
			router.Get("/categories/{id}", getCategory.New(s.logger, s.categoryService))

			router.Get("/burgers", getAllBurger.New(s.logger, s.burgerService))
			router.Get("/burgers/{id}", getBurger.New(s.logger, s.burgerService))
			router.Post("/burgers", saveBurger.New(s.logger, s.burgerService))
		})

		s.logger.Info("starting server", slog.String("address", s.conf.Address))

		s.server = &http.Server{
			Addr:         s.conf.Address,
			Handler:      router,
			ReadTimeout:  s.conf.HTTPServer.Timeout,
			WriteTimeout: s.conf.HTTPServer.Timeout,
			IdleTimeout:  s.conf.HTTPServer.IdleTimeout,
		}

		s.logger.Info("server started", slog.String("address", s.conf.Address))
	}

	return s.server
}
