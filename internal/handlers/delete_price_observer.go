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

type deleteAccountPriceObserverReq struct {
	TelegramUserID string `json:"telegramUserID"`
	ObserverID     string `json:"observerID"`
}

// DeleteAccountPriceObserver Delete account price observers
// @summary Delete price observers
// @description Delete account price observers
// @Security Bearer
// @Tags Account - Price observers
// @Accept  json
// @Produce  json
// @Success 200 {object} object
// @Failure 401 {object} object
// @Failure 404 {object} object
// @router /account/price_observers/delete [get]
func (h *Handler) DeleteAccountPriceObserver(c *gin.Context) {
	var req deleteAccountPriceObserverReq

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
	observerID, _ := primitive.ObjectIDFromHex(req.ObserverID)
	if err != nil {
		log.Printf("Failed to convert observer ID to ObjectID: %v\n", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
		return
	}

	// log request
	ctx := c.Request.Context()
	err = h.LogService.LogRequestFromBot(ctx, "/account/price_observers/delete", int64(telegramUserID), req)
	if err != nil {
		log.Printf("Failed to log request from bot: %v\n", err.Error())
	}

	// delete user price observer
	err = h.ObserverService.DeletePriceObserver(ctx, observerID)
	if err != nil {
		log.Printf("Failed to delete user observer: %v\n", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"result": "ok",
	})
}
