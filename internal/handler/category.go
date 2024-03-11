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

	id, err := h.CategoryService.Add(c.Request.Context(), input.Name)
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

func (h *Handler) deleteCategory(c *gin.Context) {

}

func (h *Handler) editCategory(c *gin.Context) {

}

func (h *Handler) getAllCategory(c *gin.Context) {

}
