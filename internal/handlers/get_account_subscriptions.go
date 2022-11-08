package handlers

import (
	"cryptocurrency/internal/middleware"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
	"strconv"
)

type getAccountSubscriptionReq struct {
	TelegramUserID string `json:"telegramUserID"`
}

// GetAccountSubscription View data about subscription
// @summary View subscription
// @description View data about subscription
// @Security Bearer
// @Tags Account - Subscription
// @Accept  json
// @Produce  json
// @Success 200 {object} object
// @Failure 401 {object} object
// @Failure 404 {object} object
// @router /account/subscription/view [get]
func (h *Handler) GetAccountSubscription(c *gin.Context) {
	var req getAccountSubscriptionReq

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
	err = h.LogService.LogRequestFromBot(ctx, "/account/subscription/view", int64(telegramUserID), req)
	if err != nil {
		log.Printf("Failed to log request from bot: %v\n", err.Error())
	}

	// get info about user subscription
	dataSubscription, err := h.AccountService.ViewUserSubscription(ctx, int64(telegramUserID))
	if err != nil {
		log.Printf("Failed to get info about user subscription: %v\n", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
		return
	}

	// view count of active user observers
	dataCountTier1, dataCountTier2, dataCountPercentage, err := h.AccountService.ViewCountOfUserActiveObserver(ctx, int64(telegramUserID))
	if err != nil {
		log.Printf("Failed to  view count of active user observers: %v\n", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
		return
	}

	result := make(map[string]string)
	result["subscriptionType"] = dataSubscription
	result["count_observers_tier1"] = strconv.Itoa(dataCountTier1)
	result["count_observers_tier2"] = strconv.Itoa(dataCountTier2)
	result["count_observers_percentage"] = strconv.Itoa(dataCountPercentage)

	if dataSubscription == "free" {
		result["limit_observers_tier1"] = os.Getenv("FreeSubscription_CountOfObserversTier1")
		result["limit_observers_tier2"] = os.Getenv("FreeSubscription_CountOfObserversTier2")
		result["limit_observers_percentage"] = os.Getenv("FreeSubscription_CountOfPercentObservers")
		result["time_to_observe"] = os.Getenv("FreeSubscription_TimeToObserve_Seconds")
	} else if dataSubscription == "premium" {
		result["limit_observers_tier1"] = os.Getenv("PaidSubscription_CountOfObserversTier1")
		result["limit_observers_tier2"] = os.Getenv("PaidSubscription_CountOfObserversTier2")
		result["limit_observers_percentage"] = os.Getenv("PaidSubscription_CountOfPercentObservers")
		result["time_to_observe"] = os.Getenv("PaidSubscription_TimeToObserve_Seconds")
	}

	c.JSON(http.StatusOK, gin.H{
		"result": result,
	})
}
