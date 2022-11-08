package handlers

import (
	"cryptocurrency/internal/middleware"
	"cryptocurrency/internal/models/mongo_models"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"net/http"
	"strconv"
	"time"
)

type createPercentageObserverReq struct {
	TelegramUserID     string `json:"telegramUserID"`
	Observe20Minutes   string `json:"observe20Minutes"`
	Observe60Minutes   string `json:"observe60Minutes"`
	FirstFilterType    string `json:"firstFilterType"`
	FirstFilterAmount  string `json:"firstFilterAmount"`
	SecondFilterType   string `json:"secondFilterType"`
	SecondFilterAmount string `json:"secondFilterAmount"`
}

// CreatePercentageObserver Create percentage observer
// @summary Create percentage observer
// @description Create percentage observer
// @Security Bearer
// @Tags Observers
// @Accept  json
// @Produce  json
// @Success 200 {object} object
// @Failure 401 {object} object
// @Failure 404 {object} object
// @router /observers/percentage_observer/create [get]
func (h *Handler) CreatePercentageObserver(c *gin.Context) {
	var req createPercentageObserverReq

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

	// convert FirstFilterAmount to float64
	firstFilterAmount, err := strconv.ParseFloat(req.FirstFilterAmount, 64)
	if err != nil {
		log.Printf("Failed to convert first filter value from string to float: %v\n", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
		return
	}
	secondFilterAmount, err := strconv.ParseFloat(req.SecondFilterAmount, 64)
	if err != nil {
		log.Printf("Failed to convert second filter value from string to float: %v\n", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
		return
	}

	// convert boolean
	var observe20Minutes bool
	if req.Observe20Minutes == "True" {
		observe20Minutes = true
	} else if req.Observe20Minutes == "False" {
		observe20Minutes = false
	}
	var observe60Minutes bool
	if req.Observe60Minutes == "True" {
		observe60Minutes = true
	} else if req.Observe60Minutes == "False" {
		observe60Minutes = false
	}

	o := &mongo_models.PercentageObserver{
		Id:                 primitive.NewObjectID(),
		TelegramUserID:     int64(telegramUserID),
		Observe20Minutes:   observe20Minutes,
		Observe60Minutes:   observe60Minutes,
		FirstFilterType:    req.FirstFilterType,
		FirstFilterAmount:  firstFilterAmount,
		SecondFilterType:   req.SecondFilterType,
		SecondFilterAmount: secondFilterAmount,
		CreatedAt:          time.Now(),
	}

	// log request
	ctx := c.Request.Context()
	err = h.LogService.LogRequestFromBot(ctx, "/observers/percentage_observer/create", int64(telegramUserID), req)
	if err != nil {
		log.Printf("Failed to log request from bot: %v\n", err.Error())
	}

	// save percentage observer to db
	status, err := h.ObserverService.CreatePercentageObserver(ctx, o)
	if err != nil {
		log.Printf("Failed to create percentage observer: %v\n", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
		return
	}

	// return result with status
	c.JSON(http.StatusCreated, gin.H{
		"result": status,
	})

}
