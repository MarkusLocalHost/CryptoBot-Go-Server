package handlers

import (
	"cryptocurrency/internal/middleware"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
)

type viewPortfolioInfoReq struct {
	TelegramUserID string `json:"telegramUserID"`
}

// ViewPortfolio View currencies in user portfolio
// @summary View portfolio
// @description View currencies in user portfolio
// @Security Bearer
// @Tags Account - Portfolio
// @Accept  json
// @Produce  json
// @Success 200 {object} object
// @Failure 401 {object} object
// @Failure 404 {object} object
// @router /account/portfolio/view [get]
func (h *Handler) ViewPortfolio(c *gin.Context) {
	var req viewPortfolioInfoReq

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
	err = h.LogService.LogRequestFromBot(ctx, "/account/portfolio/view", int64(telegramUserID), req)
	if err != nil {
		log.Printf("Failed to log request from bot: %v\n", err.Error())
	}

	// get data in user portfolio
	data, err := h.AccountService.ViewUserPortfolio(ctx, int64(telegramUserID))
	if err != nil {
		log.Printf("Failed to get user portfolio: %v\n", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
		return
	}

	// get actual price for currency in portfolio
	dataWithActualPrices, err := h.InfoService.GetPriceForSymbolCurrenciesInPortfolio(data)
	if err != nil {
		log.Printf("Failed to get actual prices: %v\n", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"result": dataWithActualPrices,
	})
}
