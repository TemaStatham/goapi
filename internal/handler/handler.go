package handler

import (
	"context"
	"github.com/gin-gonic/gin"
	"goapi/internal/model"
	"log/slog"
)

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
	AddProduct(ctx context.Context, name string, categoryies []string) (productID int64, err error)
	DeleteProduct(ctx context.Context, name string) error
	EditProduct(ctx context.Context, name string) (productID int64, err error)
	GetAllProducts(ctx context.Context, tag string) (product []model.Product, err error)
}

type CategoryService interface {
	AddCategory(ctx context.Context, name string) (categoryID int64, err error)
	DeleteCategory(ctx context.Context, name string) error
	EditCategory(ctx context.Context, name string) (categoryID int64, err error)
	GetAllCategoryies(ctx context.Context, tag string) (categoryies []model.Category, err error)
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

	product := router.Group("/product")
	{
		product.POST("/add", h.addProduct)
		product.POST("/delete", h.editProduct)
		product.POST("/edit", h.deleteProduct)
		product.POST("/get-all", h.getAllProducts)
	}

	category := router.Group("/category")
	{
		category.POST("/add", h.addCategory)
		category.POST("/delete", h.editCategory)
		category.POST("/edit", h.deleteCategory)
		category.POST("/get-all", h.getAllCategory)
	}

	log.Info("Handler init")

	return router
}
