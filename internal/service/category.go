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
	ErrCategoryId        = -1
	TagGetAllCategoryies = "get all products"
)

//go:generate mockgen -source=category.go -destination=mock/category_mock.go

var (
	ErrCategoryNameIsEmpty = errors.New("category name is empty")
	ErrCategoryIDIsEmpty   = errors.New("category id is empty")
	ErrCategoryUnknownTag  = errors.New("unknown tag get all products")
)

type CategoryService struct {
	adder   AdderCategory
	deleter DeleterCategory
	updater UpdaterCategory
	getter  GetterCategory
	log     *slog.Logger
}

type AdderCategory interface {
	AddCategory(ctx context.Context, name string) (int64, error)
}

type DeleterCategory interface {
	DeleteCategory(ctx context.Context, id int64) error
}

type UpdaterCategory interface {
	UpdateCategoryName(ctx context.Context, id int64, name string) (int64, error)
}

type GetterCategory interface {
	GetAllCategoryies(ctx context.Context) ([]model.Category, error)
}

func NewCategoryService(
	a AdderCategory,
	d DeleterCategory,
	u UpdaterCategory,
	g GetterCategory,
	l *slog.Logger,
) *CategoryService {
	return &CategoryService{
		adder:   a,
		deleter: d,
		updater: u,
		getter:  g,
		log:     l,
	}
}

func (s *CategoryService) AddCategory(ctx context.Context, name string) (int64, error) {
	const op = "category.AddCategory"

	log := s.log.With(
		slog.String("op", op),
		slog.String("name", name),
	)

	log.Info("add category")

	if name == "" {
		log.Info("name is empty", ErrCategoryNameIsEmpty)
		return ErrCategoryId, fmt.Errorf("%s %w", op, ErrCategoryNameIsEmpty)
	}

	categoryID, err := s.adder.AddCategory(ctx, name)
	if err != nil {
		if errors.Is(err, repository.ErrCategoryExist) {
			s.log.Warn("category already exist", err)
			return ErrCategoryId, fmt.Errorf("%s: %w", op, ErrorInvalidCredentials)
		}

		log.Info("category didnt added")
		return ErrCategoryId, fmt.Errorf("%s %w", op, err)
	}

	log.Info("category is saved")

	return categoryID, nil
}
func (s *CategoryService) DeleteCategory(ctx context.Context, id int64) error {
	const op = "category.DeleteCategory"

	log := s.log.With(
		slog.String("op", op),
		slog.Int64("id", id),
	)

	log.Info("delete category")

	if id <= 0 {
		log.Info("id is empty", ErrCategoryIDIsEmpty)
		return fmt.Errorf("%s %w", op, ErrCategoryIDIsEmpty)
	}

	err := s.deleter.DeleteCategory(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrCategoryDelete) {
			s.log.Warn("user not found", err)
			return fmt.Errorf("%s: %w", op, ErrorInvalidCredentials)
		}
		if errors.Is(err, repository.ErrDeleteProductCategory) {
			s.log.Warn("category not exist", err)
			return fmt.Errorf("%s: %w", op, ErrorInvalidCredentials)
		}

		log.Info("category didnt deleted")
		return fmt.Errorf("%s %w", op, err)
	}

	log.Info("category is deleted")

	return nil
}

func (s *CategoryService) EditCategory(ctx context.Context, id int64, name string) (int64, error) {
	const op = "category.DeleteCategory"

	log := s.log.With(
		slog.String("op", op),
		slog.Int64("id", id),
	)

	log.Info("delete category")

	if id <= 0 {
		log.Info("id is empty", ErrCategoryIDIsEmpty)
		return ErrCategoryId, fmt.Errorf("%s %w", op, ErrCategoryIDIsEmpty)
	}

	if name == "" {
		log.Info("name is empty", ErrCategoryNameIsEmpty)
		return ErrCategoryId, fmt.Errorf("%s %w", op, ErrCategoryNameIsEmpty)
	}

	categoryID, err := s.updater.UpdateCategoryName(ctx, id, name)
	if err != nil {
		if errors.Is(err, repository.ErrUpdateCategory) {
			s.log.Warn("categoryies not updated", err)
			return ErrCategoryId, fmt.Errorf("%s: %w", op, ErrorInvalidCredentials)
		}

		log.Info("category didnt deleted")
		return ErrCategoryId, fmt.Errorf("%s %w", op, err)
	}

	log.Info("category is deleted")

	return categoryID, nil
}
func (s *CategoryService) GetAllCategoryies(ctx context.Context, tag string) ([]model.Category, error) {
	const op = "category.GetAllCategoryies"

	log := s.log.With(
		slog.String("op", op),
		slog.String("tag", tag),
	)

	log.Info("get all categoryies")

	if tag != TagGetAllCategoryies {
		log.Error("categoryies didnt get", ErrCategoryUnknownTag)
		return []model.Category{}, fmt.Errorf("%s %w", op, ErrCategoryUnknownTag)
	}

	categoryies, err := s.getter.GetAllCategoryies(ctx)
	if err != nil {
		if errors.Is(err, repository.ErrAllCategoryies) {
			s.log.Warn("categoryies empty", err)
			return []model.Category{}, fmt.Errorf("%s: %w", op, ErrorInvalidCredentials)
		}

		log.Error("categoryies didnt get", err)
		return []model.Category{}, fmt.Errorf("%s %w", op, err)
	}

	log.Info("all categoryies is getter")

	return categoryies, nil
}
