package main

import (
	"context"
	"embed"
	"fmt"
	"log"
	"os"
	"os/signal"

	"go.uber.org/zap"

	"github.com/alesr/urltinyizer/app"
	"github.com/alesr/urltinyizer/internal/repository"
	"github.com/alesr/urltinyizer/internal/service"
	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"
	envars "github.com/netflix/go-env"
	"github.com/pressly/goose/v3"
)

//go:embed migrations/*.sql
var embedMigrations embed.FS

const (
	postgresDriverName string = "postgres"
	dbMigrationsDir    string = "migrations"
)

type config struct {
	AppHost string `env:"APP_HOST,default=http://localhost:8080/"`
	DBUser  string `env:"POSTGRES_USER,default=user"`
	DBPass  string `env:"POSTGRES_PASSWORD,default=password"`
	DBName  string `env:"POSTGRES_DB,default=urltinyizer"`
	DBHost  string `env:"POSTGRES_HOST,default=db"`
	DBPort  string `env:"POSTGRES_PORT,default=5432"`
}

func newConfig() *config {
	var cfg config
	if _, err := envars.UnmarshalFromEnviron(&cfg); err != nil {
		log.Fatal(err)
	}
	return &cfg
}

func main() {
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalln("failed to initialize zap logger:", err)
	}
	defer logger.Sync()

	cfg := newConfig()

	db, err := sqlx.Open(postgresDriverName, fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPass, cfg.DBName),
	)
	if err != nil {
		logger.Fatal("failed to connect to database", zap.Error(err))
	}

	defer db.Close()

	goose.SetBaseFS(embedMigrations)

	if err := goose.SetDialect(postgresDriverName); err != nil {
		logger.Fatal("failed to set goose dialect", zap.Error(err))
	}

	if err := goose.Up(db.DB, dbMigrationsDir); err != nil {
		logger.Fatal("failed to run goose migrations", zap.Error(err))
	}

	repo := repository.NewPostgreSQL(logger, db)
	service := service.NewServiceDefault(logger, cfg.AppHost, repo)
	router := chi.NewRouter()
	app := app.NewREST(logger, router, service)

	app.RegisterRoutes()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	defer func() {
		signal.Stop(c)
		cancel()
	}()
	go func() {
		select {
		case <-c:
			cancel()
		case <-ctx.Done():
		}
	}()

	if err := app.Run(ctx); err != nil {
		logger.Fatal("failed to run app", zap.Error(err))
	}
}
