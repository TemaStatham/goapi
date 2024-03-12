package handler

import (
	"github.com/gin-gonic/gin"
	"goapi/internal/model"
	"log/slog"
	"net/http"
)

type addProductType struct {
	Name        string           `json:"name" binding:"required"`
	Categoryies []model.Category `json:"categoryies" binding:"required"`
}

func (h *Handler) addProduct(c *gin.Context) {
	const op = "handler.addProduct"

	log := h.log.With(
		slog.String("op", op),
	)

	var input addProductType

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, InvalidInputBodyErr)
		log.Error("error bind json:", InvalidInputBodyErr)
		return
	}

	id, err := h.product.AddProduct(c.Request.Context(), input.Name, input.Categoryies)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		log.Error("error added product", err.Error())
		return
	}

	log.Info("Handler product added")

	c.JSON(http.StatusOK, map[string]interface{}{
		"productID": id,
	})
}

type deleteProductType struct {
	ID int64 `json:"id" binding:"required"`
}

func (h *Handler) deleteProduct(c *gin.Context) {
	const op = "handler.deleteProduct"

	log := h.log.With(
		slog.String("op", op),
	)

	var input deleteProductType

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, InvalidInputBodyErr)
		log.Error("error bind json:", InvalidInputBodyErr)
		return
	}

	err := h.product.DeleteProduct(c.Request.Context(), input.ID)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		log.Error("error delete product", err.Error())
		return
	}

	log.Info("Handler product deleted")

	c.JSON(http.StatusOK, "success")
}

type editProduct struct {
	ID   int64  `json:"id" binding:"required"`
	Name string `json:"name" binding:"required"`
}

func (h *Handler) editProductName(c *gin.Context) {
	const op = "handler.editProduct"

	log := h.log.With(
		slog.String("op", op),
	)

	var input editProduct

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, InvalidInputBodyErr)
		log.Error("error bind json:", InvalidInputBodyErr)
		return
	}

	productID, err := h.product.EditProductName(c.Request.Context(), input.ID, input.Name)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		log.Error("error edit product", err.Error())
		return
	}

	log.Info("Handler product edited")

	c.JSON(http.StatusOK, map[string]interface{}{
		"productID": productID,
	})
}

type editProductCategoryiesType struct {
	ID          int64            `json:"id" binding:"required"`
	Categoryies []model.Category `json:"categoryies" binding:"required"`
}

func (h *Handler) editProductCategoryies(c *gin.Context) {
	const op = "handler.editProduct"

	log := h.log.With(
		slog.String("op", op),
	)

	var input editProductCategoryiesType

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, InvalidInputBodyErr)
		log.Error("error bind json:", InvalidInputBodyErr)
		return
	}

	productID, err := h.product.EditProductCategory(c.Request.Context(), input.ID, input.Categoryies)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		log.Error("error edit product", err.Error())
		return
	}

	log.Info("Handler product edited")

	c.JSON(http.StatusOK, map[string]interface{}{
		"productID": productID,
	})
}

type getAllProductsType struct {
	Tag string `json:"tag" binding:"required"`
}

func (h *Handler) getAllProducts(c *gin.Context) {
	const op = "handler.getAllProducts"

	log := h.log.With(
		slog.String("op", op),
	)

	var input getAllProductsType

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, InvalidInputBodyErr)
		log.Error("error bind json:", InvalidInputBodyErr)
		return
	}

	products, err := h.product.GetAllProducts(c.Request.Context(), input.Tag)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		log.Error("error getting all products", err.Error())
		return
	}

	log.Info("Handler getting all product")

	c.JSON(http.StatusOK, map[string]interface{}{
		"products": products,
	})
}
