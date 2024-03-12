package app

import (
	"context"
	"goapi/internal/app/logger"
	"goapi/internal/app/server"
	"goapi/internal/config"
	"goapi/internal/handler"
	"goapi/internal/repository/postgres"
	"goapi/internal/service/auth"
	"goapi/internal/service/category"
	"goapi/internal/service/product"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

type App struct {
	log *slog.Logger
}

func New(env string) *App {
	return &App{
		log: logger.SetupLogger(env),
	}
}

func (a *App) Run(cfg *config.Config) {
	const op = "app.run"

	log := a.log.With(
		slog.String("op", op),
	)

	db, err := postgres.NewPostgresDB(postgres.Config{
		Host:     cfg.DBConfig.Host,
		Port:     cfg.DBConfig.Port,
		Username: cfg.DBConfig.Username,
		Password: cfg.DBConfig.Password,
		DBName:   cfg.DBConfig.DBName,
		SSLMode:  cfg.DBConfig.SSLMode,
	})
	if err != nil {
		log.Error("failed to initialize db: %s", err)
	}

	authRep := postgres.NewAuthPostgres(db, a.log)
	productRep := postgres.NewProductRepository(db, a.log)
	categoryRep := postgres.NewCategoryRepository(db, a.log)

	authServ := auth.NewService(authRep, authRep, a.log, cfg.TokenTTL)
	productServ := product.NewProductService(productRep, productRep, productRep, productRep, a.log)
	categoryServ := category.NewCategoryService(categoryRep, categoryRep, categoryRep, categoryRep, a.log)

	handlers := handler.NewHandler(authServ, productServ, categoryServ, a.log)

	srv := new(server.Server)
	go func() {
		if err := srv.Run(cfg.SConfig.Port, handlers.Init()); err != nil {
			log.Error("error occured while running http server: %s", err.Error())
		}
	}()

	log.Info("Application started")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	log.Info("Application Shutting Down")

	if err := srv.Shutdown(context.Background()); err != nil {
		log.Info("error occured on server shutting down: %s", err.Error())
	}

	if err := db.Close(); err != nil {
		log.Info("error occured on db connection close: %s", err.Error())
	}
}
