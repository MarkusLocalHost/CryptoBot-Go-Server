package handlers

import (
	"cryptocurrency/internal/middleware"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
)

type getIndexRatingReq struct {
	TelegramUserID string `json:"telegramUserID"`
	Page           string `json:"page"`
	Currency       string `json:"currency"`
}

// GetIndexRating Get data to rating currency
// @summary Currency rating
// @description Get data to rating currency
// @Security Bearer
// @Tags Info - Rating
// @Accept  json
// @Produce  json
// @Success 200 {object} object
// @Failure 401 {object} object
// @Failure 404 {object} object
// @router /info/index_rating [get]
func (h *Handler) GetIndexRating(c *gin.Context) {
	var req getIndexRatingReq

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

	// convert page to int
	page, err := strconv.Atoi(req.Page)
	if err != nil {
		log.Printf("Failed to convert page from string to integer: %v\n", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
		return
	}

	// log request
	ctx := c.Request.Context()
	err = h.LogService.LogRequestFromBot(ctx, "/info/index_rating", int64(telegramUserID), req)
	if err != nil {
		log.Printf("Failed to log request from bot: %v\n", err.Error())
	}

	// get index price by page
	data, err := h.InfoService.GetIndexPriceByPage(ctx, page, req.Currency)
	if err != nil {
		log.Printf("Failed to get index price by page: %v\n", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"result": data,
	})
}
