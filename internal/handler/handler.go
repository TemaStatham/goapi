package handler

import (
	"context"
	"github.com/gin-gonic/gin"
	"goapi/internal/model"
	"log/slog"
)

//go:generate mockgen -source=handler.go -destination=mock/mock.go

const (
	InvalidInputBodyErr = "invalid input body"
)

type Handler struct {
	auth     AuthService
	product  ProductService
	category CategoryService
	log      *slog.Logger
}

type AuthService interface {
	Login(ctx context.Context, email, password string) (token string, err error)
	Register(ctx context.Context, email, password string) (userID int64, err error)
}

type ProductService interface {
	AddProduct(ctx context.Context, name string, categoryies []string) (int64, error)
	DeleteProduct(ctx context.Context, id int64) error
	EditProductName(ctx context.Context, id int64, name string) (int64, error)
	EditProductCategory(ctx context.Context, id int64, categoryies []model.Category) (int64, error)
	GetAllProducts(ctx context.Context, tag string) ([]model.Product, error)
	GetCategoryProducts(ctx context.Context, category string) ([]model.Product, error)
}

type CategoryService interface {
	AddCategory(ctx context.Context, name string) (int64, error)
	DeleteCategory(ctx context.Context, id int64) error
	EditCategory(ctx context.Context, id int64, name string) (int64, error)
	GetAllCategoryies(ctx context.Context, tag string) ([]model.Category, error)
}

func NewHandler(a AuthService, p ProductService, c CategoryService, l *slog.Logger) *Handler {
	return &Handler{
		auth:     a,
		product:  p,
		category: c,
		log:      l,
	}
}

func (h *Handler) Init() *gin.Engine {
	const op = "handler.init"

	log := h.log.With(
		slog.String("op", op),
	)

	router := gin.New()

	auth := router.Group("/auth")
	{
		auth.POST("/sign-up", h.signUp)
		auth.POST("/sign-in", h.signIn)
	}

	api := router.Group("/api", h.userIdentity)
	{
		product := api.Group("/product")
		{
			product.POST("/add", h.addProduct)
			product.POST("/delete", h.deleteProduct)
			product.POST("/edit-name", h.editProductName)
			product.POST("/edit-categoryies", h.editProductCategoryies)
			product.POST("/get-all", h.getAllProducts)
			product.POST("/get", h.getProducts)
		}

		category := api.Group("/category")
		{
			category.POST("/add", h.addCategory)
			category.POST("/delete", h.editCategory)
			category.POST("/edit", h.deleteCategory)
			category.POST("/get-all", h.getAllCategory)
		}
	}

	log.Info("Handler init")

	return router
}
