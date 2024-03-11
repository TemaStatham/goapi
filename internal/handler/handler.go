package handler

import (
	"context"
	"github.com/gin-gonic/gin"
	"log/slog"
)

type Handler struct {
	AuthService
	ProductService
	CategoryService
	log *slog.Logger
}

type AuthService interface {
	Login(ctx context.Context, email, password string) (token string, err error)
	Register(ctx context.Context, email, name, password string) (userID int64, err error)
}

type ProductService interface {
	Add(ctx context.Context)
	Delete(ctx context.Context)
	Edit(ctx context.Context)
	GetAllProducts(ctx context.Context)
}

type CategoryService interface {
	Add(ctx context.Context)
	Delete(ctx context.Context)
	Edit(ctx context.Context)
	GetAllProducts(ctx context.Context)
}

func NewHandler(a *AuthService, p *ProductService, c *CategoryService, l *slog.Logger) *Handler {
	return &Handler{
		AuthService:     *a,
		ProductService:  *p,
		CategoryService: *c,
		log:             l,
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
		auth.POST("/sign-up")
		auth.POST("/sign-in")
	}

	product := router.Group("/product")
	{
		product.POST("/add")
		product.POST("/delete")
		product.POST("/edit")
		product.POST("/get-all")
	}

	category := router.Group("/category")
	{
		category.POST("/add")
		category.POST("/delete")
		category.POST("/edit")
		category.POST("/get-all")
	}

	log.Info("Handler init")

	return router
}
