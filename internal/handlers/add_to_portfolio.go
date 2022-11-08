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

type addToPortfolioReq struct {
	TelegramUserID string `json:"telegramUserID"  bson:"telegram_user_id"`
	Cryptocurrency string `json:"cryptocurrency"  bson:"cryptocurrency"`
	Value          string `json:"value"           bson:"value"`
	Type           string `json:"type"            bson:"type"`
	Price          string `json:"price"           bson:"price"`
}

// AddToPortfolio Add currency in user portfolio
// @summary Add to portfolio
// @description Add currency in user portfolio
// @Security Bearer
// @Tags Account - Portfolio
// @Accept  json
// @Produce  json
// @Success 200 {object} object
// @Failure 401 {object} object
// @Failure 404 {object} object
// @router /account/portfolio/add [get]
func (h *Handler) AddToPortfolio(c *gin.Context) {
	var req addToPortfolioReq

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

	// convert value to float64
	value, err := strconv.ParseFloat(req.Value, 64)
	if err != nil {
		log.Printf("Failed to convert value from string to float: %v\n", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
		return
	}

	// convert price to float64
	price, err := strconv.ParseFloat(req.Price, 64)
	if err != nil {
		log.Printf("Failed to convert price from string to float: %v\n", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
		return
	}

	p := &mongo_models.Portfolio{
		Id:             primitive.NewObjectID(),
		Cryptocurrency: req.Cryptocurrency,
		TelegramUserID: int64(telegramUserID),
		Value:          value,
		Price:          price,
		Type:           req.Type,
		CreatedAt:      time.Now(),
	}

	// log request
	ctx := c.Request.Context()
	err = h.LogService.LogRequestFromBot(ctx, "/account/portfolio/add", int64(telegramUserID), req)
	if err != nil {
		log.Printf("Failed to log request from bot: %v\n", err.Error())
	}

	// add currency to user portfolio
	err = h.AccountService.AddToUserPortfolio(ctx, p)
	if err != nil {
		log.Printf("Failed to add currency to portfolio: %v\n", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"result": "ok",
	})
}
