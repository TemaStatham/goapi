package handler

import (
	"github.com/gin-gonic/gin"
	"log/slog"
	"net/http"
)

type addCategory struct {
	Name string `json:"name" binding:"required"`
}

func (h *Handler) addCategory(c *gin.Context) {
	const op = "handler.addCategory"

	log := h.log.With(
		slog.String("op", op),
	)

	var input addCategory

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, InvalidInputBodyErr)
		log.Error("error bind json:", InvalidInputBodyErr)
		return
	}

	id, err := h.CategoryService.AddCategory(c.Request.Context(), input.Name)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		log.Error("error added category", err.Error())
		return
	}

	log.Info("Handler category added")

	c.JSON(http.StatusOK, map[string]interface{}{
		"categoryID": id,
	})
}

type deleteCategory struct {
	Name string `json:"name" binding:"required"`
}

func (h *Handler) deleteCategory(c *gin.Context) {
	const op = "handler.deleteCategory"

	log := h.log.With(
		slog.String("op", op),
	)

	var input deleteCategory

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, InvalidInputBodyErr)
		log.Error("error bind json:", InvalidInputBodyErr)
		return
	}

	err := h.CategoryService.DeleteCategory(c.Request.Context(), input.Name)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		log.Error("error added category", err.Error())
		return
	}

	log.Info("Handler category deleted")

	c.JSON(http.StatusOK, "success")
}

type editCategory struct {
	Name string `json:"name" binding:"required"`
}

func (h *Handler) editCategory(c *gin.Context) {
	const op = "handler.editCategory"

	log := h.log.With(
		slog.String("op", op),
	)

	var input editCategory

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, InvalidInputBodyErr)
		log.Error("error bind json:", InvalidInputBodyErr)
		return
	}

	id, err := h.CategoryService.EditCategory(c.Request.Context(), input.Name)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		log.Error("error added category", err.Error())
		return
	}

	log.Info("Handler category edited")

	c.JSON(http.StatusOK, map[string]interface{}{
		"categoryID": id,
	})
}

type getAllCategoryiesType struct {
	Tag string `json:"tag" binding:"required"`
}

func (h *Handler) getAllCategory(c *gin.Context) {
	const op = "handler.getAllCategoryies"

	log := h.log.With(
		slog.String("op", op),
	)

	var input getAllCategoryiesType

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, InvalidInputBodyErr)
		log.Error("error bind json:", InvalidInputBodyErr)
		return
	}

	categoryies, err := h.CategoryService.GetAllCategoryies(c.Request.Context(), input.Tag)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		log.Error("error getting all products", err.Error())
		return
	}

	log.Info("Handler getting all categoryies")

	c.JSON(http.StatusOK, map[string]interface{}{
		"categoryies": categoryies,
	})
}
