package handlers

import (
	"cryptocurrency/internal/middleware"
	"cryptocurrency/internal/utils/apperrors"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
)

type extendUserSubscriptionReq struct {
	TelegramUserID string `json:"telegramUserID"`
	Hours          string `json:"hours"`
}

// ExtendUserSubscription Change account price observers
// @summary Change price observers
// @description Change account price observers
// @Security Bearer
// @Tags Account - Price observers
// @Accept  json
// @Produce  json
// @Success 200 {object} object
// @Failure 401 {object} object
// @Failure 404 {object} object
// @router /account/price_observers/change [get]
func (h *Handler) ExtendUserSubscription(c *gin.Context) {
	var req extendUserSubscriptionReq

	// get context data
	middlewareData, _ := c.Get(middleware.AuthorizationPayloadKey)
	err := json.Unmarshal([]byte(middlewareData.(string)), &req)
	if err != nil {
		log.Println(err)
	}

	//convert telegramUserId to int64
	telegramUserID, _ := strconv.Atoi(req.TelegramUserID)

	// convert hours to int
	hours, _ := strconv.Atoi(req.Hours)
	if err != nil {
		log.Fatal(err)
	}

	ctx := c.Request.Context()
	err = h.AccountService.ExtendUserSubscription(ctx, int64(telegramUserID), hours)
	if err != nil {
		log.Printf("Failed to change status observer: %v\n", err.Error())
		c.JSON(
			apperrors.Status(err),
			gin.H{"error": err})
		return
	}

	// type of status
	// time of active is ended
	// you already activate this promo code
	// promo code activated
	// limit of activation
	// no promo code in db

	c.JSON(http.StatusCreated, gin.H{
		"result": "ok",
	})
}
