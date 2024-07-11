package app

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/RomanLevBy/BurgersAPI/internal/config"
	"github.com/RomanLevBy/BurgersAPI/internal/http-server/handlers/category/get"
	categoryRepo "github.com/RomanLevBy/BurgersAPI/internal/repository/category/postgres"
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
}

func newServiceProvider(ctx context.Context, config *config.Config) *serviceProvider {
	provider := serviceProvider{}
	provider.conf = config

	provider.InitLogger()
	provider.InitCategoryService()
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

func (s *serviceProvider) InitPostgres() *sql.DB {
	const fn = "init.postgres-db"

	fmt.Println("config config config")

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
			w.Write([]byte("welcome"))
		})

		router.Get("/categories/{id}", get.New(s.logger, s.categoryService))

		s.logger.Info("starting server", slog.String("address", s.conf.Address))

		s.server = &http.Server{
			Addr:         s.conf.Address,
			Handler:      router,
			ReadTimeout:  s.conf.HTTPServer.Timeout,
			WriteTimeout: s.conf.HTTPServer.Timeout,
			IdleTimeout:  s.conf.HTTPServer.IdleTimeout,
		}
	}

	return s.server
}
