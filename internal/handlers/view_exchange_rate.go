package handlers

import (
	"cryptocurrency/internal/middleware"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type viewExchangeRateReq struct {
	TelegramUserID   string `json:"telegramUserID"`
	ExchangeFrom     string `json:"exchangeFrom"`
	ExchangeFromType string `json:"exchangeFromType"`
	ExchangeTo       string `json:"exchangeTo"`
	ExchangeToType   string `json:"exchangeToType"`
	LimitCurrency    string `json:"limitCurrency"`
	LimitValue       string `json:"limitValue"`
}

// ViewExchangeRate Create price observer
// @summary Create price observer
// @description Create price observer
// @Security Bearer
// @Tags Observers
// @Accept  json
// @Produce  json
// @Param   telegramUserId   body int64   true  "Telegram User ID"
// @Param   CryptoCurrency   body string  true  "Name of cryptocurrency which observe"
// @Param   CurrencyOfValue  body string  true  "Name of cryptocurrency to track"
// @Param   ExpectedValue    body float64 true  "Value that expected"
// @Param   IsUpDirection    body bool    true  "Up or down direction"
// @Success 200 {object} object
// @Failure 401 {object} object
// @Failure 404 {object} object
// @router /observers/price_observer/create [get]
func (h *Handler) ViewExchangeRate(c *gin.Context) {
	var req viewExchangeRateReq

	// get context data
	middlewareData, _ := c.Get(middleware.AuthorizationPayloadKey)
	data := strings.ReplaceAll(middlewareData.(string), "'", "\"")
	err := json.Unmarshal([]byte(data), &req)
	if err != nil {
		log.Printf("Failed to parse request body: %v\n", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
		return
	}

	// convert telegramUserID to int
	telegramUserId, err := strconv.Atoi(req.TelegramUserID)
	if err != nil {
		log.Printf("Failed to convert Telegram User ID from string to integer: %v\n", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
		return
	}

	// convert limitValue to float64
	limitValue, err := strconv.ParseFloat(req.LimitValue, 64)
	if err != nil {
		log.Printf("Failed to convert value of limit to : %v\n", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
		return
	}

	// log request
	ctx := c.Request.Context()
	err = h.LogService.LogRequestFromBot(ctx, "/info/exchange/bestchange", int64(telegramUserId), req)
	if err != nil {
		log.Printf("Failed to log request from bot: %v\n", err.Error())
	}

	// get variants for exchange
	exchangeVariants, err := h.InfoService.GetExchangeRateFromBestchange(ctx, req.ExchangeFrom, req.ExchangeFromType, req.ExchangeTo, req.ExchangeToType, req.LimitCurrency, limitValue)
	if err != nil {
		log.Printf("Failed to get variants for exchange: %v\n", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
		return
	}

	// return result with status
	c.JSON(http.StatusOK, gin.H{
		"result": exchangeVariants,
	})
}
