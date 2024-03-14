package productcollector

import (
	"context"
	"encoding/json"
	"fmt"
	"goapi/internal/model"
	"log/slog"
	"net/http"
	"time"
)

const (
	apiURL         = "https://petstore.swagger.io/v2/pet/findByStatus?status=available"
	tickerInterval = 30 * time.Minute
)

type ProductCollector struct {
	ProductSaver
	log *slog.Logger
}

type ProductSaver interface {
	AddProducts(ctx context.Context, products []model.Product) error
}

func NewProductCollector(saver ProductSaver, log *slog.Logger) *ProductCollector {
	return &ProductCollector{
		ProductSaver: saver,
		log:          log,
	}
}

func (p *ProductCollector) Collect(ctx context.Context) {
	const op = "productcollector.Collect"

	log := p.log.With(
		slog.String("op", op),
	)

	log.Info("collect product")

	ticker := time.NewTicker(tickerInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			fmt.Println("Stopping product data collection...")
			return
		case <-ticker.C:
			log.Info("start collect product")

			response, err := http.Get(apiURL)
			if err != nil {
				log.Error("HTTP request failed: %v", err)
				continue
			}
			defer response.Body.Close()

			if response.StatusCode != http.StatusOK {
				log.Error("HTTP request failed: %s", response.Status)
				continue
			}

			var products []model.Product
			if err := json.NewDecoder(response.Body).Decode(&products); err != nil {
				log.Error("Failed to decode JSON: %v", err)
				continue
			}

			err = p.ProductSaver.AddProducts(ctx, products)
			if err != nil {
				log.Error("Failed to save product from api: %v", err)
				continue
			}
			log.Info("collect product successfully")
		default:
			continue
		}
	}
}
