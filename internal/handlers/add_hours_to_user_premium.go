package handlers

import (
	"cryptocurrency/internal/middleware"
	"cryptocurrency/internal/utils/apperrors"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type addHoursToUserSubscriptionReq struct {
	UserID         string `json:"userID"`
	TelegramUserId string `json:"telegramUserId"`
	Hours          string `json:"hours"`
}

// AddHoursToUserSubscription Create price observer
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
func (h *Handler) AddHoursToUserSubscription(c *gin.Context) {
	var req addHoursToUserSubscriptionReq

	// get context data
	middlewareData, _ := c.Get(middleware.AuthorizationPayloadKey)
	data := middlewareData.(string)
	data = strings.ReplaceAll(data, "'", "\"")
	err := json.Unmarshal([]byte(data), &req)
	if err != nil {
		log.Println(err)
	}

	// convert telegramUserId to int64
	telegramUserID, _ := strconv.Atoi(req.TelegramUserId)
	log.Println(telegramUserID)

	// convert hours to int
	hours, _ := strconv.Atoi(req.Hours)
	if err != nil {
		log.Fatal(err)
	}

	// log request
	ctx := c.Request.Context()
	err = h.LogService.LogRequestFromSite(ctx, "/manager/account/set_premium", req.UserID, req)
	if err != nil {
		log.Fatal(err)
	}

	// save observer to db
	err = h.ManagerService.AddHoursToUserSubscription(ctx, int64(telegramUserID), hours)
	if err != nil {
		log.Printf("Failed to create promo code: %v\n", err.Error())
		c.JSON(
			apperrors.Status(err),
			gin.H{"error": err})
		return
	}

	// return result with status
	c.JSON(http.StatusCreated, gin.H{
		"result": "promo code created",
	})
}
