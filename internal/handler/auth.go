package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log/slog"
	"net/http"
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
		log.Error("error bind json:", fmt.Errorf("%s", InvalidInputBodyErr))
		return
	}

	token, err := h.auth.Login(c.Request.Context(), input.Email, input.Password)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		log.Error("error login", err)
		return
	}

	log.Info("Handler sign in")

	c.JSON(http.StatusOK, map[string]interface{}{
		"token": token,
	})
}

type signUpType struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (h *Handler) signUp(c *gin.Context) {
	const op = "handler.signUp"

	log := h.log.With(
		slog.String("op", op),
	)

	var input signUpType

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, InvalidInputBodyErr)
		log.Error("error bind json:", fmt.Errorf("%s", InvalidInputBodyErr))
		return
	}

	id, err := h.auth.Register(c.Request.Context(), input.Email, input.Password)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		log.Error("error register ", err)
		return
	}

	log.Info("Handler sign up")

	c.JSON(http.StatusOK, map[string]interface{}{
		"id": id,
	})
}
