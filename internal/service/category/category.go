package category

import (
	"context"
	"errors"
	"fmt"
	"goapi/internal/model"
	"log/slog"
)

const (
	ErrCategoryId        = -1
	TagGetAllCategoryies = "get all products"
)

var (
	ErrNameIsEmpty = errors.New("category name is empty")
	ErrIDIsEmpty   = errors.New("category id is empty")
	ErrUnknownTag  = errors.New("unknown tag get all products")
)

type Service struct {
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
) *Service {
	return &Service{
		adder:   a,
		deleter: d,
		updater: u,
		getter:  g,
		log:     l,
	}
}

func (s *Service) AddCategory(ctx context.Context, name string) (int64, error) {
	const op = "category.AddCategory"

	log := s.log.With(
		slog.String("op", op),
		slog.String("name", name),
	)

	log.Info("add category")

	if name == "" {
		log.Info("name is empty", ErrNameIsEmpty)
		return ErrCategoryId, fmt.Errorf("%s %w", op, ErrNameIsEmpty)
	}

	categoryID, err := s.adder.AddCategory(ctx, name)
	if err != nil {
		log.Info("category didnt added")
		return ErrCategoryId, fmt.Errorf("%s %w", op, err)
	}

	log.Info("category is saved")

	return categoryID, nil
}
func (s *Service) DeleteCategory(ctx context.Context, id int64) error {
	const op = "category.DeleteCategory"

	log := s.log.With(
		slog.String("op", op),
		slog.Int64("id", id),
	)

	log.Info("delete category")

	if id <= 0 {
		log.Info("id is empty", ErrIDIsEmpty)
		return fmt.Errorf("%s %w", op, ErrIDIsEmpty)
	}

	err := s.deleter.DeleteCategory(ctx, id)
	if err != nil {
		log.Info("category didnt deleted")
		return fmt.Errorf("%s %w", op, err)
	}

	log.Info("category is deleted")

	return nil
}

func (s *Service) EditCategory(ctx context.Context, id int64, name string) (int64, error) {
	const op = "category.DeleteCategory"

	log := s.log.With(
		slog.String("op", op),
		slog.Int64("id", id),
	)

	log.Info("delete category")

	if id <= 0 {
		log.Info("id is empty", ErrIDIsEmpty)
		return ErrCategoryId, fmt.Errorf("%s %w", op, ErrIDIsEmpty)
	}

	if name == "" {
		log.Info("name is empty", ErrNameIsEmpty)
		return ErrCategoryId, fmt.Errorf("%s %w", op, ErrNameIsEmpty)
	}

	categoryID, err := s.updater.UpdateCategoryName(ctx, id, name)
	if err != nil {
		log.Info("category didnt deleted")
		return ErrCategoryId, fmt.Errorf("%s %w", op, err)
	}

	log.Info("category is deleted")

	return categoryID, nil
}
func (s *Service) GetAllCategoryies(ctx context.Context, tag string) ([]model.Category, error) {
	const op = "category.GetAllCategoryies"

	log := s.log.With(
		slog.String("op", op),
		slog.String("tag", tag),
	)

	log.Info("get all categoryies")

	if tag != TagGetAllCategoryies {
		log.Error("categoryies didnt get", ErrUnknownTag)
		return []model.Category{}, fmt.Errorf("%s %w", op, ErrUnknownTag)
	}

	categoryies, err := s.getter.GetAllCategoryies(ctx)
	if err != nil {
		log.Error("categoryies didnt get", err)
		return []model.Category{}, fmt.Errorf("%s %w", op, err)
	}

	log.Info("all categoryies is getter")

	return categoryies, nil
}
