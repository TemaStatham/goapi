package product

import (
	"context"
	"errors"
	"fmt"
	"goapi/internal/model"
	"log/slog"
)

const (
	ErrProductId      = -1
	TagGetAllProducts = "get all products"
)

var (
	ErrNameIsEmpty      = errors.New("name is empty")
	ErrCategoryiesEmpty = errors.New("categoryies is empty")
	ErrUnknownTag       = errors.New("unknown tag get all products")
)

type Service struct {
	adder   AdderProduct
	deleter DeleterProduct
	updater UpdaterProduct
	getter  GetterProduct
	log     *slog.Logger
}

type AdderProduct interface {
	AddProduct(ctx context.Context, name string, categoryies []model.Category) (int64, error)
}

type DeleterProduct interface {
	DeleteProduct(ctx context.Context, id int64) error
}

type UpdaterProduct interface {
	UpdateProductName(ctx context.Context, id int64, name string) (int64, error)
	UpdateProductCategoryies(ctx context.Context, id int64, category []model.Category) (int64, error)
}

type GetterProduct interface {
	GetAllProducts(ctx context.Context) ([]model.Product, error)
}

func (s *Service) AddProduct(ctx context.Context, name string, categoryies []model.Category) (int64, error) {
	const op = "product.AddProduct"

	log := s.log.With(
		slog.String("op", op),
		slog.String("name", name),
	)

	log.Info("add product")

	if name == "" {
		log.Error("data is invalid: ", ErrNameIsEmpty)
		return ErrProductId, fmt.Errorf("%s %w", op, ErrNameIsEmpty)
	}

	if len(categoryies) == 0 {
		log.Error("data is invalid: ", ErrCategoryiesEmpty)
		return ErrProductId, fmt.Errorf("%s %w", op, ErrCategoryiesEmpty)
	}

	productID, err := s.adder.AddProduct(ctx, name, categoryies)
	if err != nil {
		log.Error("product dont saved", err)
		return ErrProductId, fmt.Errorf("%s %w", op, err)
	}

	log.Info("product is added")

	return productID, nil
}

func (s *Service) DeleteProduct(ctx context.Context, id int64) error {
	const op = "product.DeleteProduct"

	log := s.log.With(
		slog.String("op", op),
		slog.Int64("id", id),
	)

	log.Info("delete product")

	if id <= 0 {
		log.Error("data is invalid: ", ErrNameIsEmpty)
		return fmt.Errorf("%s %w", op, ErrNameIsEmpty)
	}

	err := s.deleter.DeleteProduct(ctx, id)
	if err != nil {
		log.Error("product didnt deleted", err)
		return fmt.Errorf("%s %w", op, err)
	}

	log.Info("product is deleted")

	return nil
}

func (s *Service) EditProductName(ctx context.Context, id int64, name string) (int64, error) {
	const op = "product.EditProductName"

	log := s.log.With(
		slog.String("op", op),
		slog.Int64("id", id),
	)

	log.Info("edit product name")

	if id <= 0 {
		log.Error("data is invalid: ", ErrNameIsEmpty)
		return ErrProductId, fmt.Errorf("%s %w", op, ErrNameIsEmpty)
	}

	if name == "" {
		log.Error("data is invalid: ", ErrNameIsEmpty)
		return ErrProductId, fmt.Errorf("%s %w", op, ErrNameIsEmpty)
	}

	productID, err := s.updater.UpdateProductName(ctx, id, name)
	if err != nil {
		log.Error("product name didnt edited", err)
		return ErrProductId, fmt.Errorf("%s %w", op, err)
	}

	log.Info("product is edited")

	return productID, nil
}

func (s *Service) EditProductCategoryies(ctx context.Context, id int64, categoryies []model.Category) (int64, error) {
	const op = "product.EditProductCategoryies"

	log := s.log.With(
		slog.String("op", op),
		slog.Int64("id", id),
	)

	log.Info("edit product categoryies")

	if id <= 0 {
		log.Error("data is invalid: ", ErrNameIsEmpty)
		return ErrProductId, fmt.Errorf("%s %w", op, ErrNameIsEmpty)
	}

	if len(categoryies) == 0 {
		log.Error("data is invalid: ", ErrCategoryiesEmpty)
		return ErrProductId, fmt.Errorf("%s %w", op, ErrCategoryiesEmpty)
	}

	productID, err := s.updater.UpdateProductCategoryies(ctx, id, categoryies)
	if err != nil {
		log.Error("product categoryies didnt edited", err)
		return ErrProductId, fmt.Errorf("%s %w", op, err)
	}

	log.Info("product is edited")

	return productID, nil
}

func (s *Service) GetAllProducts(ctx context.Context, tag string) ([]model.Product, error) {
	const op = "product.EditProduct"

	log := s.log.With(
		slog.String("op", op),
		slog.String("tag", tag),
	)

	log.Info("get all product")

	if tag != TagGetAllProducts {
		log.Error("products didnt get", ErrUnknownTag)
		return []model.Product{}, fmt.Errorf("%s %w", op, ErrUnknownTag)
	}

	products, err := s.getter.GetAllProducts(ctx)
	if err != nil {
		log.Error("products didnt get", err)
		return []model.Product{}, fmt.Errorf("%s %w", op, err)
	}

	log.Info("all products is getter")

	return products, nil
}
