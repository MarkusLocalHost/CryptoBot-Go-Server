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

type viewUserActionsReq struct {
	UserID         string `json:"userID"`
	TelegramUserId string `bson:"telegramUserId"`
}

// ViewUserActions Create price observer
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
func (h *Handler) ViewUserActions(c *gin.Context) {
	var req viewUserActionsReq

	// get context data
	middlewareData, _ := c.Get(middleware.AuthorizationPayloadKey)
	data := middlewareData.(string)
	data = strings.ReplaceAll(data, "'", "\"")
	err := json.Unmarshal([]byte(data), &req)
	if err != nil {
		log.Println(err)
	}

	// convert count of activation to int
	telegramUserId, err := strconv.Atoi(req.TelegramUserId)
	if err != nil {
		log.Fatal(err)
	}

	// log request
	ctx := c.Request.Context()
	err = h.LogService.LogRequestFromSite(ctx, "/manager/info/user_actions", req.UserID, req)
	if err != nil {
		log.Fatal(err)
	}

	// save observer to db
	userActionsInfo, err := h.ManagerService.ViewUserActions(ctx, int64(telegramUserId))
	if err != nil {
		log.Printf("Failed to create promo code: %v\n", err.Error())
		c.JSON(
			apperrors.Status(err),
			gin.H{"error": err})
		return
	}

	// return result with status
	c.JSON(http.StatusCreated, gin.H{
		"result": userActionsInfo,
	})
}
