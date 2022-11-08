package handlers

import (
	"cryptocurrency/internal/middleware"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"net/http"
	"strconv"
)

type deleteRecordPortfolioReq struct {
	TelegramUserID string `json:"telegramUserID"`
	RecordId       string `json:"record_id"`
}

// DeleteRecordPortfolio Delete currency record in portfolio
// @summary Delete in portfolio
// @description Delete currency record in portfolio
// @Security Bearer
// @Tags Account - Portfolio
// @Accept  json
// @Produce  json
// @Success 200 {object} object
// @Failure 401 {object} object
// @Failure 404 {object} object
// @router /account/portfolio/delete [get]
func (h *Handler) DeleteRecordPortfolio(c *gin.Context) {
	var req deleteRecordPortfolioReq

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

	// convert observerID to ObjectID
	recordID, err := primitive.ObjectIDFromHex(req.RecordId)
	if err != nil {
		log.Printf("Failed to convert record ID from string to ObjectID: %v\n", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
		return
	}

	// log request
	ctx := c.Request.Context()
	err = h.LogService.LogRequestFromBot(ctx, "/account/portfolio/delete", int64(telegramUserID), req)
	if err != nil {
		log.Printf("Failed to log request from bot: %v\n", err.Error())
	}

	// delete record from portfolio
	err = h.AccountService.DeleteElementUserPortfolio(ctx, recordID)
	if err != nil {
		log.Printf("Failed to delete record from portfolio: %v\n", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"result": "ok",
	})
}
