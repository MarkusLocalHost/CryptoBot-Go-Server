package handlers

import (
	"cryptocurrency/internal/middleware"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
)

type getSupportedVSCurrenciesReq struct {
	TelegramUserID string `json:"telegramUserID"`
}

// GetSupportedVSCurrencies Get currency which enable to track
// @summary View supported vs currencies
// @description View currency which enable to track
// @Security Bearer
// @Tags Info
// @Accept  json
// @Produce  json
// @Success 200 {object} object
// @Failure 401 {object} object
// @Failure 404 {object} object
// @router /info/supported_vs_currencies/view [get]
func (h *Handler) GetSupportedVSCurrencies(c *gin.Context) {
	var req getSupportedVSCurrenciesReq

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
	err = h.LogService.LogRequestFromBot(ctx, "/info/supported_vs_currencies/view", int64(telegramUserID), req)
	if err != nil {
		log.Printf("Failed to log request from bot: %v\n", err.Error())
	}

	// get info about supported currencies for price observers
	data, err := h.InfoService.GetSupportedVSCurrencies(ctx)
	if err != nil {
		log.Printf("Failed to get info about supported currencies for price observers: %v\n", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"result": "observer started",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"result": data,
	})
}
