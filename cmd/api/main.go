package main

import (
	"goapi/internal/app"
	"goapi/internal/config"
)

// go run ./cmd/api/main.go --config="./config/config.yaml"

func main() {
	cfg := config.MustLoad()

	a := app.New(cfg.Env)

	a.Run()
}
