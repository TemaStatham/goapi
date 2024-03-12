package main

import (
	"goapi/internal/app"
	"goapi/internal/config"
)

// Запуск бд
// docker run --name=api-db -e POSTGRES_PASSWORD='qwerty' -p 5436:5432 -d --rm postgres
// migrate -path ./schema -database 'postgres://postgres:qwerty@localhost:5436/postgres?sslmode=disable' up

// Запуск приложения
// go run ./cmd/api/main.go --config="./config/config.yaml"

func main() {
	cfg := config.MustLoad()

	a := app.New(cfg.Env)

	a.MustRun(cfg)
}
