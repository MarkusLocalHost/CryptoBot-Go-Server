package handlers

import (
	"cryptocurrency/internal/middleware"
	"cryptocurrency/internal/models/mongo_models"
	"cryptocurrency/internal/utils/apperrors"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type createPromoCodeReq struct {
	UserID            string `json:"userID"`
	Title             string `json:"title"`
	Value             string `json:"value"`
	SubscriptionHours string `json:"subscriptionHours"`
	CountOfActivation string `json:"countOfActivation"`
	ActiveBefore      string `json:"activeBefore"`
}

// CreatePromoCode Create price observer
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
func (h *Handler) CreatePromoCode(c *gin.Context) {
	var req createPromoCodeReq

	// get context data
	middlewareData, _ := c.Get(middleware.AuthorizationPayloadKey)
	data := middlewareData.(string)
	data = strings.ReplaceAll(data, "'", "\"")
	err := json.Unmarshal([]byte(data), &req)
	if err != nil {
		log.Println(err)
	}

	// convert subscriptionHours to string format "xh"
	subscriptionHours := fmt.Sprintf("%sh", req.SubscriptionHours)

	// convert count of activation to int
	countOfActivation, err := strconv.Atoi(req.CountOfActivation)
	if err != nil {
		log.Fatal(err)
	}

	// convert active before in hours to time.Time
	duration, err := strconv.Atoi(req.ActiveBefore)
	if err != nil {
		log.Fatal(err)
	}
	activeBefore := time.Now().Add(time.Duration(duration))

	p := &mongo_models.PromoCode{
		Id:                primitive.NewObjectID(),
		Title:             req.Title,
		Value:             req.Value,
		SubscriptionHours: subscriptionHours,
		CountOfActivation: countOfActivation,
		CreatedAt:         time.Now(),
		ActiveBefore:      activeBefore,
	}

	// log request
	ctx := c.Request.Context()
	err = h.LogService.LogRequestFromSite(ctx, "/manager/promo_code/create", req.UserID, req)
	if err != nil {
		log.Fatal(err)
	}

	// save observer to db
	err = h.ManagerService.CreatePromoCode(ctx, p)
	if err != nil {
		log.Printf("Failed to create promo code: %v\n", err.Error())
		c.JSON(
			apperrors.Status(err),
			gin.H{"result": err.Error()})
		return
	}

	// return result with status
	c.JSON(http.StatusCreated, gin.H{
		"result": "promo code created",
	})
}
