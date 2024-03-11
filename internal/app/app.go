package app

import (
	"context"
	"goapi/internal/app/logger"
	"goapi/internal/app/server"
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

func (a *App) Run() {
	const op = "app.run"

	log := a.log.With(
		slog.String("op", op),
	)

	//db, err := repository.NewPostgresDB(repository.Config{
	//	Host:     viper.GetString("db.host"),
	//	Port:     viper.GetString("db.port"),
	//	Username: viper.GetString("db.username"),
	//	DBName:   viper.GetString("db.dbname"),
	//	SSLMode:  viper.GetString("db.sslmode"),
	//	Password: os.Getenv("DB_PASSWORD"),
	//})
	//if err != nil {
	//	log.Error("failed to initialize db: %s", err)
	//}
	//
	//repos := repository.NewRepository(db)
	//services := service.NewService(repos)
	//handlers := handler.NewHandler(services)

	srv := new(server.Server)
	//go func() {
	//	if err := srv.Run(viper.GetString("port"), handlers.InitRoutes()); err != nil {
	//		log.Error("error occured while running http server: %s", err.Error())
	//	}
	//}()

	log.Info("Application started")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	log.Info("Application Shutting Down")

	if err := srv.Shutdown(context.Background()); err != nil {
		log.Info("error occured on server shutting down: %s", err.Error())
	}

	//if err := db.Close(); err != nil {
	//	log.Info("error occured on db connection close: %s", err.Error())
	//}
}
