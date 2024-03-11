package handler

import (
	"github.com/gin-gonic/gin"
	"goapi/internal/model"
	"log/slog"
	"net/http"
)

const (
	InvalidInputBodyErr = "invalid input body"
)

type signInInput struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (h *Handler) signIn(c *gin.Context) {
	const op = "handler.signIn"

	log := h.log.With(
		slog.String("op", op),
	)

	var input signInInput

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, InvalidInputBodyErr)
		log.Error("error bind json:", InvalidInputBodyErr)
		return
	}

	token, err := h.AuthService.Login(c.Request.Context(), input.Email, input.Password)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		log.Error("error login", err.Error())
		return
	}

	log.Info("Handler sign in")

	c.JSON(http.StatusOK, map[string]interface{}{
		"token": token,
	})
}

func (h *Handler) signUp(c *gin.Context) {
	const op = "handler.signUp"

	log := h.log.With(
		slog.String("op", op),
	)

	var input model.User

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, InvalidInputBodyErr)
		log.Error("error bind json:", InvalidInputBodyErr)
		return
	}

	id, err := h.Register(c.Request.Context(), input.Email, input.Name, input.Password)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		log.Error("error register ", err.Error())
		return
	}

	log.Info("Handler sign up")

	c.JSON(http.StatusOK, map[string]interface{}{
		"id": id,
	})
}
