package handlers

import (
	"cryptocurrency/internal/middleware"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
)

type checkPromoCodeReq struct {
	TelegramUserID string `json:"telegramUserID"`
	PromoCode      string `json:"promoCode"`
}

// CheckPromoCode Change account price observers
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
func (h *Handler) CheckPromoCode(c *gin.Context) {
	var req checkPromoCodeReq

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

	//convert telegramUserId to int64
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
	err = h.LogService.LogRequestFromBot(ctx, "/account/promo_code/check", int64(telegramUserID), req)
	if err != nil {
		log.Fatal(err)
	}

	// get status of the promocode
	status, err := h.AccountService.CheckPromoCode(ctx, req.PromoCode, int64(telegramUserID))
	if err != nil {
		log.Printf("Failed to get status of the promocode: %v\n", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
		return
	}

	// type of status
	// time of active is ended
	// you already activate this promocode
	// promocode activated
	// limit of activation
	// no promo code in db

	c.JSON(http.StatusOK, gin.H{
		"result": status,
	})
}
