package app

import (
	"context"
	"goapi/internal/app/logger"
	"goapi/internal/app/productcollector"
	"goapi/internal/app/server"
	"goapi/internal/config"
	"goapi/internal/handler"
	"goapi/internal/repository/postgres"
	"goapi/internal/service"
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

func (a *App) MustRun(cfg *config.Config) {
	if err := a.Run(cfg); err != nil {
		panic(err)
	}
}

func (a *App) Run(cfg *config.Config) error {
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
		return err
	}

	authRep := postgres.NewAuthPostgres(db, a.log)
	productRep := postgres.NewProductRepository(db, a.log)
	categoryRep := postgres.NewCategoryRepository(db, a.log)

	authServ := service.NewAuthService(authRep, authRep, a.log, cfg.TokenTTL)
	productServ := service.NewProductService(productRep, productRep, productRep, productRep, a.log)
	categoryServ := service.NewCategoryService(categoryRep, categoryRep, categoryRep, categoryRep, a.log)

	handlers := handler.NewHandler(authServ, productServ, categoryServ, a.log)

	srv := new(server.Server)
	go func() {
		if err := srv.Run(cfg.SConfig.Port, handlers.Init()); err != nil {
			log.Error("error occured while running http server: %s", err.Error())
		}
	}()

	ctx, cancel := context.WithCancel(context.Background())
	collector := productcollector.NewProductCollector(productServ, a.log)
	go collector.Collect(ctx)

	log.Info("Application started")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	log.Info("Application Shutting Down")

	cancel()

	if err := srv.Shutdown(context.Background()); err != nil {
		log.Info("error occured on server shutting down: %s", err.Error())
		return err
	}

	if err := db.Close(); err != nil {
		log.Info("error occured on db connection close: %s", err.Error())
		return err
	}

	return nil
}
