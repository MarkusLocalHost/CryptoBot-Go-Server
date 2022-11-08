package handlers

import (
	"cryptocurrency/internal/middleware"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
)

type newAccountReq struct {
	TelegramUserID string `json:"telegramUserID"`
	Language       string `json:"language"`
}

// NewAccount Create account for new user
// @summary Create account
// @description Create account for new user
// @Security Bearer
// @Tags Account
// @Accept  json
// @Produce  json
// @Success 200 {object} object
// @Failure 401 {object} object
// @Failure 404 {object} object
// @router /account/new_account [get]
func (h *Handler) NewAccount(c *gin.Context) {
	var req newAccountReq

	// get context data
	middlewareData, _ := c.Get(middleware.AuthorizationPayloadKey)
	err := json.Unmarshal([]byte(middlewareData.(string)), &req)
	if err != nil {
		log.Printf("Failed to parse request body: %v\n", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
		return
	}

	// convert telegramUserId to int64
	telegramUserID, err := strconv.Atoi(req.TelegramUserID)
	if err != nil {
		log.Printf("Failed to convert Telegram User ID from string to integer: %v\n", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
		return
	}

	// log request
	ctx := c.Request.Context()
	err = h.LogService.LogRequestFromBot(ctx, "/account/new_account", int64(telegramUserID), req)
	if err != nil {
		log.Printf("Failed to log request from bot: %v\n", err.Error())
	}

	// create account
	err = h.AccountService.CreateAccount(ctx, int64(telegramUserID), req.Language)
	if err != nil {
		log.Printf("Failed to create account: %v\n", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"result": "user in bd",
	})
}
