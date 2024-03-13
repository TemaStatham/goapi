package service

import (
	"context"
	"errors"
	"fmt"
	"goapi/internal/model"
	"goapi/internal/repository"
	"log/slog"
)

const (
	ErrProductId      = -1
	TagGetAllProducts = "get all products"
)

var (
	ErrProductNameIsEmpty = errors.New("product name is empty")
	ErrProductIDIsEmpty   = errors.New("product id is empty")
	ErrCategoryiesEmpty   = errors.New("product categoryies is empty")
	ErrProductsEmpty      = errors.New("products is empty")
	ErrProductUnknownTag  = errors.New("unknown tag get all products")
)

type ProductService struct {
	adder   AdderProduct
	deleter DeleterProduct
	updater UpdaterProduct
	getter  GetterProduct
	log     *slog.Logger
}

type AdderProduct interface {
	AddProduct(ctx context.Context, name string, categoryies []string) (int64, error)
	AddProducts(ctx context.Context, products []model.Product) error
}

type DeleterProduct interface {
	DeleteProduct(ctx context.Context, id int64) error
}

type UpdaterProduct interface {
	UpdateProductName(ctx context.Context, id int64, name string) (int64, error)
	UpdateProductCategoryies(ctx context.Context, id int64, categoryies []model.Category) (int64, error)
}

type GetterProduct interface {
	GetAllProducts(ctx context.Context) ([]model.Product, error)
	GetCategoryProducts(ctx context.Context, category string) ([]model.Product, error)
}

func NewProductService(
	a AdderProduct,
	d DeleterProduct,
	u UpdaterProduct,
	g GetterProduct,
	l *slog.Logger,
) *ProductService {
	return &ProductService{
		adder:   a,
		deleter: d,
		updater: u,
		getter:  g,
		log:     l,
	}
}

func (s *ProductService) AddProduct(ctx context.Context, name string, categoryies []string) (int64, error) {
	const op = "product.AddProduct"

	log := s.log.With(
		slog.String("op", op),
		slog.String("name", name),
	)

	log.Info("add product")

	if name == "" {
		log.Error("data is invalid: ", ErrProductNameIsEmpty)
		return ErrProductId, fmt.Errorf("%s %w", op, ErrProductNameIsEmpty)
	}

	if len(categoryies) == 0 {
		log.Error("data is invalid: ", ErrCategoryiesEmpty)
		return ErrProductId, fmt.Errorf("%s %w", op, ErrCategoryiesEmpty)
	}

	productID, err := s.adder.AddProduct(ctx, name, categoryies)
	if err != nil {
		if errors.Is(err, repository.ErrSaveProduct) {
			s.log.Warn("product isnt saved", err)
			return ErrProductId, fmt.Errorf("%s: %w", op, ErrorInvalidCredentials)
		}

		log.Error("product dont saved", err)
		return ErrProductId, fmt.Errorf("%s %w", op, err)
	}

	log.Info("product is added")

	return productID, nil
}

func (s *ProductService) DeleteProduct(ctx context.Context, id int64) error {
	const op = "product.DeleteProduct"

	log := s.log.With(
		slog.String("op", op),
		slog.Int64("id", id),
	)

	log.Info("delete product")

	if id <= 0 {
		log.Error("data is invalid: ", ErrProductIDIsEmpty)
		return fmt.Errorf("%s %w", op, ErrProductIDIsEmpty)
	}

	err := s.deleter.DeleteProduct(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrDeleteProduct) {
			s.log.Warn("product isnt deleted", err)
			return fmt.Errorf("%s: %w", op, ErrorInvalidCredentials)
		}
		if errors.Is(err, repository.ErrDeleteProductCategory) {
			s.log.Warn("product isnt deleted", err)
			return fmt.Errorf("%s: %w", op, ErrorInvalidCredentials)
		}

		log.Error("product didnt deleted", err)
		return fmt.Errorf("%s %w", op, err)
	}

	log.Info("product is deleted")

	return nil
}

func (s *ProductService) EditProductName(ctx context.Context, id int64, name string) (int64, error) {
	const op = "product.EditProductName"

	log := s.log.With(
		slog.String("op", op),
		slog.Int64("id", id),
	)

	log.Info("edit product name")

	if id <= 0 {
		log.Error("data is invalid: ", ErrProductIDIsEmpty)
		return ErrProductId, fmt.Errorf("%s %w", op, ErrProductIDIsEmpty)
	}

	if name == "" {
		log.Error("data is invalid: ", ErrProductNameIsEmpty)
		return ErrProductId, fmt.Errorf("%s %w", op, ErrProductNameIsEmpty)
	}

	productID, err := s.updater.UpdateProductName(ctx, id, name)
	if err != nil {
		if errors.Is(err, repository.ErrUpdateProduct) {
			s.log.Warn("product isnt edited", err)
			return ErrProductId, fmt.Errorf("%s: %w", op, ErrorInvalidCredentials)
		}

		log.Error("product name didnt edited", err)
		return ErrProductId, fmt.Errorf("%s %w", op, err)
	}

	log.Info("product is edited")

	return productID, nil
}

func (s *ProductService) EditProductCategory(ctx context.Context, id int64, categoryies []model.Category) (int64, error) {
	const op = "product.EditProductCategoryies"

	log := s.log.With(
		slog.String("op", op),
		slog.Int64("id", id),
	)

	log.Info("edit product categoryies")

	if id <= 0 {
		log.Error("data is invalid: ", ErrProductIDIsEmpty)
		return ErrProductId, fmt.Errorf("%s %w", op, ErrProductIDIsEmpty)
	}

	if len(categoryies) == 0 {
		log.Error("data is invalid: ", ErrCategoryiesEmpty)
		return ErrProductId, fmt.Errorf("%s %w", op, ErrCategoryiesEmpty)
	}

	productID, err := s.updater.UpdateProductCategoryies(ctx, id, categoryies)
	if err != nil {
		if errors.Is(err, repository.ErrDeleteProductCategory) {
			s.log.Warn("product isnt edited", err)
			return ErrProductId, fmt.Errorf("%s: %w", op, ErrorInvalidCredentials)
		}
		if errors.Is(err, repository.ErrSaveProductCategory) {
			s.log.Warn("product isnt edited", err)
			return ErrProductId, fmt.Errorf("%s: %w", op, ErrorInvalidCredentials)
		}

		log.Error("product categoryies didnt edited", err)
		return ErrProductId, fmt.Errorf("%s %w", op, err)
	}

	log.Info("product is edited")

	return productID, nil
}

func (s *ProductService) GetAllProducts(ctx context.Context, tag string) ([]model.Product, error) {
	const op = "product.GetAllProducts"

	log := s.log.With(
		slog.String("op", op),
		slog.String("tag", tag),
	)

	log.Info("get all product")

	if tag != TagGetAllProducts {
		log.Error("products didnt get", ErrProductUnknownTag)
		return []model.Product{}, fmt.Errorf("%s %w", op, ErrProductUnknownTag)
	}

	products, err := s.getter.GetAllProducts(ctx)
	if err != nil {
		if errors.Is(err, repository.ErrProductNotFound) {
			s.log.Warn("products empty", err)
			return []model.Product{}, fmt.Errorf("%s: %w", op, ErrorInvalidCredentials)
		}

		log.Error("products didnt get", err)
		return []model.Product{}, fmt.Errorf("%s %w", op, err)
	}

	log.Info("all products is getter")

	return products, nil
}

func (s *ProductService) GetCategoryProducts(ctx context.Context, category string) ([]model.Product, error) {
	const op = "product.GetAllProducts"

	log := s.log.With(
		slog.String("op", op),
		slog.String("category", category),
	)

	log.Info("get category product")

	if category == "" {
		log.Error("data is invalid: ", ErrCategoryiesEmpty)
		return []model.Product{}, fmt.Errorf("%s %w", op, ErrCategoryiesEmpty)
	}

	products, err := s.getter.GetCategoryProducts(ctx, category)
	if err != nil {
		if errors.Is(err, repository.ErrProductNotFound) {
			s.log.Warn("products empty", err)
			return []model.Product{}, fmt.Errorf("%s: %w", op, ErrorInvalidCredentials)
		}

		log.Error("products didnt get", err)
		return []model.Product{}, fmt.Errorf("%s %w", op, err)
	}

	log.Info("products is getter")

	return products, nil
}

func (s *ProductService) AddProducts(ctx context.Context, products []model.Product) error {
	const op = "postgres.AddProducts"

	log := s.log.With(
		slog.String("op", op),
	)

	log.Info("add products")

	if len(products) == 0 {
		log.Error("data is invalid: ", ErrProductsEmpty)
		return fmt.Errorf("%s %w", op, ErrProductsEmpty)
	}

	err := s.adder.AddProducts(ctx, products)
	if err != nil {
		if errors.Is(err, repository.ErrSaveProduct) {
			s.log.Warn("products empty", err)
			return fmt.Errorf("%s: %w", op, ErrorInvalidCredentials)
		}

		log.Error("products didnt added", err)
		return fmt.Errorf("%s %w", op, err)
	}

	log.Info("products added")

	return nil
}
