package handlers

import (
	"cryptocurrency/internal/middleware"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
)

type tryFindCurrencyInfoReq struct {
	TelegramUserID string `json:"telegramUserID"`
	CurrencyName   string `json:"currency_name"`
	CurrencySlug   string `json:"currency_slug"`
}

func (h *Handler) TryFindCurrency(c *gin.Context) {
	var req tryFindCurrencyInfoReq

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
	err = h.LogService.LogRequestFromBot(ctx, "/search/try_find_currency", int64(telegramUserID), req)
	if err != nil {
		log.Printf("Failed to log request from bot: %v\n", err.Error())
	}

	// try to find currency in bd
	data, err := h.InfoService.TryFindCurrencyByNameOrSlug(ctx, req.CurrencyName, req.CurrencySlug)
	if err != nil {
		log.Printf("Failed to find currency: %v\n", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"result": data,
	})
}
