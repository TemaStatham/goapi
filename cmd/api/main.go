package main

import (
	"fmt"
	"goapi/internal/app"
	"goapi/internal/config"
)

// go run ./cmd/api/main.go --config="./config/config.yaml"

func main() {
	cfg := config.MustLoad()
	fmt.Println(cfg)
	a := app.New(cfg.Env)

	a.Run(cfg)
}
