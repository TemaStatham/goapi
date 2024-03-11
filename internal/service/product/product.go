package product

import (
	"context"
	"goapi/internal/model"
	"log/slog"
)

type Service struct {
	adder   AdderProduct
	deleter DeleterProduct
	updater UpdaterProduct
	getter  GetterProduct
	log     *slog.Logger
}

type AdderProduct interface {
	AddProduct(product model.Product) error
}

type DeleterProduct interface {
	DeleteProduct() error
}

type UpdaterProduct interface {
	UpdateProduct(product model.Product) error
}

type GetterProduct interface {
	GetAllProducts() ([]model.Product, error)
}

func (s *Service) AddProduct(ctx context.Context, name string, categoryies []string) (productID int64, err error) {
}

func (s *Service) DeleteProduct(ctx context.Context, name string) error {

}
func (s *Service) EditProduct(ctx context.Context, name string) (productID int64, err error) {

}
func (s *Service) GetAllProducts(ctx context.Context, tag string) (product []model.Product, err error) {
}
