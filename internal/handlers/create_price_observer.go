package handlers

import (
	"cryptocurrency/internal/middleware"
	"cryptocurrency/internal/models/mongo_models"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"net/http"
	"strconv"
	"time"
)

type createPriceObserverReq struct {
	TelegramUserID  string `json:"telegramUserID"`
	CryptoID        string `json:"cryptoID"`
	CurrencyOfValue string `json:"currencyOfValue"`
	ExpectedValue   string `json:"expectedValue"`
	IsUpDirection   string `json:"isUpDirection"`
}

// CreatePriceObserver Create price observer
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
func (h *Handler) CreatePriceObserver(c *gin.Context) {
	var req createPriceObserverReq

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

	// convert expectedValue to float64
	expectedValue, err := strconv.ParseFloat(req.ExpectedValue, 64)
	if err != nil {
		log.Printf("Failed to convert expected value from string to float: %v\n", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
		return
	}

	// convert boolean
	var isUpDirection bool
	if req.IsUpDirection == "True" {
		isUpDirection = true
	} else if req.IsUpDirection == "False" {
		isUpDirection = false
	}

	o := &mongo_models.Observer{
		Id:              primitive.NewObjectID(),
		TelegramUserID:  int64(telegramUserID),
		CryptoID:        req.CryptoID,
		CurrencyOfValue: req.CurrencyOfValue,
		ExpectedValue:   expectedValue,
		IsUpDirection:   isUpDirection,
		CreatedAt:       time.Now(),
	}

	// log request
	ctx := c.Request.Context()
	err = h.LogService.LogRequestFromBot(ctx, "/observers/price_observer/create", int64(telegramUserID), req)
	if err != nil {
		log.Printf("Failed to log request from bot: %v\n", err.Error())
	}

	// save price observer to db
	observerIsActive, err := h.ObserverService.CreatePriceObserver(ctx, o)
	if err != nil {
		log.Printf("Failed to create price observer: %v\n", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
		return
	}

	// return result with status
	if observerIsActive {
		c.JSON(http.StatusCreated, gin.H{
			"result": "observer started",
		})
	} else {
		c.JSON(http.StatusAccepted, gin.H{
			"result": "observer created",
		})
	}
}
