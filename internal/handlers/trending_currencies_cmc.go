package handlers

import (
	"cryptocurrency/internal/middleware"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
)

type getTrendingCurrenciesReq struct {
	TelegramUserID string `json:"telegramUserID"`
	Source         string `json:"source"`
}

// GetTrendingCurrencies Get trending currency from coinmarketplace
// @summary Trending from coinmarketplace
// @description Get trending currency from coinmarketplace
// @Security Bearer
// @Tags Info - Trending
// @Accept  json
// @Produce  json
// @Success 200 {object} object
// @Failure 401 {object} object
// @Failure 404 {object} object
// @router /info/trending/cmc [get]
func (h *Handler) GetTrendingCurrencies(c *gin.Context) {
	var req getTrendingCurrenciesReq

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
	err = h.LogService.LogRequestFromBot(ctx, "/info/trending", int64(telegramUserID), req)
	if err != nil {
		log.Printf("Failed to log request from bot: %v\n", err.Error())
	}

	// get trending currencies
	data, err := h.InfoService.GetTrendingCurrencies(ctx, req.Source)
	if err != nil {
		log.Printf("Failed to get trending currencies: %v\n", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"result": data,
	})
}
