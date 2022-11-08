package handlers

import (
	"cryptocurrency/internal/middleware"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
)

type getAccountPercentageObserversReq struct {
	TelegramUserID string `json:"telegramUserID"`
}

// GetAccountPercentageObservers View account percentage observers
// @summary View percentage observers
// @description View account percentage observers
// @Security Bearer
// @Tags Account - Percentage observers
// @Accept  json
// @Produce  json
// @Success 200 {object} object
// @Failure 401 {object} object
// @Failure 404 {object} object
// @router /account/percentage_observers/list [get]
func (h *Handler) GetAccountPercentageObservers(c *gin.Context) {
	var req getAccountPercentageObserversReq

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
	err = h.LogService.LogRequestFromBot(ctx, "/account/percentage_observers/list", int64(telegramUserID), req)
	if err != nil {
		log.Printf("Failed to log request from bot: %v\n", err.Error())
	}

	// get user percentage observers
	data, err := h.AccountService.GetUserPercentageObservers(ctx, int64(telegramUserID))
	if err != nil {
		log.Printf("Failed to get user percentage observers: %v\n", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"result": data,
	})
}
